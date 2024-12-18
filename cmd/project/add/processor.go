package add

import (
	"strings"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/processor"
	"github.com/Red-Sock/rscli/plugins/project/actions"
	"github.com/Red-Sock/rscli/plugins/project/actions/git"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
)

type Proc struct {
	processor.Processor
	ActionPerformer actions.ActionPerformer
}

func NewCommand(basicProc processor.Processor) *cobra.Command {
	proc := &Proc{
		Processor:       basicProc,
		ActionPerformer: actions.NewActionPerformer(basicProc.IO),
	}
	c := &cobra.Command{
		Use:   "add",
		Short: "Adds resource dependency to project",
		Long:  `Can be used to add a datasource or external API dependency to project`,

		RunE: proc.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(
		processor.PathFlag,
		processor.PathFlag[:1],
		proc.WD,
		`path to folder with project`)

	return c
}

func (p *Proc) run(cmd *cobra.Command, args []string) error {
	project, err := p.LoadProject(cmd)
	if err != nil {
		return rerrors.Wrap(err, "error loading project for add action")
	}

	p.IO.Println(preparingMsg)

	deps := dependencies.GetDependencies(p.RscliConfig, args)
	if len(deps) == 0 {
		//	TODO return with help message
	}

	for _, d := range deps {
		err = d.AppendToProject(project)
		if err != nil {
			return rerrors.Wrap(err, "error adding dependency to project")
		}
	}

	p.IO.Println(startingMsg)

	err = p.ActionPerformer.Tidy(project)
	if err != nil {
		return rerrors.Wrap(err, "error tidying project")
	}

	p.IO.Println(endMsg)

	err = git.CommitWithUntracked(project.GetProjectPath(), "added "+strings.Join(args, "; "))
	if err != nil {
		return rerrors.Wrap(err, "error performing git commit")
	}

	return nil
}
