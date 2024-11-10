package app_struct_generators

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns/generators"
)

type AppFileGenArgs struct {
	Imports    map[string]string
	AppContent []AppContent
	Starters   []AppStarter
}

type AppContent struct {
	Comment              string
	Fields               []generators.KeyValue
	InitFunc             string
	InitFuncErrorMessage string
	Imports              map[string]string
}

type AppStarter struct {
	FieldName string
	StartCall string
	StopCall  string
}

func GenerateAppFiles(p project.IProject) (map[string][]byte, error) {
	initAppArgs := AppFileGenArgs{
		Imports: make(map[string]string),
	}

	cfg := p.GetConfig()

	out := make(map[string][]byte)

	// init data sources
	if len(cfg.DataSources) != 0 {
		initDataSourcesArgs, initDataSourcesFile, err := generateDataSourceInitFileAndArgs(cfg.DataSources)
		if err != nil {
			return nil, errors.Wrap(err, "error generation data source init file")
		}

		out[patterns.AppInitDataSourcesFileName] = initDataSourcesFile
		if initDataSourcesArgs != nil {
			initAppArgs.AppContent = append(initAppArgs.AppContent, *initDataSourcesArgs)
		}

	}
	// init server
	if len(cfg.Servers) != 0 {
		initServerArgs, initServerFile, err := generateServerInitFileAndArgs(cfg.Servers)
		if err != nil {
			return nil, errors.Wrap(err, "error generating server init file")
		}

		out[patterns.AppInitServerFileName] = initServerFile
		initAppArgs.addAppContent(
			"/* Servers managers */",
			"error during server initialization",
			initServerArgs)

		initAppArgs.Starters = append(initAppArgs.Starters, AppStarter{
			FieldName: initServerArgs.ServerName,
		})
	}

	for _, ac := range initAppArgs.AppContent {
		for importPath, dependencyAlias := range ac.Imports {
			appAlias, ok := initAppArgs.Imports[importPath]
			if ok && appAlias != dependencyAlias {
				return nil, errors.New("Fatal error: app already imported package " +
					importPath + " with alias " + appAlias +
					". But dependency requires this package to be imported as " +
					dependencyAlias)
			}
			initAppArgs.Imports[importPath] = dependencyAlias
		}
	}

	mainAppFile := &rw.RW{}
	err := appTemplate.Execute(mainAppFile, initAppArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error generating app file")
	}

	out[patterns.AppFileName] = mainAppFile.Bytes()
	out[patterns.AppConfigFileName] = appConfigPattern

	if p.GetFolder().GetByPath(patterns.InternalFolder, patterns.AppFolder, patterns.AppCustomFileName) == nil {
		out[patterns.AppCustomFileName] = customPattern
	}

	return out, nil
}

func (a *AppFileGenArgs) addAppContent(comment, errMsg string, args InitDepFuncGenArgs) {
	serverAppContent := AppContent{
		Comment:              comment,
		Fields:               make([]generators.KeyValue, 0, len(args.Functions)),
		InitFunc:             args.InitFunctionName,
		InitFuncErrorMessage: errMsg,
	}
	for _, serverInitFunc := range args.Functions {
		serverAppContent.Fields = append(serverAppContent.Fields,
			generators.KeyValue{
				Key:   serverInitFunc.ResultName,
				Value: serverInitFunc.ResultType,
			})
	}

	for importPath, importAlias := range args.Imports {
		a.Imports[importPath] = importAlias
	}
	a.AppContent = append(a.AppContent, serverAppContent)
}
