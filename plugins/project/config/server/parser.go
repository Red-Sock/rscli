package server

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/utils/copier"
	projpatterns "github.com/Red-Sock/rscli/plugins/project/patterns"
)

var ErrUnknownResource = errors.New("unknown resource")

type Server interface {
	GetName() string
	GetPort() uint16
}

func ParseServerOption(name string, in interface{}) (Server, error) {
	dataSourceType := strings.Split(name, "_")[0]

	var r Server

	switch dataSourceType {
	case projpatterns.TelegramServer:
		r = &Telegram{
			name: name,
		}
	case projpatterns.RESTHTTPServer:
		r = &Rest{
			name: name,
		}
	case projpatterns.GRPCServer:
		r = &GRPC{
			name: name,
		}
	default:
		return nil, errors.Wrapf(ErrUnknownResource, "unknown server option type %s", dataSourceType)
	}

	err := copier.Copy(in, r)
	if err != nil {
		return nil, err
	}

	return r, nil

}
