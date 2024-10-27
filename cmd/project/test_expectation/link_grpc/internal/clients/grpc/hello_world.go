package grpc

import (
	errors "github.com/Red-Sock/trace-errors"
	pb "github.com/godverv/hello_world/pkg/api"
	"github.com/godverv/matreshka/resources"
	"google.golang.org/grpc"
)

type HelloWorldAPIClient pb.HelloWorldAPIClient

func NewHelloWorldAPIClient(grpcConn *resources.GRPC, opts ...grpc.DialOption) (HelloWorldAPIClient, error) {
	conn, err := connect(grpcConn.ConnectionString, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "error crating grpc client")
	}

	return pb.NewHelloWorldAPIClient(conn), nil
}
