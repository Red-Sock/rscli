package grpc

import (
	errors "github.com/Red-Sock/trace-errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connect(connString string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	dial, err := grpc.NewClient(connString, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "error dialing")
	}

	return dial, nil
}
