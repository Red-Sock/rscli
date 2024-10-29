package init_new

import (
	"fmt"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/colors"
	"github.com/Red-Sock/rscli/plugins/project/validators"
)

const (
	askUserForNameMessagePattern = `
What would it be called?
hint: You can specify name with custom git url like "github.com/RedSock/rscli" 
      or just print name without spec symbols and spaces like "rscli"
      in this case default git-url will be "%[1]s" and final result is "%[1]s/rscli"
>`
	ackProjectNameMessagePattern = `Wonderful!!! "%s" it is!`
)

var (
	emptyNameErr = errors.New("no name entered")
)

type nameCollector struct {
	io io.IO

	defaultProjectGitPath string
}

func newNameCollector(io io.IO, defaultProjectGitPath string) nameCollector {
	return nameCollector{
		io:                    io,
		defaultProjectGitPath: defaultProjectGitPath,
	}
}

func (p *nameCollector) collect(args []string) (name string, err error) {
	if len(args) > 0 {
		name = args[0]
	} else {
		name, err = p.askUserForName()
		if err != nil {
			return "", errors.Wrap(err, "error while asking user for name")
		}
	}

	if name == "" {
		return "", errors.Wrap(emptyNameErr)
	}

	name = p.removeHttpProtoc(name)
	name = p.preAppendHost(name)
	err = validators.ValidateProjectNameStr(name)
	if err != nil {
		return "", errors.Wrap(err, "error validating project name")
	}

	name = path.Join(path.Dir(name), strings.ToLower(path.Base(name)))

	p.io.PrintlnColored(colors.ColorCyan, fmt.Sprintf(`Wonderful!!! "%s" it is!`, name))

	return name, nil
}

func (p *nameCollector) askUserForName() (name string, err error) {
	p.io.Print(fmt.Sprintf(askUserForNameMessagePattern, p.defaultProjectGitPath))

	name, err = p.io.GetInput()
	if err != nil {
		return "", errors.Wrap(err, "error obtaining project name")
	}

	return name, nil
}

func (p *nameCollector) removeHttpProtoc(name string) string {
	if strings.HasPrefix(name, "http") {
		return name[strings.Index(name, "://")+3:]
	}

	return name
}

func (p *nameCollector) preAppendHost(name string) string {
	firstDot := strings.Index(name, ".")
	firstSlash := strings.Index(name, "/")

	// if firstSlash comes after first dot - consider name already has host
	if firstSlash-firstDot > 2 {
		return name
	}

	return p.defaultProjectGitPath + "/" + name
}
