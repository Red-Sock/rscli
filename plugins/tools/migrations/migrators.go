package migrations

import (
	"github.com/Red-Sock/rscli/plugins/project/config/resources"
)

type MigrationTool interface {
	Install() error
	Version() (version string, err error)
	GetLatestVersion() (version string, err error)
	Migrate(pathToFolder string, resource resources.Resource) error
}
