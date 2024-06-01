package perun

import (
	"context"

	"github.com/godverv/matreshka"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"proj_name/pkg/example_api"
)

type Implementation struct {
	example_api.UnimplementedProjNameAPIServer

	version string
}

func New(cfg matreshka.Config) *Implementation {
	return &Implementation{
		version: cfg.AppInfo().Version,
	}
}

func (impl *Implementation) Register(server grpc.ServiceRegistrar) {
	example_api.RegisterProjNameAPIServer(server, impl)
}

func (impl *Implementation) RegisterGw(ctx context.Context, mux *runtime.ServeMux, addr string) error {
	return example_api.RegisterProjNameAPIHandlerFromEndpoint(
		ctx,
		mux,
		addr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		})
}
