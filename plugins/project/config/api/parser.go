package resources

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/utils/copier"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

var ErrUnknownResource = errors.New("unknown resource")

type ServerOption interface{}

func ParseServerOption(name string, in interface{}) (ServerOption, error) {

	dataSourceType := strings.Split(name, "_")[0]

	var r ServerOption

	switch dataSourceType {
	case patterns.TelegramServer:
		r = &Telegram{}
	default:
		return nil, errors.Wrapf(ErrUnknownResource, "unknown server option type %s", dataSourceType)
	}

	err := copier.Copy(in, &r)
	if err != nil {
		return nil, err
	}

	return r, nil

}
