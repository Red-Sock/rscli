package server

const DefaultGrpcPort = 50051

type GRPC struct {
	name string
	Port uint16
}

func (r *GRPC) GetName() string {
	return r.name
}

func (r *GRPC) GetPort() uint16 {
	if r.Port != 0 {
		return r.Port
	}

	return DefaultGrpcPort
}
