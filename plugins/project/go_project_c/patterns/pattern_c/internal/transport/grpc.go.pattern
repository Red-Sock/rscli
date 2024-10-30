package transport

import (
	"context"
	"net"
	"net/http"

	errors "github.com/Red-Sock/trace-errors"
	"google.golang.org/grpc"
)

type GrpcImpl interface {
	Register(srv *grpc.Server)
}

type GrpcWithGateway interface {
	Gateway(ctx context.Context) (rootRoute string, handler http.Handler)
}

type grpcServer struct {
	ctx    context.Context
	server *grpc.Server

	listener net.Listener

	gatewayMux *http.ServeMux
}

func newGrpcServer(ctx context.Context, listener net.Listener, gatewayMux *http.ServeMux) grpcServer {
	return grpcServer{
		ctx:        ctx,
		server:     grpc.NewServer(),
		listener:   listener,
		gatewayMux: gatewayMux,
	}
}

func (s *grpcServer) start() error {
	err := s.server.Serve(s.listener)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "error serving grpc server")
		}
	}

	return nil
}

func (s *grpcServer) stop() error {
	s.server.GracefulStop()
	return nil
}

func (s *grpcServer) AddGrpcServer(grpcImpl GrpcImpl) {
	grpcImpl.Register(s.server)

	grpcWithGateway, ok := grpcImpl.(GrpcWithGateway)
	if ok {
		s.gatewayMux.Handle(grpcWithGateway.Gateway(s.ctx))
	}
}
