package tidy

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/internal/utils/slices"
	"github.com/Red-Sock/rscli/plugins/project/config"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

var ErrNoMainFile = errors.New("no main file was found")

const (
	waitingForTheEndFunc = "waitingForTheEnd"
	// with \n to be sure if correct import
	// and not a commentary will be recognized as a start point of import
	importWord          = "\nimport"
	goFuncWord          = "go func() {"
	transportNewManager = "transport.NewManager()"
)

func Api(p interfaces.Project) error {
	cfg := p.GetConfig()

	pathToMainFile := []string{projpatterns.CmdFolder, p.GetShortName(), projpatterns.MainFile.Name}

	projMainFile := p.GetFolder().GetByPath(pathToMainFile...)
	if projMainFile == nil {
		return errors.Wrap(ErrNoMainFile, strings.Join(pathToMainFile, "/"))
	}

	err := tidyAPI(p, cfg, projMainFile)
	if err != nil {
		return errors.Wrap(err, "error tiding API")
	}

	return nil
}

func tidyAPI(p interfaces.Project, cfg *config.Config, projMainFile *folder.Folder) error {
	// todo
	//serverFolders, err := cfg.GetServerFolders()
	//if err != nil {
	//	return err
	//}

	//if serverFolders == nil {
	//	return nil
	//}
	//
	//insertApiSetupInMainIfNotExists(p, projMainFile)
	//
	//err = tidyAPIFile(p, serverFolders)
	//if err != nil {
	//	return errors.Wrap(err, "error tiding api file")
	//}

	return nil
}

func insertApiSetupInMainIfNotExists(p interfaces.Project, projMainFile *folder.Folder) {
	bootstrapFolderPath := path.Join(projpatterns.CmdFolder, p.GetShortName(), projpatterns.BootStrapFolder)
	apiFile := p.GetFolder().GetByPath(bootstrapFolderPath, projpatterns.APISetupFile.Name)

	// if bootstrap for server doesn't exist - adding it
	if apiFile == nil {
		apiFile = &folder.Folder{
			Name:    path.Join(bootstrapFolderPath, projpatterns.APISetupFile.Name),
			Content: projpatterns.APISetupFile.Content,
		}
		p.GetFolder().Add(apiFile)
	}

	const (
		// key lines for starting and stopping servers
		apiEntryPointCall = `bootstrap.ApiEntryPoint`

		apiEntryPointStop = `stopFunc(context.Background())`
	)

	var insertBeforeEnd []byte
	var insertAfterEnd []byte

	// add to main file api entry call if not exists
	if bytes.Index(projMainFile.Content, []byte(apiEntryPointCall)) == -1 {
		insertBeforeEnd = []byte(fmt.Sprintf(`
stopFunc, err := %s(ctx, cfg)
if err != nil {
	logrus.Fatal(err)
}
`, apiEntryPointCall))
	}
	wfteBytes := []byte(waitingForTheEndFunc)

	endFuncIdx := bytes.Index(projMainFile.Content, wfteBytes)
	endFuncIdx = bytes.LastIndex(projMainFile.Content[:endFuncIdx], []byte("\n"))
	if len(insertBeforeEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertBeforeEnd, endFuncIdx)
		endFuncIdx = bytes.Index(projMainFile.Content, wfteBytes) + len(wfteBytes)
		endFuncIdx = endFuncIdx + bytes.Index(projMainFile.Content[endFuncIdx:], []byte("\n")) + 1
	}

	if bytes.Index(projMainFile.Content, []byte(apiEntryPointStop)) == -1 {
		insertAfterEnd = append(
			insertAfterEnd,
			[]byte(fmt.Sprintf(`
err = %s
if err != nil {
	logrus.Fatal(err)
}
`, apiEntryPointStop))...)
	}
	if len(insertAfterEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertAfterEnd, endFuncIdx)
	}

	{
		// add import on boostrap if doesn't exists
		importBootstrap := []byte("\"" + p.GetName() + "/cmd/" + p.GetShortName() + "/bootstrap\"\n")
		if bytes.Index(projMainFile.Content, importBootstrap) == -1 {
			importStartIdx := bytes.Index(projMainFile.Content, []byte(importWord))
			importEndIdx := importStartIdx + bytes.Index(projMainFile.Content[importStartIdx:], []byte(")"))
			projMainFile.Content = slices.InsertSlice(
				projMainFile.Content,
				importBootstrap,
				importEndIdx,
			)
		}
	}

	return
}
func tidyAPIFile(p interfaces.Project, serverFolders []*folder.Folder) error {
	bootstrapFolderPath := path.Join(projpatterns.CmdFolder, p.GetShortName(), projpatterns.BootStrapFolder)

	apiFile := p.GetFolder().GetByPath(bootstrapFolderPath, projpatterns.APISetupFile.Name)
	if apiFile == nil {
		apiFile = projpatterns.APISetupFile.CopyWithNewName(
			path.Join(bootstrapFolderPath, projpatterns.APISetupFile.Name))
		p.GetFolder().Add(apiFile)
	}

	err := insertMissingAPI(p, serverFolders, apiFile)
	if err != nil {
		return errors.Wrap(err, "error inserting missing api")
	}

	apiMgr := p.GetFolder().GetByPath(projpatterns.InternalFolder, projpatterns.TransportFolder, projpatterns.ApiManagerFileName)
	if apiMgr == nil {
		serverFolders = append(serverFolders,
			projpatterns.ServerManagerPatternFile.CopyWithNewName(
				path.Join(projpatterns.InternalFolder, projpatterns.TransportFolder, projpatterns.ServerManagerPatternFile.Name)),
		)
	}

	p.GetFolder().Add(serverFolders...)

	return nil
}

