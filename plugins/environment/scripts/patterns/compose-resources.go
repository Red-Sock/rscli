package patterns

import (
	"bytes"
	_ "embed"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/utils/nums"
)

const (
	servicesPart = "services"
)

var (
	ErrInvalidComposeFileFormat = errors.New("invalid compose format. MUST be a VALID compose file with \"services:\" field")
	ErrInvalidComposeEnvFormat  = errors.New("invalid environment variable format in docker-compose file. \"${\" must be followed by \"}\"")
)

type ComposePatterns struct {
	Name    string
	content ContainerSettings
	envs    *EnvService
}

func NewComposePatterns(wd string) (services map[string]ComposePatterns, err error) {
	var out map[string]ComposePatterns
	{
		out, err = extractComposePatternsFromServices(buildInExamples)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting patterns from prepared file")
		}
	}

	{
		var userDockerComposeExample []byte
		userDockerComposeExample, err = os.ReadFile(path.Join(wd, EnvDir, DockerComposeFile.Name))
		if err != nil {
			return nil, errors.Wrap(err, "error reading user defined compose file")
		}

		var userServicesDefinitions map[string]ComposePatterns
		userServicesDefinitions, err = extractComposePatternsFromServices(userDockerComposeExample)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting patterns from prepared file")
		}

		for serviceName, content := range userServicesDefinitions {
			// override predefined service description with user's defined one
			out[serviceName] = content
		}
	}

	return out, nil
}

func (c *ComposePatterns) GetEnvs() *EnvService {
	return c.envs
}

func (c *ComposePatterns) GetCompose() ContainerSettings {
	return c.content
}

func extractComposePatternsFromServices(dockerComposeFile []byte) (out map[string]ComposePatterns, err error) {
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

	services := make(map[string]ComposePatterns, len(composeServices))

	for serviceName, content := range composeServices {
		cs := ComposePatterns{
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
		cs.envs, err = extractEnvsFromComposeService(bts)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting environment variables")
		}

		services[serviceName] = cs
	}

	return services, nil
}

// extractEnvsFromComposeService walks thought compose service description
// in order to find environment variable usage such as '${ENV_VAR}'
func extractEnvsFromComposeService(b []byte) (*EnvService, error) {
	startIdx := bytes.Index(b, []byte{36, 123}) // search for "${"
	var endIdx int

	out, _ := NewEnvService(nil)

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
			return &EnvService{}, errors.Wrapf(ErrInvalidComposeEnvFormat, string(b[startIdx:nums.Min(len(b)-1, 10)]))
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
