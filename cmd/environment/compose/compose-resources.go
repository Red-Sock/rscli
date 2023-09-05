package compose

import (
	"bytes"
	_ "embed"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/compose/env"
	"github.com/Red-Sock/rscli/cmd/environment/patterns"
	"github.com/Red-Sock/rscli/internal/helpers/nums"
)

const (
	servicesPart = "services"
)

var (
	ErrInvalidComposeFileFormat = errors.New("invalid compose format. MUST be a VALID compose file with \"services:\" field")
	ErrInvalidComposeEnvFormat  = errors.New("invalid environment variable format in docker-compose file. \"${\" must be followed by \"}\"")
)

type ComposePattern struct {
	Name    string
	content ContainerSettings
	envs    env.Container
}

func ReadComposePatternsFromFile(pth string) (map[string]ComposePattern, error) {
	// Basic compose examples: rscli built-in
	out, err := extractComposePatternsFromFile(patterns.BuildInComposeExamples)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting composePatterns from prepared file")
	}

	// User's defined compose examples from
	userDockerComposeExample, err := os.ReadFile(pth)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return out, nil
		}

		return nil, errors.Wrap(err, "error reading user defined compose file")
	}

	userServicesDefinitions, err := extractComposePatternsFromFile(userDockerComposeExample)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting composePatterns from prepared file")
	}

	for serviceName, content := range userServicesDefinitions {
		// override predefined service description with user's defined one
		out[serviceName] = content
	}

	return out, nil
}

func (c *ComposePattern) GetEnvs() env.Container {
	return c.envs
}

func (c *ComposePattern) GetCompose() ContainerSettings {
	return c.content
}

func extractComposePatternsFromFile(dockerComposeFile []byte) (out map[string]ComposePattern, err error) {
	composeServices := map[string]interface{}{}
	{
		// validating prepared config
		err = yaml.Unmarshal(dockerComposeFile, composeServices)
		if err != nil {
			return nil, errors.Wrap(err, "error validating prepared config")
		}

		examplesMap, ok := composeServices[servicesPart]
		if !ok {
			return nil, errors.Wrapf(ErrInvalidComposeFileFormat, "expected to have \"%s\" object", servicesPart)
		}
		composeServices = examplesMap.(map[string]interface{})
	}

	services := make(map[string]ComposePattern, len(composeServices))

	for serviceName, content := range composeServices {
		cs := ComposePattern{
			Name: serviceName,
		}

		var bts []byte
		bts, err = yaml.Marshal(content)
		if err != nil {
			return nil, errors.Wrapf(err, "error marshaling service to yaml")
		}

		err = yaml.Unmarshal(bts, &cs.content)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling service "+serviceName+" to struct")
		}

		cs.envs, err = extractEnvsFromComposeFile(bts)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting environment variables")
		}

		services[serviceName] = cs
	}

	return services, nil
}

// extractEnvsFromComposeFile walks thought compose service description
// in order to find environment variable usage such as '${ENV_VAR}'
func extractEnvsFromComposeFile(b []byte) (env.Container, error) {
	startIdx := bytes.Index(b, []byte{36, 123}) // search for "${"
	var endIdx int

	out := env.Container{}

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
			return env.Container{}, errors.Wrapf(ErrInvalidComposeEnvFormat, string(b[startIdx:nums.Min(len(b)-1, 10)]))
		}

		// validate if after env variable goes ":" -> extract value
		var val []byte
		if b[endIdx+1] == 58 {
			val = b[endIdx+2:][:bytes.IndexByte(b[endIdx+1:], byte('\n'))]
			if val[len(val)-1] == '\n' {
				val = val[:len(val)-1]
			}
		}

		out.Append(string(b[startIdx+1:endIdx]), string(val))

		startIdx = bytes.Index(b[endIdx:], []byte{36, 123}) // "${"

	}

	return out, nil
}
