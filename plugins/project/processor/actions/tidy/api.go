package tidy

import (
	"bytes"
	"strings"

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
	goFuncWord          = "go func() {"
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
	serverFolders, err := cfg.GetServerFolders()
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
	wfteBytes := []byte(waitingForTheEndFunc)

	endFuncIdx := bytes.Index(projMainFile.Content, wfteBytes)
	endFuncIdx = bytes.LastIndex(projMainFile.Content[:endFuncIdx], []byte("\n"))
	if len(insertBeforeEnd) != 0 {
		projMainFile.Content = slices.InsertSlice(projMainFile.Content, insertBeforeEnd, endFuncIdx)
		endFuncIdx = bytes.Index(projMainFile.Content, wfteBytes) + len(wfteBytes)
		endFuncIdx = endFuncIdx + bytes.Index(projMainFile.Content[endFuncIdx:], []byte("\n")) + 1
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
	removeExtraAPI(p, serverFolders, httpFile)

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

		newAPIImportInsert = append(newAPIImportInsert, []byte("\n\t\""+p.GetName()+"/internal/transport/"+serv.Name+"\"")...)
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
func removeExtraAPI(p interfaces.Project, serverFolders []*folder.Folder, httpFile *folder.Folder) {
	var aliasesInFile []string
	{
		goFuncWordBytes := []byte(goFuncWord)
		startIdx := bytes.Index(httpFile.Content, []byte(transportNewManager)) + len(transportNewManager) + 2

		endIdx := bytes.Index(httpFile.Content, goFuncWordBytes)

		splitedNames := strings.Split(string(httpFile.Content[startIdx:endIdx]), "\n")

		replacer := strings.NewReplacer(
			"\n", "",
			"\t", "",
			"mngr.AddServer(", "",
		)

		for _, item := range splitedNames {
			item = replacer.Replace(item)
			if item != "" {
				if idx := strings.Index(item, ".NewServer("); idx != -1 {
					item = item[:idx]
					aliasesInFile = append(aliasesInFile, item)
				}

			}
		}
	}

	aliasesFromConfig := make(map[string]struct{}, len(serverFolders))
	{
		for _, serv := range serverFolders {
			aliasesFromConfig[serv.Name] = struct{}{}
		}
	}

	for _, aliasInFile := range aliasesInFile {

		if _, ok := aliasesFromConfig[aliasInFile]; !ok {
			abbB := []byte(aliasInFile)
			idx := bytes.Index(httpFile.Content, abbB)
			for idx != -1 {
				startIdx := bytes.LastIndexByte(httpFile.Content[:idx], '\n') + 1
				endIdx := idx + bytes.IndexByte(httpFile.Content[idx:], '\n')
				httpFile.Content = slices.RemovePart(httpFile.Content, startIdx, endIdx)

				idx = bytes.Index(httpFile.Content, abbB)
			}
		}
	}

	transports := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.TransportFolder)
	if transports == nil {
		return
	}
	for idx := range transports.Inner {
		if _, ok := aliasesFromConfig[transports.Inner[idx].Name]; !ok && len(transports.Inner[idx].Content) == 0 {
			transports.Inner[idx].Delete()
		}
	}
}