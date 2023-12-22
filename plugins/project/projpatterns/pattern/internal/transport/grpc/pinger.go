package grpc

import (
	"context"
	"time"

	api "proj_name/pkg/api/example_api"
)

func (p *ExampleApi) Version(_ context.Context, req *api.PingRequest) (rsp *api.PingResponse, err error) {
	diff := time.Since(req.ClientTimestamp.AsTime())

	return &api.PingResponse{Took: uint32(diff)}, nil
}
