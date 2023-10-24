package resources

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/utils/copier"
)

var ErrUnknownResource = errors.New("unknown resource")

type DataSourceName string

const (
	DataSourcePostgres = "postgres"
	DataSourceRedis    = "redis"
	DataSourceCustom   = "custom"
)

type Resource interface {
	// GetName - returns name defined in config file
	// e.g. if there is a need to use two different resources of the same type
	// user has to specify it a following way
	// data_sources:
	// 		redis_user_cache:
	//			....
	// 		redis_user_cart:
	//			....
	// then GetName() will return "user_cache" and "user_cart"
	GetName() string
	// GetType - returns on of DataSourceName types based on
	GetType() DataSourceName
	// GetEnv - returns a set of NAME-VALUE environment variables required for this resource to be run
	GetEnv() map[string]string
}

func ParseResource(name string, in interface{}) (Resource, error) {
	var dataSourceType string

	splitIdx := strings.Index(name, "_")

	if splitIdx == -1 {
		dataSourceType = name
		name = ""
	} else {
		dataSourceType = name[:splitIdx]
		name = name[splitIdx+1:]
	}

	var r Resource
	switch dataSourceType {
	case DataSourcePostgres:
		r = &Postgres{
			ResourceName: name,
		}
	case DataSourceRedis:
		r = &Redis{
			ResourceName: name,
		}
	default:
		return nil, errors.Wrapf(ErrUnknownResource, "unknown datasource type %s", dataSourceType)
	}

	err := copier.Copy(in, r)
	if err != nil {
		return nil, err
	}

	return r, nil

}
