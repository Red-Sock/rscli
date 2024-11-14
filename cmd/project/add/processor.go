package project

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/processor"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
)

const (
	pathFlag = "path"
)

type Proc struct {
	processor.Processor
}

func NewCommand(basicProc processor.Processor) *cobra.Command {
	proc := &Proc{
		Processor: basicProc,
	}
	c := &cobra.Command{
		Use:   "add",
		Short: "Adds resource dependency to project",
		Long:  `Can be used to add a datasource or external API dependency to project`,

		RunE: proc.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *Proc) run(cmd *cobra.Command, args []string) error {
	project, err := p.loadProject(cmd)
	if err != nil {
		return errors.Wrap(err, "error loading project for add action")
	}

	p.IO.Println(preparingMsg)
	deps := dependencies.GetDependencies(p.RscliConfig, args)
	if len(deps) == 0 {
		//	TODO return with help message
	}
	for _, d := range deps {
		err = d.AppendToProject(project)
		if err != nil {
			return errors.Wrap(err, "error adding dependency to project")
		}
	}

	p.IO.Println(startingMsg)

	ap := actions.NewActionPerformer(p.IO, project)

	err = ap.Tidy()
	if err != nil {
		return errors.Wrap(err, "error tidying project")
	}

	p.IO.Println(endMsg)

	err = git.CommitWithUntracked(project.GetProjectPath(), "added "+strings.Join(args, "; "))
	if err != nil {
		return errors.Wrap(err, "error performing git commit")
	}

	return nil
}
