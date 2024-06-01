package go_actions

import (
	stderrs "errors"
	"os"
	"path"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"
	"github.com/godverv/matreshka/servers"

	"github.com/Red-Sock/rscli/internal/cmd"
	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/utils/bins/makefile"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
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
	renamer.ReplaceProjectName(p.GetName(), p.GetFolder())

	err := p.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building project")
	}

	b, err := p.GetConfig().Marshal()
	if err != nil {
		return errors.Wrap(err, "error marshaling config")
	}

	err = os.WriteFile(p.GetConfig().Path, b, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error writing config to file")
	}

	_, err = cmd.Execute(cmd.Request{
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

type GenerateClientsAction struct {
	C  *rscliconfig.RsCliConfig
	IO io.IO
}

func (a GenerateClientsAction) Do(p interfaces.Project) error {
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

func (a GenerateClientsAction) NameInAction() string {
	return "Generating clients"
}

type GenerateServerAction struct {
	C  *rscliconfig.RsCliConfig
	IO io.IO
}

func (a GenerateServerAction) Do(p interfaces.Project) error {
	err := makefile.Install()
	if err != nil {
		return errors.Wrap(err, "error installing makefile")
	}

	var errs []error
	for _, f := range p.GetFolder().Inner {
		if !strings.HasSuffix(f.Name, ".mk") {
			continue
		}
		switch f.Name {
		case projpatterns.GrpcMK.Name:
			err = makefile.Run(p.GetProjectPath(), f.Name, "gen")
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) != 0 {
		return errors.Wrap(stderrs.Join(errs...), "error generating files")
	}

	return nil
}

func (a GenerateServerAction) NameInAction() string {
	return "Generating server"
}

type GenerateMakefileAction struct{}

func (a GenerateMakefileAction) Do(p interfaces.Project) error {
	scriptsFolder := p.GetFolder().GetByPath(projpatterns.ScriptsFolder)
	if scriptsFolder == nil {
		scriptsFolder = &folder.Folder{
			Name: projpatterns.ScriptsFolder,
		}
		p.GetFolder().Add(scriptsFolder)
	}

	if scriptsFolder.GetByPath(projpatterns.GrpcMK.Name) == nil {
		for _, s := range p.GetConfig().AppConfig.Servers {
			_, ok := s.(*servers.GRPC)
			if ok {
				scriptsFolder.Add(projpatterns.GrpcMK.Copy())
				return nil
			}
		}
	}

	return nil
}

func (a GenerateMakefileAction) NameInAction() string {
	return "Generating Makefile"
}
