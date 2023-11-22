package grpc_api

import (
	"context"
	"time"

	"financial-microservice/pkg/api/proj_name_api"
)

type Pinger struct {
	calcFunc func(time time.Time) (diff int32)
	proj_name_api.UnimplementedProjNameAPIServer
}

func (p *Pinger) Version(_ context.Context, req *proj_name_api.PingRequest) (rsp *proj_name_api.PingResponse, err error) {

	diff := time.Since(req.ClientTimestamp.AsTime())

	return &proj_name_api.PingResponse{Took: uint32(diff)}, nil
}
