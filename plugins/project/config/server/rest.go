package server

const DefaultRestPort = 8080

type Rest struct {
	name string
	Port uint16
}

func (r *Rest) GetName() string {
	return r.name
}

func (r *Rest) GetPort() uint16 {
	if r.Port != 0 {
		return r.Port
	}

	return DefaultRestPort
}
