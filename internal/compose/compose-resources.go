package compose

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"go.verv.tech/matreshka/resources"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/compose/env"
	"github.com/Red-Sock/rscli/internal/envpatterns"
	"github.com/Red-Sock/rscli/internal/utils/copier"
	"github.com/Red-Sock/rscli/internal/utils/nums"
)

const (
	servicesPart = "services"
	networkPart  = "networks"
)

var (
	ErrInvalidComposeFileFormat = rerrors.New("invalid compose format. MUST be a VALID compose file with \"services:\" field")
	ErrInvalidComposeEnvFormat  = rerrors.New("invalid environment variable format in docker-compose file. \"${\" must be followed by \"}\"")
)

type PatternManager struct {
	Patterns map[string]Pattern
	Network  map[string]interface{} // TODO
}

type Pattern struct {
	Name                string
	ResourceType        string
	ContainerDefinition ContainerSettings
	Envs                *env.Container
}

func ReadComposePatternsFromFile(pth string) (out *PatternManager, err error) {
	out = &PatternManager{}
	// Basic compose examples: rscli built-in
	out, err = extractComposePatternsFromFile(envpatterns.BuildInComposeExamples.Content)
	if err != nil {
		return nil, rerrors.Wrap(err, "error extracting composePatterns from prepared file")
	}

	// User's defined compose examples from
	userDockerComposeExample, err := os.ReadFile(pth)
	if err != nil {
		if rerrors.Is(err, os.ErrNotExist) {
			return out, nil
		}

		return out, rerrors.Wrap(err, "error reading user defined compose file")
	}

	userServicesDefinitions, err := extractComposePatternsFromFile(userDockerComposeExample)
	if err != nil {
		return out, rerrors.Wrap(err, "error extracting composePatterns from prepared file")
	}

	for serviceName, content := range userServicesDefinitions.Patterns {
		// override predefined service description with user's defined one
		out.Patterns[serviceName] = content
	}

	if userServicesDefinitions.Network != nil {
		out.Network = userServicesDefinitions.Network
	}

	return out, nil
}

func (p *Pattern) GetName() string {
	return p.Name
}

func (p *Pattern) GetType() string {
	return p.ResourceType
}

func (p *Pattern) GetEnvs() *env.Container {
	return p.Envs
}

func (p *Pattern) GetCompose() ContainerSettings {
	return p.ContainerDefinition
}

func (p *Pattern) RenameVariable(oldName, newName string) {
	p.Envs.Rename(oldName, newName)

	for k, v := range p.ContainerDefinition.Environment {
		if strings.Contains(v, oldName) {
			p.ContainerDefinition.Environment[k] = AddEnvironmentBrackets(newName)
			break
		}
	}

	for portIdx := range p.ContainerDefinition.Ports {
		p.ContainerDefinition.Ports[portIdx] = strings.ReplaceAll(p.ContainerDefinition.Ports[portIdx], oldName, newName)
	}
}

func (c *PatternManager) GetServiceDependencies(resource resources.Resource) (*Pattern, error) {
	resType := strings.Split(resource.GetName(), "_")[0]

	originalPattern, ok := c.Patterns[resType]
	if !ok {
		// TODO handle unknown data sources
		return nil, nil
	}

	var pattern Pattern
	err := copier.Copy(&originalPattern, &pattern)
	if err != nil {
		return nil, rerrors.Wrap(err, "error coping pattern")
	}

	pattern.Name = resource.GetName()

	// if value for environment variable is defined in source config
	// e.g. section postgres_* in dev.yaml
	// data_sources:
	// 		postgres:
	//			pwd: 123
	//
	// will set environment variable for POSTGRES_PASSWORD
	// by setting variable PROJ_NAME_CAPS_RESOURCE_NAME_CAPS_PWD in .env file to 123

	envVars, err := evon.MarshalEnvWithPrefix(resource.GetType(), resource)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshalling environment variables")
	}
	for _, v := range envVars.InnerNodes {
		if envVariable, ok := pattern.ContainerDefinition.Environment[v.Name]; ok {
			pattern.Envs.AppendRaw(removeEnvironmentBrackets(envVariable), fmt.Sprint(v.Value))
		}
	}

	return &pattern, nil

}

func AddEnvironmentBrackets(in string) string {
	return "${" + in + "}"
}

func extractComposePatternsFromFile(dockerComposeFile []byte) (out *PatternManager, err error) {
	out = &PatternManager{}
	composeServices := map[string]interface{}{}
	{
		// validating prepared config
		err = yaml.Unmarshal(dockerComposeFile, composeServices)
		if err != nil {
			return nil, rerrors.Wrap(err, "error validating prepared config")
		}

		network := composeServices[networkPart]
		out.Network, _ = network.(map[string]interface{})

		examplesMap, ok := composeServices[servicesPart]
		if !ok {
			return nil, rerrors.Wrapf(ErrInvalidComposeFileFormat, "expected to have \"%s\" object", servicesPart)
		}
		composeServices = examplesMap.(map[string]interface{})
	}

	out.Patterns = make(map[string]Pattern, len(composeServices))

	for serviceName, content := range composeServices {
		cs := Pattern{
			Name:         serviceName,
			ResourceType: serviceName,
		}

		var bts []byte
		bts, err = yaml.Marshal(content)
		if err != nil {
			return nil, rerrors.Wrapf(err, "error marshaling service to yaml")
		}

		err = yaml.Unmarshal(bts, &cs.ContainerDefinition)
		if err != nil {
			return nil, rerrors.Wrap(err, "error unmarshalling service "+serviceName+" to struct")
		}

		cs.Envs, err = extractEnvsFromComposeFile(bts)
		if err != nil {
			return nil, rerrors.Wrap(err, "error extracting environment variables")
		}

		out.Patterns[serviceName] = cs
	}

	return out, nil
}

// extractEnvsFromComposeFile walks thought compose service description
// in order to find environment variable usage such as '${ENV_VAR}'
func extractEnvsFromComposeFile(b []byte) (*env.Container, error) {
	startIdx := bytes.Index(b, []byte{36, 123}) // search for "${"
	var endIdx int

	out := &env.Container{}

	for {
		if startIdx == -1 {
			break
		}

		// Looks like a black magic BUT. NEED this for the 2nd (and others) loop.
		// After next "${" was found, move pointer off last "${" location
		// TODO maybe rewrite it to something more meaningful and easy
		startIdx += endIdx + 1

		endIdx = startIdx + bytes.IndexByte(b[startIdx:], 125) // "}"
		if endIdx == -1 {
			return nil, rerrors.Wrapf(ErrInvalidComposeEnvFormat, string(b[startIdx:nums.Min(len(b)-1, 10)]))
		}

		// validate if after env variable goes ":" -> extract value
		var val []byte
		if b[endIdx+1] == 58 {
			val = b[endIdx+2:][:bytes.IndexByte(b[endIdx+1:], byte('\n'))]
			if val[len(val)-1] == '\n' {
				val = val[:len(val)-1]
			}
		}

		out.AppendRaw(string(b[startIdx+1:endIdx]), string(val))

		startIdx = bytes.Index(b[endIdx:], []byte{36, 123}) // "${"

	}

	return out, nil
}

func removeEnvironmentBrackets(in string) string {
	if in[:2] == "${" && in[len(in)-1] == '}' {
		return in[2 : len(in)-1]
	}

	return in
}
