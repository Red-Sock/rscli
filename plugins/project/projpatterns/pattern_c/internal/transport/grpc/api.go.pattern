package grpc

import (
	"time"

	api "proj_name/pkg/api/example_api"
)

type ExampleApi struct {
	calcFunc func(time time.Time) (diff int32)
	api.UnimplementedProjNameAPIServer
}
