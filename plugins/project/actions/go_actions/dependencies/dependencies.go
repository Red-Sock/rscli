package dependencies

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"
	"github.com/godverv/matreshka/servers"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
)

var (
	ErrNoFolderInConfig = errors.New("no folder path in rscli config")
)

type Dependency interface {
	AppendToProject(proj interfaces.Project) error
}

func GetDependencies(c *config.RsCliConfig, args []string) []Dependency {
	serverOpts := make([]Dependency, 0, len(args))

	for _, name := range args {
		idx := strings.Index(name, "_")
		resourceName := name
		if idx != -1 {
			resourceName = name[:idx]
		}
		var dep Dependency
		switch resourceName {
		case resources.PostgresResourceName:
			dep = Postgres{Cfg: c, Name: name}
		case resources.RedisResourceName:
			dep = Redis{Cfg: c, Name: name}
		case resources.TelegramResourceName:
			dep = Telegram{Cfg: c, Name: name}
		case resources.SqliteResourceName:
			dep = Sqlite{Cfg: c, Name: name}
		case servers.RestServerType:
			dep = Rest{Cfg: c, Name: name}
		case servers.GRPSServerType:
			dep = GrpcServer{Cfg: c, Name: name}
		default:
			continue
		}

		serverOpts = append(serverOpts, dep)
	}

	return serverOpts
}

// containsDependencyFolder - searches through RSCLI_PATH_TO_CLIENTS
// folders in order to find depName
// IF Dependency already placed - returns path to it
func containsDependencyFolder(paths []string, rootF *folder.Folder, depName string) (ok bool, err error) {
	if len(paths) == 0 {
		return false, errors.Wrap(ErrNoFolderInConfig, "no client")
	}

	for _, clientPath := range paths {
		clientFolder := rootF.GetByPath(clientPath)
		if clientFolder == nil {
			continue
		}

		for _, cF := range clientFolder.Inner {
			if cF.Name == depName {
				return true, nil
			}
		}
	}

	return false, nil
}
