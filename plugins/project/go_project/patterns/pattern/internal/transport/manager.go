package transport

import (
	"context"
	"net"
	"net/http"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

type ServersManager struct {
	ctx context.Context

	mux cmux.CMux

	grpcServer
	httpServer
}

func NewServerManager(ctx context.Context, port string) (*ServersManager, error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil, errors.Wrap(err, "error opening listener")
	}

	mainMux := cmux.New(listener)
	httpMux := http.NewServeMux()

	s := &ServersManager{
		ctx: ctx,
		mux: mainMux,

		grpcServer: newGrpcServer(ctx, mainMux.Match(cmux.HTTP2()), httpMux),
		httpServer: newHttpServer(mainMux.Match(cmux.Any()), httpMux),
	}

	return s, nil
}

func (m *ServersManager) Start() error {
	logrus.Info("Starting server at http://0.0.0.0" + m.grpcServer.listener.Addr().String()[4:])
	errGroup, ctx := errgroup.WithContext(context.Background())

	errGroup.Go(m.mux.Serve)
	errGroup.Go(m.grpcServer.start)
	errGroup.Go(m.httpServer.start)

	errC := make(chan error, 1)

	select {
	case <-ctx.Done():
		return nil
	case errC <- errGroup.Wait():
		err := <-errC
		return errors.Wrap(err)
	}
}

func (m *ServersManager) Stop() error {
	eg, _ := errgroup.WithContext(m.ctx)

	eg.Go(m.grpcServer.stop)
	eg.Go(m.httpServer.stop)
	eg.Go(func() error { m.mux.Close(); return nil })

	err := eg.Wait()
	if err != nil {
		return errors.Wrap(err, "error stopping server")
	}

	return nil
}
