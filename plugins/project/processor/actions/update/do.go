package update

import (
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_18_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_20_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_21_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_23_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_24_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"

	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_10_alpha"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/update/v0_0_17_alpha"
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
		Version: v0_0_17_alpha.Version,
		do:      v0_0_17_alpha.Do,
	},
	{
		Version: v0_0_18_alpha.Version,
		do:      v0_0_18_alpha.Do,
	},
	{
		Version: v0_0_20_alpha.Version,
		do:      v0_0_20_alpha.Do,
	},
	{
		Version: v0_0_21_alpha.Version,
		do:      v0_0_21_alpha.Do,
	},
	{
		Version: v0_0_23_alpha.Version,
		do:      v0_0_23_alpha.Do,
	},

	{
		Version: v0_0_24_alpha.Version,
		do:      v0_0_24_alpha.Do,
	},
}
