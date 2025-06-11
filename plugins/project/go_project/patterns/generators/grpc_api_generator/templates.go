package grpc_api_generator

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/api.proto.pattern
	basicProto            string
	basicApiProtoTemplate *template.Template
)

func init() {
	basicApiProtoTemplate = template.Must(
		template.New("basic_proto").
			Parse(basicProto))
}
