package server

type GRPC struct {
	name string
	Port uint16
}

func (r *GRPC) GetName() string {
	return r.name
}

func (r *GRPC) GetPort() uint16 {
	return r.Port
}
