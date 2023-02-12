package patterns

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/Red-Sock/rscli/internal/utils/nums"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

var (
	ErrInvalidComposeFileFormat = errors.New("invalid compose format. MUST be a VALID compose file with \"services:\" field")
	ErrCasting                  = errors.New("type casting error")
	ErrInvalidPortFormat        = errors.New("invalid ports format in docker-compose file. Must be \"port1:port2\"")
	ErrInvalidEnvFormat         = errors.New("invalid environment variable format in docker-compose file. \"${\" must be followed by \"}\"")
)

//go:embed files/compose.examples.yaml
var composeExamples []byte

type ComposeService struct {
	content       []byte
	dependencies  []string
	startingPorts []int
	envs          string
}

func NewComposeService() (services map[string]ComposeService, err error) {
	srcFile := map[string]interface{}{}

	{
		// validating input
		err := yaml.Unmarshal(composeExamples, srcFile)
		if err != nil {
			return nil, err
		}
		smths, ok := srcFile["services"]
		if !ok {
			return nil, ErrInvalidComposeFileFormat
		}
		srcFile = smths.(map[string]interface{})
	}

	services = make(map[string]ComposeService, len(srcFile))

	for serviceName, content := range srcFile {
		cs := ComposeService{}

		cs.content = []byte(fmt.Sprintf("%s", content))

		contentMap, ok := content.(map[string]interface{})
		if !ok {
			return nil, errors.Wrapf(ErrCasting, "errors casting %v if %T type to map[string]interface{}", content, content)
		}

		cs.dependencies, err = extractDependencies(contentMap)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting dependencies on other services")
		}

		cs.startingPorts, err = extractStartingPorts(contentMap)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting ports")
		}

		cs.envs, err = extractEnvs(cs.content)
		if err != nil {
			return nil, errors.Wrap(err, "error extracting environment variables")
		}
		services[serviceName] = cs
	}

	return services, nil
}

func (c *ComposeService) GetEnvs() string {
	return c.envs
}

func (c *ComposeService) GetStartingPorts() []int {
	return c.startingPorts
}

func (c *ComposeService) GetCompose() []byte {
	return c.content
}

func extractDependencies(contentMap map[string]interface{}) ([]string, error) {
	deps, ok := contentMap["depends_on"]
	if !ok {
		return nil, nil
	}

	depsArr, ok := deps.([]interface{})
	if !ok {
		return nil, errors.Wrapf(ErrCasting, "%v to []interfaces{}. actual type is %T", deps, deps)
	}

	dependencies := make([]string, len(depsArr))
	for _, d := range depsArr {
		dependencies = append(dependencies, fmt.Sprintf("%v", d))
	}

	return dependencies, nil
}

func extractStartingPorts(b map[string]interface{}) ([]int, error) {
	ports, ok := b["ports"]
	if !ok {
		return nil, nil
	}

	portsArr, ok := ports.([]interface{})
	if !ok {
		return nil, errors.Wrapf(ErrCasting, "%v of type %T to []interface", ports, ports)
	}

	portsInt := make([]int, 0, len(portsArr))

	for _, p := range portsArr {
		portString := fmt.Sprintf("%v", p)
		portSplit := strings.Split(portString, ":")
		if len(portSplit) != 2 {
			return nil, errors.Wrapf(ErrInvalidPortFormat, "got: %s", portString)
		}

		portInt, err := strconv.Atoi(portSplit[1])
		if err != nil {
			return nil, errors.Wrapf(err, "error converting %s to int", portSplit[1])
		}
		portsInt = append(portsInt, portInt)
	}

	return portsInt, nil
}

func extractEnvs(b []byte) (string, error) {
	startIdx := bytes.Index(b, []byte{36, 123}) // "${"
	var endIdx int

	sb := strings.Builder{}

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
			return "", errors.Wrapf(ErrInvalidEnvFormat, string(b[startIdx:nums.Min(len(b)-1, 10)]))
		}

		sb.Write(b[startIdx+1 : endIdx])
		sb.WriteByte(10) // "=\n"

		startIdx = bytes.Index(b[endIdx+1:], []byte{36, 123}) // "${"

	}

	return sb.String(), nil
}
