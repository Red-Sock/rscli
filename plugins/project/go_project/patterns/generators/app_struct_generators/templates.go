package app_struct_generators

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/app.go.pattern
	appPattern  string
	appTemplate *template.Template

	//go:embed templates/config.go.pattern
	appConfigPattern []byte
	//go:embed templates/custom.go.pattern
	customPattern []byte

	//go:embed templates/init_structure.go.pattern
	initAppStructFuncPattern  string
	initAppStructFuncTemplate *template.Template

	//go:embed templates/init_server.go.pattern
	initServerPattern  string
	initServerTemplate *template.Template
)

func init() {
	appTemplate = template.Must(
		template.New("app").
			Parse(appPattern))

	initAppStructFuncTemplate = template.Must(
		template.New("init_app_struct_func").
			Parse(initAppStructFuncPattern))

	initServerTemplate = template.Must(
		template.New("init_server_struct_func").
			Parse(initServerPattern))
}
