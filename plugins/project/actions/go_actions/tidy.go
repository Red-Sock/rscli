package go_actions

import (
	"bytes"
	stderrs "errors"
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/cmd"
	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/utils/bins/makefile"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

const goBin = "go"

type InitGoModAction struct{}

func (a InitGoModAction) Do(p interfaces.Project) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "init", p.GetName()},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return errors.Wrap(err, "error executing go mod init")
	}

	goMod, err := os.OpenFile(path.Join(p.GetProjectPath(), "go.mod"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer func() {
		err2 := goMod.Close()
		if err2 != nil {
			if err != nil {
				err = errors.Wrap(err, "error on closing"+err2.Error())
			} else {
				err = err2
			}
		}
	}()

	return nil
}
func (a InitGoModAction) NameInAction() string {
	return "Initiating go project"
}

type FormatAction struct{}

func (a FormatAction) Do(p interfaces.Project) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"fmt", "./..."},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return err
	}

	return nil
}
func (a FormatAction) NameInAction() string {
	return "Performing project fix up"
}

type TidyAction struct{}

func (a TidyAction) Do(p interfaces.Project) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "tidy"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return errors.Wrap(err, "error executing go mod tidy")
	}

	err = FormatAction{}.Do(p)
	if err != nil {
		return errors.Wrap(err, "error formatting project")
	}

	return nil
}
func (a TidyAction) NameInAction() string {
	return "Cleaning up the project"
}

type PrepareClientsAction struct {
	C  *rscliconfig.RsCliConfig
	IO io.IO
}

func (a PrepareClientsAction) Do(p interfaces.Project) error {
	if a.C == nil {
		a.C = rscliconfig.GetConfig()
	}

	if a.IO == nil {
		a.IO = io.StdIO{}
	}

	var simpleClients []string
	var grpcClients []string

	for _, r := range p.GetConfig().DataSources {
		grpcC, ok := r.(*resources.GRPC)
		if ok {
			grpcClients = append(grpcClients, grpcC.Module)
		} else {
			simpleClients = append(simpleClients, r.GetName())
		}
	}
	var errs []error

	deps := dependencies.GetDependencies(a.C, simpleClients)
	if len(deps) != 0 {
		for _, item := range deps {
			err := item.AppendToProject(p)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	err := dependencies.GrpcClient{
		Modules: grpcClients,
		Cfg:     a.C,
		Io:      a.IO,
	}.AppendToProject(p)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return stderrs.Join(errs...)
	}

	return nil
}
func (a PrepareClientsAction) NameInAction() string {
	return "Generating clients"
}

type GenerateServerAction struct {
	C  *rscliconfig.RsCliConfig
	IO io.IO
}

func (a GenerateServerAction) Do(p interfaces.Project) error {
	if len(p.GetConfig().Servers) == 0 {
		return nil
	}

	err := makefile.Install()
	if err != nil {
		return errors.Wrap(err, "error installing makefile")
	}

	err = makefile.Run(p.GetProjectPath(), projpatterns.Makefile, projpatterns.GenCommand)
	if err != nil {
		return errors.Wrap(err, "error generating")
	}
	return nil
}
func (a GenerateServerAction) NameInAction() string {
	return "Generating server"
}

type PrepareMakefileAction struct{}

func (a PrepareMakefileAction) Do(p interfaces.Project) error {

	genScriptSummary := make([]string, 0)

	// first part for summary scripts
	makefileContent := make([][]byte, 1, 4)

	{
		// basic info
		rscliCopy := make([]byte, len(projpatterns.RscliMK))
		copy(rscliCopy, projpatterns.RscliMK)
		makefileContent = append(makefileContent, append([]byte(`### General Rscli info`+"\n"), rscliCopy...))
	}

	if len(p.GetConfig().Servers) != 0 {
		// basic info
		serverGenCopy := make([]byte, len(projpatterns.GrpcServerGenMK))
		copy(serverGenCopy, projpatterns.GrpcServerGenMK)

		makefileContent = append(makefileContent, append([]byte(`### Grpc server generation`+"\n"), serverGenCopy...))
		genScriptSummary = append(genScriptSummary, projpatterns.GenGrpcServerCommand)
	}

	makeFile := p.GetFolder().GetByPath(projpatterns.Makefile)
	if makeFile == nil {
		p.GetFolder().Add(&folder.Folder{
			Name: projpatterns.Makefile,
		})
		makeFile = p.GetFolder().GetByPath(projpatterns.Makefile)
	}

	if len(genScriptSummary) != 0 {
		makefileContent[0] = []byte(projpatterns.GenCommand + ": " + strings.Join(genScriptSummary, " "))
	}

	makeFile.Content = bytes.Join(makefileContent, []byte{'\n', '\n'})

	return nil
}
func (a PrepareMakefileAction) NameInAction() string {
	return "Generating Makefile"
}
