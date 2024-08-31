package dependencies

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/config"
	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
)

var (
	ErrNoFolderInConfig = errors.New("no folder path in rscli config")
)

type Dependency interface {
	AppendToProject(proj project.Project) error
}

type dependencyBase struct {
	Name string
	Cfg  *rscliconfig.RsCliConfig
}

const (
	DependencyNamePostgres = "postgres"
	DependencyNameRedis    = "redis"
	DependencyNameTelegram = "telegram"
	DependencyNameSqlite   = "sqlite"
	DependencyNameRest     = "rest"
	DependencyNameGrpc     = "grpc"
)

var nameToDependencyConstructor = map[string]func(dep dependencyBase) Dependency{
	DependencyNamePostgres: func(dep dependencyBase) Dependency { return &Postgres{dependencyBase: dep} },
	DependencyNameRedis:    func(dep dependencyBase) Dependency { return &Redis{dep} },

	DependencyNameTelegram: func(dep dependencyBase) Dependency { return &Telegram{dep} },
	DependencyNameSqlite:   func(dep dependencyBase) Dependency { return &Sqlite{dep} },
	DependencyNameRest:     func(dep dependencyBase) Dependency { return &Rest{dep} },
	DependencyNameGrpc:     func(dep dependencyBase) Dependency { return &GrpcServer{dep} },
}

func GetDependencies(c *config.RsCliConfig, args []string) []Dependency {
	serverOpts := make([]Dependency, 0, len(args))

	for _, name := range args {
		idx := strings.Index(name, "_")
		resourceName := name
		if idx != -1 {
			resourceName = name[:idx]
		}
		depConstr, ok := nameToDependencyConstructor[resourceName]
		if !ok {
			continue
		}
		base := dependencyBase{
			Name: name,
			Cfg:  c,
		}
		serverOpts = append(serverOpts, depConstr(base))
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

func containsDependency(dataSources matreshka.DataSources, resource resources.Resource) bool {
	for _, ds := range dataSources {
		if ds.GetName() == resource.GetName() {
			return true
		}
	}

	return false
}
