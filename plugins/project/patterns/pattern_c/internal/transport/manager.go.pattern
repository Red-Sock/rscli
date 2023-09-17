package transport

import (
	"context"

	"github.com/pkg/errors"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type ServersManager struct {
	serverPool []Server
}

func NewManager() *ServersManager {
	return &ServersManager{}
}

func (m *ServersManager) AddServer(server Server) {
	m.serverPool = append(m.serverPool, server)
}

func (m *ServersManager) Start(ctx context.Context) error {
	var errs []error

	for sID := range m.serverPool {
		err := m.serverPool[sID].Start(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	finalError := errors.New("error starting servers")

	for _, err := range errs {
		finalError = errors.Wrap(finalError, err.Error())
	}

	return finalError
}

func (m *ServersManager) Stop(ctx context.Context) error {
	var errs []error

	for sID := range m.serverPool {
		err := m.serverPool[sID].Stop(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	finalError := errors.New("error stopping servers")

	for _, err := range errs {
		finalError = errors.Wrap(finalError, err.Error())
	}

	return finalError
}
