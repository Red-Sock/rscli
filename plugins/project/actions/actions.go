package actions

import (
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
)

type Action interface {
	Do(p interfaces.Project) error
	NameInAction() string
}
