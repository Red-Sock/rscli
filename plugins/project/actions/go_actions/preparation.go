package go_actions

import (
	"bytes"
	stderrs "errors"
	"path"
	"strings"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/utils/renamer"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators/server_generators/impl_gen"
)

type PrepareProjectStructure struct {
}

func (a PrepareProjectStructure) Do(p project.IProject) error {
	rootF := p.GetFolder()

	cmd := &folder.Folder{Name: patterns.CmdFolder}
	cmd.Add(patterns.MainFile.CopyWithNewName(path.Join(patterns.ServiceFolder, patterns.MainFile.Name)))
	rootF.Add(cmd)

	rootF.Add(&folder.Folder{Name: patterns.ConfigsFolder})
	rootF.Add(&folder.Folder{Name: patterns.InternalFolder})

	rootF.Add(
		patterns.Dockerfile.Copy(),
		patterns.Readme.Copy(),
		patterns.GitIgnore.Copy(),
		patterns.Linter.Copy(),
	)

	return nil
}
func (a PrepareProjectStructure) NameInAction() string {
	return "Preparing project structure"
}

type PrepareClients struct {
	C  *rscliconfig.RsCliConfig
	IO io.IO
}

func (a PrepareClients) Do(p project.IProject) error {
	if a.C == nil {
		a.C = rscliconfig.GetConfig()
	}

	if a.IO == nil {
		a.IO = io.StdIO{}
	}

	var simpleClients []string
	var grpcClients []string
	cfg := p.GetConfig()

	for _, r := range cfg.DataSources {
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

	if len(errs) != 0 {
		return stderrs.Join(errs...)
	}

	return nil
}

func (a PrepareClients) NameInAction() string {
	return "Generating clients"
}

type PrepareMakefile struct{}

func (a PrepareMakefile) Do(p project.IProject) error {
	genScriptSummary := make([]string, 0)

	// first part for summary scripts
	makefileContent := make([][]byte, 1, 4)
	{
		// basic info
		rscliBasicScript := make([]byte, len(patterns.RscliMK.Content))
		copy(rscliBasicScript, patterns.RscliMK.Content)

		rscliBasicScript = renamer.ReplaceProjectNameShort(rscliBasicScript, p.GetShortName())

		makefileContent = append(makefileContent, rscliBasicScript)
	}

	if len(p.GetConfig().Servers) != 0 {
		// basic info
		serverGenCopy := make([]byte, len(patterns.GrpcServerGenMK))
		copy(serverGenCopy, patterns.GrpcServerGenMK)

		makefileContent = append(makefileContent, append([]byte(`### Grpc server generation`+"\n"), serverGenCopy...))
		genScriptSummary = append(genScriptSummary, patterns.GenGrpcServerCommand)
	}

	rscliMk := p.GetFolder().GetByPath(patterns.RscliMakefileFile)
	if rscliMk == nil {
		p.GetFolder().Add(&folder.Folder{
			Name: patterns.RscliMakefileFile,
		})
		rscliMk = p.GetFolder().GetByPath(patterns.RscliMakefileFile)
	}

	if len(genScriptSummary) != 0 {
		makefileContent[0] = []byte(patterns.GenCommand + ": " + strings.Join(genScriptSummary, " "))
	} else {
		makefileContent = makefileContent[1:]
	}

	rscliMk.Content = bytes.Join(makefileContent, []byte{'\n'})

	makefile := p.GetFolder().GetByPath(patterns.MakefileFile)
	if makefile == nil {
		p.GetFolder().Add(patterns.Makefile.Copy())
	}

	return nil
}
func (a PrepareMakefile) NameInAction() string {
	return "Generating Makefile"
}

type PrepareServer struct{}

func (a PrepareServer) Do(p project.IProject) error {
	if len(p.GetConfig().Servers) == 0 {
		return nil
	}

	rootF := p.GetFolder()

	if rootF.GetByPath(patterns.EasyP.Name) == nil {
		rootF.Add(patterns.EasyP.Copy())
	}

	transportFolder := rootF.GetByPath(patterns.InternalFolder, patterns.TransportFolder)
	if transportFolder == nil {
		transportFolder = &folder.Folder{}
	}

	transportFolder.Add(patterns.ServerManager.Copy())
	transportFolder.Add(patterns.GrpcServerManager.Copy())
	transportFolder.Add(patterns.HttpServerManager.Copy())

	implFolders, err := impl_gen.GenerateImpl(rscliconfig.GetConfig(), p)
	if err != nil {
		return rerrors.Wrap(err, "error during stub generation")
	}
	for _, implF := range implFolders {
		exists := false
		for _, tF := range transportFolder.Inner {
			if tF.Name == implF.Name {
				exists = true
				break
			}
		}

		if !exists {
			transportFolder.Add(implFolders...)
		}
	}

	return nil
}
func (a PrepareServer) NameInAction() string {
	return "Preparing server files"
}
