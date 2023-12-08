// Code generated by RedSock CLI

package grpc

import (
	"context"
	"net"

	"github.com/godverv/matreshka/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"proj_name/internal/config"
)

type Server struct {
	srv *grpc.Server

	networkType string
	address     string
}

func NewServer(cfg config.Config, server *api.GRPC) *Server {
	srv := grpc.NewServer()

	// Register your servers here

	return &Server{
		srv:         srv,
		networkType: "tcp",
		address:     "0.0.0.0:" + server.GetPortStr(),
	}
}

func (s *Server) Start(_ context.Context) error {
	lis, err := net.Listen(s.networkType, s.address)
	if err != nil {
		return errors.Wrapf(err, "error when tried to listen for %s, %s", s.networkType, s.address)
	}

	go func() {
		logrus.Infof("Starting GRPC Server at %s (%s)", s.address, s.networkType)
		err = s.srv.Serve(lis)
		if err != nil {
			logrus.Errorf("error serving grpc: %s", err)
		} else {
			logrus.Infof("GRPC Server at %s is Stopped", s.address)
		}
	}()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	logrus.Infof("Stopping GRPC server at %s", s.address)
	s.srv.GracefulStop()
	logrus.Infof("GRPC server at %s is stopped", s.address)
	return nil
}
