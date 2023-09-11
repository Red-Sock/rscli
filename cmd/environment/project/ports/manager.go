package ports

import (
	"sync"
)

type PortManager struct {
	ports map[uint16]string
	m     sync.Mutex
}

func NewPortManager() *PortManager {
	return &PortManager{
		ports: map[uint16]string{},
	}
}

func (p *PortManager) GetPort(port uint16) (string, bool) {
	res, ok := p.ports[port]
	return res, ok
}

func (p *PortManager) GetNextPort(in uint16, key string) uint16 {
	p.m.Lock()
	defer p.m.Unlock()

	for {
		// if such port already exists - increment it
		if portName, ok := p.ports[in]; !ok {
			p.ports[in] = key
			return in
		} else {
			if portName == key {
				return in
			}
		}
		in++
	}
}

func (p *PortManager) SaveIfNotExist(port uint16, name string) (conflict string) {
	p.m.Lock()
	defer p.m.Unlock()

	if existingName, ok := p.ports[port]; ok {
		return existingName

	}

	p.ports[port] = name
	return ""
}
