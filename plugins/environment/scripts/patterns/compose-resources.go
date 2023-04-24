package patterns

import (
	"bytes"
	_ "embed"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/Red-Sock/rscli/internal/utils/nums"
)

const (
	servicesPart = "services"
)

var ErrCasting = errors.New("type casting error")

var (
	ErrInvalidComposeConstructorArgument = errors.New("invalid argument in compose constructor")
	ErrInvalidComposeFileFormat          = errors.New("invalid compose format. MUST be a VALID compose file with \"services:\" field")
	ErrInvalidComposePortFormat          = errors.New("invalid ports format in docker-compose file. Must be \"port1:port2\"")
	ErrInvalidComposeEnvFormat           = errors.New("invalid environment variable format in docker-compose file. \"${\" must be followed by \"}\"")
)

//go:embed files/compose.examples.yaml
var composeExamples []byte

type ComposeService struct {
	Name    string
	content ContainerSettings
	envs    *EnvService
}

func NewComposeServices(src ...[]byte) (services map[string]ComposeService, err error) {
	examples := map[string]interface{}{}

	{
		var srcFile []byte
		switch len(src) {
		case 0:
			srcFile = composeExamples
		case 1:
			srcFile = src[0]
		default:
			return nil, errors.Wrapf(ErrInvalidComposeConstructorArgument, "Expected 0 or 1 src files to inherit from. Got %d ", len(srcFile))
		}

		// validating input
		err = yaml.Unmarshal(srcFile, examples)
		if err != nil {
			return nil, err
		}

		smths, ok := examples[servicesPart]
		if !ok {
			return nil, errors.Wrapf(ErrInvalidComposeFileFormat, "expected to have %s ", servicesPart)
		}

		examples = smths.(map[string]interface{})
	}

	services = make(map[string]ComposeService, len(examples))

	for serviceName, content := range examples {
		cs := ComposeService{
			Name: serviceName,
		}
		var bts []byte
		bts, err = yaml.Marshal(content)
		if err != nil {
			return nil, errors.Wrapf(err, "error marshaling service to yaml")
		}
		err = yaml.Unmarshal(bts, &cs.content)

		cs.envs, err = extractEnvs(bts)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting environment variables")
		}

		services[serviceName] = cs
	}

	return services, nil
}

func (c *ComposeService) GetEnvs() *EnvService {
	return c.envs
}

func (c *ComposeService) GetCompose() ContainerSettings {
	return c.content
}

func extractEnvs(b []byte) (*EnvService, error) {

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
