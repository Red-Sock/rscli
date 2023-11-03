package migrations

import (
	"github.com/godverv/matreshka/resources"
)

type MigrationTool interface {
	Install() error
	Version() (version string, err error)
	GetLatestVersion() (version string, err error)
	Migrate(pathToFolder string, resource resources.Resource) error
}
