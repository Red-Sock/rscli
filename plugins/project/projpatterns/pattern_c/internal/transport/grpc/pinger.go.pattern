package grpc

import (
	"context"
	"time"

	"proj_name/pkg/api"
)

type Pinger struct {
	calcFunc func(time time.Time) (diff int32)
	api.UnimplementedProjNameAPIServer
}

func (p *Pinger) Version(_ context.Context, req *api.PingRequest) (rsp *api.PingResponse, err error) {
	diff := time.Since(req.ClientTimestamp.AsTime())

	return &api.PingResponse{Took: uint32(diff)}, nil
}
