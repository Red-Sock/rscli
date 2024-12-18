package transport

import (
	"context"
	"net"
	"net/http"

	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcImpl interface {
	Register(srv grpc.ServiceRegistrar)
}

type GrpcWithGateway interface {
	Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (rootRoute string, handler http.Handler)
}

type grpcServer struct {
	ctx      context.Context
	listener net.Listener

	gatewayMux *http.ServeMux

	opts            []grpc.ServerOption
	implementations []GrpcImpl

	// AvailableAfter start is called
	stopCall func()
}

func newGrpcServer(
	ctx context.Context,
	listener net.Listener,
	gatewayMux *http.ServeMux) grpcServer {
	return grpcServer{
		ctx:        ctx,
		listener:   listener,
		stopCall:   func() {},
		gatewayMux: gatewayMux,
	}
}

func (s *grpcServer) start() error {
	server := grpc.NewServer(s.opts...)

	for _, impl := range s.implementations {
		impl.Register(server)
	}

	err := server.Serve(s.listener)
	if err != nil {
		if !rerrors.Is(err, http.ErrServerClosed) {
			return rerrors.Wrap(err, "error serving grpc server")
		}
	}

	s.stopCall = server.GracefulStop

	return nil
}

func (s *grpcServer) stop() error {
	s.stopCall()
	return nil
}

func (s *grpcServer) AddImplementation(grpcImpl GrpcImpl, opts ...grpc.ServerOption) {
	s.implementations = append(s.implementations, grpcImpl)

	grpcWithGateway, ok := grpcImpl.(GrpcWithGateway)
	if ok {
		s.gatewayMux.Handle(grpcWithGateway.Gateway(s.ctx,
			s.listener.Addr().String(),
			grpc.WithTransportCredentials(insecure.NewCredentials())))
	}

	s.opts = append(s.opts, opts...)
}
