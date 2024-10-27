package app

import (
	"github.com/Red-Sock/link_grpc/internal/clients/grpc"
	errors "github.com/Red-Sock/trace-errors"
)

func (a *App) InitDataSources() (err error) {
	a.GrpcHelloWorld, err = grpc.NewHelloWorldAPIClient(a.Cfg.DataSources.GrpcHelloWorld)
	if err != nil {
		return errors.Wrap(err, "error during grpc client initialization")
	}

	return nil
}
