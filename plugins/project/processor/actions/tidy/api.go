package tidy

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/internal/utils/slices"
	"github.com/Red-Sock/rscli/pkg/folder"
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

	httpFile := insertApiSetupIfNotExists(p, projMainFile)

	tidyAPIFile(p, serverFolders, httpFile)

	return nil
}
func insertApiSetupIfNotExists(p interfaces.Project, projMainFile *folder.Folder) *folder.Folder {
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

	if bytes.Index(projMainFile.Content, []byte(apiEntryPointCall)) == -1 {

		insertBeforeEnd = append(insertBeforeEnd, []byte(apiEntryPointStopFunc)...)
		insertBeforeEnd = append(insertBeforeEnd, []byte(apiEntryPointCall)...)
		insertBeforeEnd = append(insertBeforeEnd, []byte(apiEntryPointArgs)...)

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

	return apiFile
}
func tidyAPIFile(p interfaces.Project, serverFolders []*folder.Folder, httpFile *folder.Folder) {
	insertMissingAPI(p, serverFolders, httpFile)

	apiMgr := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.TransportFolder, patterns.ApiManagerFileName)
	if apiMgr == nil {
		serverFolders = append(serverFolders, &folder.Folder{
			Name:    patterns.ApiManagerFileName,
			Content: patterns.ServerManagerPattern,
		})
	}

	p.GetFolder().AddWithPath([]string{patterns.InternalFolder, patterns.TransportFolder}, serverFolders...)
}

func insertMissingAPI(p interfaces.Project, serverFolders []*folder.Folder, httpFile *folder.Folder) {
	serverInit := extractInitApi(httpFile.Content)

	var newAPIInsert []byte
	var newAPIImportInsert []byte
	for _, serv := range serverFolders {
		if bytes.Contains(serverInit, []byte(serv.Name)) {
			continue
		}

		newAPIImportInsert = append(newAPIImportInsert, []byte("\n\t\""+p.GetProjectModName()+"/internal/transport/"+serv.Name+"\"")...)
		newAPIInsert = append(newAPIInsert, []byte("mngr.AddServer("+serv.Name+".NewServer(cfg))\n\t")...)
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
}

func extractInitApi(httpFile []byte) (out []byte) {
	goFuncWordBytes := []byte(goFuncWord)
	// indexes between creation of transport manager
	// and starting it in goroutine
	startIdx := bytes.Index(httpFile, []byte(transportNewManager)) + len(transportNewManager) + 2
	endIdx := bytes.Index(httpFile, goFuncWordBytes)

	out = make([]byte, endIdx-startIdx)
	copy(out, httpFile[startIdx:endIdx])

	return out
}
