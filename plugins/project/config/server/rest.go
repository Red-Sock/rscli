package server

type Rest struct {
	name string
	Port uint16
}

func (r *Rest) GetName() string {
	return r.name
}

func (r *Rest) GetPort() uint16 {
	return r.Port
}
