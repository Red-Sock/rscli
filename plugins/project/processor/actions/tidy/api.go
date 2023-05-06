package tidy

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/internal/utils/slices"
	"github.com/Red-Sock/rscli/pkg/cmd"
	"github.com/Red-Sock/rscli/pkg/folder"
	_const "github.com/Red-Sock/rscli/plugins/config/pkg/const"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
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

	projMainFile := p.GetFolder().GetByPath(patterns.CmdFolder, p.GetName(), patterns.MainFileName)
	if projMainFile == nil {
		return errors.Wrap(ErrNoMainFile, strings.Join([]string{patterns.CmdFolder, p.GetName(), patterns.MainFileName}, "/"))
	}

	err := tidyAPI(p, cfg, projMainFile)
	if err != nil {
		return errors.Wrap(err, "error tiding API")
	}

	return nil
}

func tidyAPI(p interfaces.Project, cfg interfaces.ProjectConfig, projMainFile *folder.Folder) error {
	serverFolders, err := cfg.GetServerFolders()
	if err != nil {
		return err
	}

	if serverFolders == nil {
		return nil
	}

	insertApiSetupInMainIfNotExists(p, projMainFile)

	err = tidyAPIFile(p, serverFolders)
	if err != nil {
		return errors.Wrap(err, "error tiding api file")
	}

	return nil
}

func insertApiSetupInMainIfNotExists(p interfaces.Project, projMainFile *folder.Folder) {
	apiFile := p.GetFolder().GetByPath(patterns.CmdFolder, p.GetName(), patterns.BootStrapFolder, patterns.ApiConstructorFileName)

	// if bootstrap for server doesn't exist - adding it
	if apiFile == nil {
		apiFile = &folder.Folder{
			Name:    patterns.ApiConstructorFileName,
			Content: patterns.APISetupFile,
		}
		p.GetFolder().ForceAddWithPath([]string{patterns.CmdFolder, p.GetName(), patterns.BootStrapFolder}, apiFile)
	}

	const (
		// key lines for starting and stopping servers
		apiEntryPointStopFunc = `stopFunc := `
		apiEntryPointCall     = `bootstrap.ApiEntryPoint`
		apiEntryPointArgs     = `(ctx, cfg)`

		apiEntryPointStop         = `stopFunc(context.Background())`
		apiEntryPointStopFuncCall = `
	err = %s
	if err != nil {
		logrus.Fatal(err)
	}
`
	)

	var insertBeforeEnd []byte
	var insertAfterEnd []byte

	// add to main file api entry call if not exists
	if bytes.Index(projMainFile.Content, []byte(apiEntryPointCall)) == -1 {
		insertBeforeEnd = bytes.Join([][]byte{
			insertBeforeEnd,
			[]byte(apiEntryPointStopFunc),
			[]byte(apiEntryPointCall),
			[]byte(apiEntryPointArgs),
		}, []byte{})
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
			[]byte(fmt.Sprintf(apiEntryPointStopFuncCall, apiEntryPointStop))...)
	}
	if len(insertAfterEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertAfterEnd, endFuncIdx)
	}

	{
		// add import on boostrap if doesn't exists
		importBootstrap := []byte("\"" + p.GetProjectModName() + "/cmd/" + p.GetName() + "/bootstrap\"\n")
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
	apiFile := p.GetFolder().GetByPath(patterns.CmdFolder, p.GetName(), patterns.BootStrapFolder, patterns.ApiConstructorFileName)
	if apiFile == nil {
		apiFile = &folder.Folder{
			Name:    patterns.ApiConstructorFileName,
			Content: patterns.APISetupFile,
		}
		p.GetFolder().ForceAddWithPath([]string{patterns.CmdFolder, p.GetName(), patterns.BootStrapFolder}, apiFile)
	}

	err := insertMissingAPI(p, serverFolders, apiFile)
	if err != nil {
		return errors.Wrap(err, "error inserting missing api")
	}

	apiMgr := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.TransportFolder, patterns.ApiManagerFileName)
	if apiMgr == nil {
		serverFolders = append(serverFolders, &folder.Folder{
			Name:    patterns.ApiManagerFileName,
			Content: patterns.ServerManagerPatternFile,
		})
	}

	p.GetFolder().AddWithPath([]string{patterns.InternalFolder, patterns.TransportFolder}, serverFolders...)

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

		newAPIImportInsert = append(newAPIImportInsert, []byte("\n\t\""+p.GetProjectModName()+"/internal/transport/"+serv.Name+"\"")...)
		newAPIInsert = append(newAPIInsert, []byte("mngr.AddServer("+serv.Name+".NewServer(cfg))\n\t")...)

		if strings.Contains(serv.Name, _const.GRPCServer) {
			grpcServers = append(grpcServers, serv.Name)
		}
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
			protoFolder := p.GetFolder().GetByPath(patterns.PkgFolder, patterns.ProtoFolder, serverName)
			if protoFolder == nil {

				exampleFile := bytes.ReplaceAll(
					patterns.GrpcProtoExampleFile,
					[]byte(patterns.ImportProjectNamePatternSnakeCase),
					[]byte(cases.KebabToSnake(serverName)))

				exampleFile = bytes.ReplaceAll(exampleFile,
					[]byte("grpc_realisation"),
					[]byte(serverName),
				)
				p.GetFolder().AddWithPath(
					[]string{patterns.PkgFolder, patterns.ProtoFolder},
					&folder.Folder{
						Name: serverName,
						Inner: []*folder.Folder{
							{
								Name:    serverName + patterns.ProtoFileExtension,
								Content: exampleFile,
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
	rscliMkF := p.GetFolder().GetByPath(patterns.RsCliMkFileName)
	if rscliMkF == nil {
		return ErrNoMakeFile
	}

	protocMkfile := bytes.Join([][]byte{
		patterns.GRPCSection,
		patterns.GRPCUtilityInstallGoProtocHeader,
		patterns.GRPCGenerateGoCodeWithDependencies,
		patterns.GRPCInstallProtocViaGolangEnvOSBased,
		patterns.GRPCInstallGolangProtoc,
		patterns.GRPCGenerateGoCode,
		patterns.GRPCGatewayDependency,
	},
		[]byte{})

	idxStart := bytes.Index(rscliMkF.Content, patterns.GRPCUtilityInstallGoProtocHeader)
	if idxStart == -1 {
		rscliMkF.Content = append(rscliMkF.Content, protocMkfile...)
	} else {
		idxEnd := idxStart + bytes.Index(rscliMkF.Content[idxStart:], patterns.SectionSeparator)
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
