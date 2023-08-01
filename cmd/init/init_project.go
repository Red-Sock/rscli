package init

import (
	"fmt"
	"strings"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/pkg/colors"
	"github.com/Red-Sock/rscli/plugins/project/processor"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Initializes project",
	Long:  `Can be used to init a project via configuration file, constructor or global config`,
	RunE:  projectConstructorImp.initProject,
}

var projectConstructorImp = projectConstructor{
	cfg: config.GetConfig(),
	io:  stdio.StdIO{},
}

type projectConstructor struct {
	cfg *config.RsCliConfig
	io  stdio.IO
}

func (p *projectConstructor) initProject(_ *cobra.Command, _ []string) error {

	constructor := processor.CreateArgs{}

	var err error
	constructor.Name, err = p.obtainName()
	if err != nil {
		return errors.Wrap(err, "error obtaining name")
	}

	p.io.PrintlnColored(colors.ColorCyan, fmt.Sprintf(`Wonderful!!! "%s" it is!`, constructor.Name))

	return nil
}

func (p *projectConstructor) obtainName() (string, error) {
	p.io.Print(fmt.Sprintf(`
What would it be called?
hint: You can specify name with custom git url like "github.com/RedSock/rscli" 
      or just print name without spec symbols and spaces like "rscli"
      in this case default git-url will be "%s" and final result is "%s/rscli"
>`, p.cfg.DefaultProjectGitPath, p.cfg.DefaultProjectGitPath))

	name, err := p.io.GetInput()
	if err != nil {
		return "", errors.Wrap(err, "error obtaining project name")
	}

	if strings.HasPrefix(name, "http") {
		name = name[strings.Index(name, "://")+3:]
	}

	var containsHost bool

	// if first part (before first "/" symbol) contains dot "." -> consider it's overriding default repository
	if leftPartIdx := strings.Index(name, "/"); leftPartIdx != -1 {
		hostSeparatorIdx := strings.Index(name[:leftPartIdx], ".")
		containsHost = hostSeparatorIdx != -1
	}

	if !containsHost {
		name = p.cfg.DefaultProjectGitPath + "/" + name
	}

	err = p.validateName(name)
	if err != nil {
		return "", errors.Wrap(err, "error validating name")
	}

	return name, nil
}
func (p *projectConstructor) validateName(name string) error {
	if name == "" {
		return errors.New("no name entered")
	}

	// starting and ending ascii symbols ranges that are applicable to project name
	availableRanges := [][]int32{
		{45, 47},
		{48, 57},
		{65, 90},
		{97, 122},
	}
	for _, s := range name {
		var hasHitRange = false
		for _, r := range availableRanges {
			if s >= r[0] && s <= r[1] {
				hasHitRange = true
				break
			}
		}
		if !hasHitRange {
			return errors.New("name contains \"" + string(s) + "\" symbol")
		}
	}

	return nil
}
