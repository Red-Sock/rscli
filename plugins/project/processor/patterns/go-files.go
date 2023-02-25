package patterns

import (
	_ "embed"
	"github.com/Red-Sock/rscli/plugins/project/processor/consts"
)

var DatasourceClients = map[consts.DataSourcePrefix][]byte{}
var ServerOptsPatterns = map[consts.ServerOptsPrefix]map[string][]byte{}

const (
	ServerGoFile  = "server.go"
	versionGoFile = "version.go"
)

func init() {
	DatasourceClients[consts.RedisDataSourcePrefix] = RedisConn
	DatasourceClients[consts.PostgresDataSourcePrefix] = PgConn

	ServerOptsPatterns[consts.RESTServerPrefix] = map[string][]byte{ServerGoFile: RestServ, versionGoFile: RestServVersion}
	// TODO
	ServerOptsPatterns[consts.GRPCServerPrefix] = map[string][]byte{}
}

// Basic files
var (
	//go:embed pattern_c/cmd/financial-microservice/main.go.pattern
	MainFile []byte
	//go:embed pattern_c/cmd/financial-microservice/api.go.pattern
	APISetupFile []byte
)

// DataStorage connection files
var (
	//go:embed pattern_c/internal/clients/redis/conn.go.pattern
	RedisConn []byte
	//go:embed pattern_c/internal/clients/postgres/conn.go.pattern
	PgConn []byte
)

// Config parser files
var (
	//go:embed pattern_c/internal/config/config.go.pattern
	Configurator string
	//go:embed pattern_c/internal/config/keys.go.pattern
	ConfigKeys []byte
)

// Server files
var (
	//go:embed pattern_c/internal/transport/manager.go.pattern
	ServerManagerPattern []byte

	//go:embed pattern_c/internal/transport/rest_realisation/server.go.pattern
	RestServ []byte
	//go:embed pattern_c/internal/transport/rest_realisation/version.go.pattern
	RestServVersion []byte

	// TODO
	GrpcServ []byte
)
