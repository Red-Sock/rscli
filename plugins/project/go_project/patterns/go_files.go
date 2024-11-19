package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

// Constants naming: Purpose+Type (File)

const (
	GithubFolder    = ".github"
	WorkflowsFolder = "workflows"

	CmdFolder     = "cmd"
	ServiceFolder = "service"

	InternalFolder             = "internal"
	AppFolder                  = "app"
	AppFileName                = "app.go"
	AppInitServerFileName      = "server.go"
	AppInitDataSourcesFileName = "data_sources.go"
	AppConfigFileName          = "config.go"
	AppCustomFileName          = "custom.go"

	ConnFileName = "conn.go"

	TransportFolder = "transport"
	WebFolder       = "web"
	DistFolder      = "dist"
	IndexHtmlFile   = "index.html"
	AboutFolder     = "about"

	HandlersFolderName = "handlers"
	VersionFolderName  = "version"

	ConfigsFolder      = "config"
	ConfigTemplateYaml = "config_template.yaml"

	ConfigLoadFileName        = "load.go"
	ConfigDataSourcesFileName = "data_sources.go"
	ConfigEnvironmentFileName = "environment.go"
	ConfigServersFileName     = "servers.go"

	GoMod = "go.mod"

	ExampleFile = ".example"
)

// Basic files
var (
	//go:embed pattern_c/cmd/service/main.go.pattern
	mainFile []byte
	MainFile = &folder.Folder{
		Name:    "main.go",
		Content: mainFile,
	}
)
