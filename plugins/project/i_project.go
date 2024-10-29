package project

import (
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/config"
)

const (
	TypeUnknown Type = "Unknown"
	TypeGo      Type = "go"
)

type Type string

type IProject interface {
	GetName() string
	GetShortName() string

	GetConfig() *config.Config

	GetFolder() *folder.Folder
	GetProjectPath() string

	GetType() Type
}