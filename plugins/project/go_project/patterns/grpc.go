package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

// Proto contract
var (
	//go:embed pattern_c/api/grpc/api.proto
	protoServer []byte
	ProtoServer = &folder.Folder{
		Name:    "api.proto",
		Content: protoServer,
	}
)

// Dependencies and generator
var (
	//go:embed pattern_c/easyp.yaml
	easyp []byte
	EasyP = &folder.Folder{
		Name:    "easyp.yaml",
		Content: easyp,
	}

	//go:embed pattern_c/grpc.mk
	GrpcServerGenMK []byte
)
