package tidy

import (
	"bytes"
	"github.com/Red-Sock/rscli/internal/utils/slices"
	"strings"

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

func Api(p interfaces.Project) error {
	cfg := p.GetConfig()

	projMainFile := p.GetFolder().GetByPath(patterns.CmdFolder, p.GetName(), patterns.MainFileName)
	if projMainFile == nil {
		return ErrNoMainFile
	}

	err := tidyAPI(p, cfg, projMainFile)
	if err != nil {
		return errors.Wrap(err, "error tiding API")
	}

	return nil
}

func tidyAPI(p interfaces.Project, cfg interfaces.Config, projMainFile *folder.Folder) error {
	serverFolders, err := cfg.ExtractServerOptions()
	if err != nil {
		return err
	}

	if serverFolders == nil {
		return nil
	}

	httpFile := insertApiSetupIfNotExists(p, projMainFile)

	tidyAPIFile(p, serverFolders, httpFile)

	return nil
}
func insertApiSetupIfNotExists(p interfaces.Project, projMainFile *folder.Folder) *folder.Folder {
	httpFile := p.GetFolder().GetByPath(patterns.CmdFolder, p.GetName(), patterns.ApiConstructorFileName)

	if httpFile == nil {
		httpFile = &folder.Folder{
			Name:    patterns.ApiConstructorFileName,
			Content: patterns.APISetupFile,
		}
		p.GetFolder().GetByPath(patterns.CmdFolder, p.GetName()).Add(httpFile)
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
	endFuncIdx := bytes.Index(projMainFile.Content, []byte(waitingForTheEndFunc))
	if len(insertBeforeEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertBeforeEnd, endFuncIdx)
		endFuncIdx += len(insertBeforeEnd) + len(waitingForTheEndFunc) + 1
	}

	if bytes.Index(projMainFile.Content, []byte(apiEntryPointStop)) == -1 {
		insertAfterEnd = append(insertAfterEnd, []byte(apiEntryPointStop)...)
	}
	if len(insertAfterEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertAfterEnd, endFuncIdx)
	}

	return httpFile
}
func tidyAPIFile(p interfaces.Project, serverFolders []*folder.Folder, httpFile *folder.Folder) {
	insertMissingAPI(p, serverFolders, httpFile)
	removeExtraAPI(serverFolders, httpFile)

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
}
func removeExtraAPI(serverFolders []*folder.Folder, httpFile *folder.Folder) {
	var aliasesInFile []string
	{
		goFuncWordBytes := []byte(goFuncWord)
		startIdx := bytes.Index(httpFile.Content, []byte(transportNewManager)) + len(transportNewManager) + 2

		endIdx := bytes.Index(httpFile.Content, goFuncWordBytes)

		splitedNames := strings.Split(string(httpFile.Content[startIdx:endIdx]), "\n")

		replacer := strings.NewReplacer(
			"\n", "",
			"\t", "",
		)

		for _, item := range splitedNames {
			item = replacer.Replace(item)
			if item != "" {
				aliasesInFile = append(aliasesInFile, item)
			}
		}
	}

	aliasesFromConfig := make([]string, len(serverFolders))
	{
		for idx, serv := range serverFolders {
			aliasesFromConfig[idx] = serv.Name
		}
	}

	for _, aliasInFile := range aliasesInFile {

		aliasExistsInConfig := false
		for _, aliasFromConfig := range aliasesFromConfig {
			if strings.Contains(aliasInFile, aliasFromConfig) {
				aliasExistsInConfig = true
				break
			}
		}

		if !aliasExistsInConfig {
			abbB := []byte(aliasInFile)
			idx := bytes.Index(httpFile.Content, abbB)
			for idx != -1 {
				startIdx := bytes.LastIndexByte(httpFile.Content[:idx], '\n') + 1
				endIdx := idx + bytes.IndexByte(httpFile.Content[idx:], '\n') + 1
				httpFile.Content = slices.RemovePart(httpFile.Content, startIdx, endIdx)

				idx = bytes.Index(httpFile.Content, abbB)
			}
		}
	}
}
