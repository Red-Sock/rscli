package update

import (
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

func Do(p interfaces.Project) error {
	return nil
}

var updates = []func(p interfaces.Project) error{}
