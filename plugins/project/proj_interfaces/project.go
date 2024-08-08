package proj_interfaces

import (
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config"
)

type Project interface {
	GetName() string
	GetShortName() string

	GetConfig() *config.Config

	GetFolder() *folder.Folder
	GetProjectPath() string

	GetVersion() Version
	SetVersion(Version)

	GetType() ProjectType
}

type ProjectType string

const (
	ProjectTypeUnknown ProjectType = "Unknown"
	ProjectTypeGo      ProjectType = "go"
)