func insertMissingAPI(p interfaces.Project, serverFolders []*folder.Folder, httpFile *folder.Folder) error {
	serversInit := extractApiInit(httpFile.Content)

	var newAPIInsert []byte
	var newAPIImportInsert []byte

	var grpcServers []string

	for _, serv := range serverFolders {
		if bytes.Contains(serversInit, []byte(serv.Name)) {
			continue
		}

		newAPIImportInsert = append(newAPIImportInsert, []byte("\n\t\""+p.GetName()+"/internal/transport/"+serv.Name+"\"")...)
		newAPIInsert = append(newAPIInsert, []byte("mngr.AddServer("+serv.Name+".NewServer(cfg))\n\t")...)

		// TODO
		//if strings.Contains(serv.Name, servers.GRp) {
		//	grpcServers = append(grpcServers, serv.Name)
		//}
	}

	if len(newAPIImportInsert) != 0 {
		var importEndIdx int
		{
			importStartIdx := bytes.Index(httpFile.Content, []byte(importWord))
			importEndIdx = importStartIdx + bytes.Index(httpFile.Content[importStartIdx:], []byte(")"))
		}

		httpFile.Content = slices.InsertSlice(
			httpFile.Content,
			newAPIImportInsert,
			importEndIdx,
		)

		httpFile.Content = slices.InsertSlice(
			httpFile.Content,
			newAPIInsert,
			bytes.Index(httpFile.Content, []byte(goFuncWord)),
		)
	}

	if len(grpcServers) != 0 {
		for _, serverName := range grpcServers {
			protoFolder := p.GetFolder().GetByPath(projpatterns.PkgFolder, projpatterns.ProtoFolder, serverName)
			if protoFolder == nil {

				exampleFile := projpatterns.GrpcProtoExampleFile.Copy()
				exampleFile.Content = bytes.ReplaceAll(
					exampleFile.Content,
					[]byte(projpatterns.ImportProjectNamePatternSnakeCase),
					[]byte(cases.KebabToSnake(serverName)))

				exampleFile.Content = bytes.ReplaceAll(exampleFile.Content,
					[]byte("grpc_realisation"),
					[]byte(serverName),
				)
				p.GetFolder().Add(
					&folder.Folder{
						Name: path.Join(projpatterns.PkgFolder, projpatterns.ProtoFolder, serverName),
						Inner: []*folder.Folder{
							{
								Name:    serverName + projpatterns.ProtoFileExtension,
								Content: exampleFile.Content,
							},
						},
					})
			}
		}

		err := insertGRPCInMakefile(p)
		if err != nil {
			return errors.Wrap(err, "error in insertGRPCInMakefile")
		}

		err = p.GetFolder().Build()
		if err != nil {
			return errors.Wrap(err, "error building project with proto and grpc calls in makefile")
		}

		_, err = cmd.Execute(cmd.Request{
			Tool:    "make",
			Args:    []string{"-f", "rscli.mk", "generate-proto"},
			WorkDir: p.GetProjectPath(),
		})
		if err != nil {
			return errors.Wrap(err, "error generating grpc via make generate-proto")
		}
	}
	return nil
}

func insertGRPCInMakefile(p interfaces.Project) error {
	rscliMkF := p.GetFolder().GetByPath(projpatterns.RscliMK.Name)
	if rscliMkF == nil {
		return ErrNoMakeFile
	}

	protocMkfile := bytes.Join([][]byte{
		projpatterns.GRPCSection,
		projpatterns.GRPCUtilityInstallGoProtocHeader,
		projpatterns.GRPCGenerateGoCodeWithDependencies,
		projpatterns.GRPCInstallProtocViaGolangEnvOSBased,
		projpatterns.GRPCInstallGolangProtoc,
		projpatterns.GRPCGenerateGoCode,
		projpatterns.GRPCGatewayDependency,
		projpatterns.SectionSeparator,
	},
		[]byte{})

	idxStart := bytes.Index(rscliMkF.Content, projpatterns.GRPCUtilityInstallGoProtocHeader)
	if idxStart == -1 {
		rscliMkF.Content = append(rscliMkF.Content, protocMkfile...)
	} else {
		idxEnd := idxStart + bytes.Index(rscliMkF.Content[idxStart:], projpatterns.SectionSeparator)
		newContent := make([]byte, idxStart+len(protocMkfile)+len(rscliMkF.Content[idxEnd:]))
		copy(newContent[:idxStart], rscliMkF.Content[:idxStart])
		copy(newContent[idxStart:idxEnd], protocMkfile)
		copy(newContent[idxEnd:], rscliMkF.Content[idxEnd:])
	}

	return nil
}

func extractApiInit(httpFile []byte) (out []byte) {
	goFuncWordBytes := []byte(goFuncWord)
	// indexes between creation of transport manager
	// and starting it in goroutine
	startIdx := bytes.Index(httpFile, []byte(transportNewManager)) + len(transportNewManager) + 2
	endIdx := bytes.Index(httpFile, goFuncWordBytes)

	out = make([]byte, endIdx-startIdx)
	copy(out, httpFile[startIdx:endIdx])

	return out
}
