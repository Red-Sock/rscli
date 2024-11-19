package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

// Server files
var (
	//go:embed pattern_c/internal/transport/manager.go.pattern
	serverManagerPatternFile []byte
	ServerManager            = &folder.Folder{
		Name:    "manager.go",
		Content: serverManagerPatternFile,
	}

	//go:embed pattern_c/internal/transport/grpc.go.pattern
	grpcServerManagerFile []byte
	GrpcServerManager     = &folder.Folder{
		Name:    "grpc.go",
		Content: grpcServerManagerFile,
	}

	//go:embed pattern_c/internal/transport/http.go.pattern
	httpServerManagerFile []byte
	HttpServerManager     = &folder.Folder{
		Name:    "http.go",
		Content: httpServerManagerFile,
	}

	//go:embed pattern_c/internal/transport/telegram/listener.go.pattern
	tgServFile []byte
	TgServFile = &folder.Folder{
		Name:    "listener.go",
		Content: tgServFile,
	}
	//go:embed pattern_c/internal/transport/telegram/version/handler.go.pattern
	tgHandlerExampleFile []byte
	TgHandlerExampleFile = &folder.Folder{
		Name:    "handler.go",
		Content: tgHandlerExampleFile,
	}
)
