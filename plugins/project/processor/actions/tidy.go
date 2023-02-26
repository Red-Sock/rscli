package actions

import (
	"bytes"
	"github.com/Red-Sock/rscli/internal/utils/slices"

	"github.com/pkg/errors"

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
	goFuncWord          = "go func() {\n"
	transportNewManager = "transport.NewManager()"
)

const (
	cmdFolder              = "cmd"
	mainFileName           = "main.go"
	apiConstructorFileName = "http.go"

	internalFolder     = "internal"
	transportFolder    = "transport"
	apiManagerFileName = "manager.go"
)

func PrepareMainFile(p interfaces.Project) error {
	cfg := p.GetConfig()

	projectStartupFolder := p.GetFolder().GetByPath(cmdFolder, p.GetName())
	projMainFile := projectStartupFolder.GetByPath(mainFileName)
	if projMainFile == nil {
		return ErrNoMainFile
	}

	err := tidyAPI(p, cfg, projMainFile, projectStartupFolder)
	if err != nil {
		return errors.Wrap(err, "error tiding API")
	}

	return nil
}

func tidyAPI(p interfaces.Project, cfg interfaces.Config, projMainFile, projectStartupFolder *folder.Folder) error {
	serverFolders, err := cfg.ExtractServerOptions()
	if err != nil {
		return err
	}

	if serverFolders == nil {
		return nil
	}

	httpFile := p.GetFolder().GetByPath(cmdFolder, p.GetName(), apiConstructorFileName)

	httpFile = tidyMainForAPI(httpFile, projMainFile)
	projectStartupFolder.Add(httpFile)

	tidyAPIFile(p, serverFolders, httpFile)

	return nil
}

func tidyMainForAPI(httpFile *folder.Folder, projMainFile *folder.Folder) *folder.Folder {

	if httpFile == nil {
		httpFile = &folder.Folder{
			Name:    apiConstructorFileName,
			Content: patterns.APISetupFile,
		}
	}

	const (
		apiEntryPointCall = "stopFunc := apiEntryPoint(ctx, cfg)\n\n\t"
		apiEntryPointStop = "\n\n\terr = stopFunc(context.Background())\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}"
	)

	var insertBeforeEnd []byte
	var insertAfterEnd []byte

	if bytes.Index(projMainFile.Content, []byte(apiEntryPointCall)) == -1 {
		insertBeforeEnd = append(insertBeforeEnd, []byte(apiEntryPointCall)...)
	}
	if bytes.Index(projMainFile.Content, []byte(apiEntryPointStop)) == -1 {
		insertAfterEnd = append(insertAfterEnd, []byte(apiEntryPointStop)...)
	}

	endFuncIdx := bytes.Index(projMainFile.Content, []byte(waitingForTheEndFunc))
	if len(insertBeforeEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertBeforeEnd, endFuncIdx)
	}
	if len(insertAfterEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertAfterEnd, endFuncIdx+len(insertBeforeEnd)+len(waitingForTheEndFunc)+1)
	}

	return httpFile
}

func tidyAPIFile(p interfaces.Project, serverFolders []*folder.Folder, httpFile *folder.Folder) {
	var apisBytes []byte
	{
		goFuncWordBytes := []byte(goFuncWord)
		startIdx := bytes.Index(httpFile.Content, []byte(transportNewManager)) + len(transportNewManager) + 2

		endIdx := bytes.Index(httpFile.Content, goFuncWordBytes)

		apisBytes = httpFile.Content[startIdx:endIdx]

	}

	var newAPIInsert []byte
	var newAPIImportInsert []byte
	for _, serv := range serverFolders {
		if bytes.Contains(apisBytes, []byte(serv.Name)) {
			continue
		}
		newAPIImportInsert = append(newAPIImportInsert, []byte("\n\t"+serv.Name+" \""+p.GetName()+"/internal/transport/"+serv.Name+"\"\n")...)
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

	apiMgr := p.GetFolder().GetByPath(internalFolder, apiManagerFileName)
	if apiMgr == nil {
		serverFolders = append(serverFolders, &folder.Folder{
			Name:    apiManagerFileName,
			Content: patterns.ServerManagerPattern,
		})
	}

	p.GetFolder().AddWithPath([]string{internalFolder, transportFolder}, serverFolders...)
}
