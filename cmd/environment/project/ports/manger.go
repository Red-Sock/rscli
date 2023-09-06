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
func (p *PortManager) GetNextPort(in uint16, projName string) uint16 {
	p.m.Lock()
	defer p.m.Unlock()

	for {
		// if such port already exists - increment it
		if _, ok := p.ports[in]; !ok {
			p.ports[in] = projName
			return in
		}
		in++
	}
}

func (p *PortManager) SaveBatch(ports []uint16, projName string) {
	p.m.Lock()
	defer p.m.Unlock()

	for _, port := range ports {
		p.ports[port] = projName
	}
}
