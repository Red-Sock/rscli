package perun

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"proj_name/internal/config"
	"proj_name/pkg/example_api"
)

type Impl struct {
	example_api.UnimplementedProjNameAPIServer

	version string
}

func New(cfg config.Config) *Impl {
	return &Impl{
		version: cfg.AppInfo.Version,
	}
}

func (impl *Impl) Register(server grpc.ServiceRegistrar) {
	example_api.RegisterProjNameAPIServer(server, impl)
}

func (impl *Impl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (route string, handler http.Handler) {
	gwHttpMux := runtime.NewServeMux()

	err := example_api.RegisterProjNameAPIHandlerFromEndpoint(
		ctx,
		gwHttpMux,
		endpoint,
		opts,
	)
	if err != nil {
		logrus.Errorf("error registering grpc2http handler: %s", err)
	}

	return "/api/", gwHttpMux
}
