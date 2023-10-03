package grpc_api

import (
	"context"
	"time"

	pb "financial-microservice/pkg/grpc-realisation"
)

type Pinger struct {
	calcFunc func(time time.Time) (diff int32)
	pb.UnimplementedFinancialAPIServer
}

func (p *Pinger) Version(_ context.Context, req *pb.PingRequest) (rsp *pb.PingResponse, err error) {

	diff := time.Since(req.ClientTimestamp.AsTime())

	return &pb.PingResponse{Took: uint32(diff)}, nil
}
