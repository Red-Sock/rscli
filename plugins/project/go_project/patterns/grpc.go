package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

const (
	GrpcImplFolder = "grpc_impl"
)

// Proto contract
var (
	//go:embed pattern_c/api/grpc/api.proto
	protoContract []byte
	ProtoContract = &folder.Folder{
		Name:    "api.proto",
		Content: protoContract,
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

var (
	//go:embed pattern_c/internal/transport/grpc/example_api_impl/impl.go.pattern
	grpcImpl []byte
	GrpcImpl = &folder.Folder{
		Name:    "impl.go",
		Content: grpcImpl,
	}
)
