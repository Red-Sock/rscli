package update

import (
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_10_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_15_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

func Do(p interfaces.Project) error {
	version := p.GetVersion()
	for _, item := range updates {
		if item.NeedUpdate(version) {
			err := item.do(p)
			if err != nil {
				return err
			}
		}
	}

	err := p.GetFolder().Build()
	if err != nil {
		return err
	}

	return nil
}

type Version struct {
	interfaces.Version
	do func(p interfaces.Project) error
}

var updates = []Version{
	{
		Version: v0_0_10_alpha.Version,
		do:      v0_0_10_alpha.Do,
	},
	{
		Version: v0_0_15_alpha.Version,
		do:      v0_0_15_alpha.Do,
	},
}
