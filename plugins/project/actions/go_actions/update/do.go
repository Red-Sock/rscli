package update

import (
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/update/v0_0_26_alpha"
	interfaces2 "github.com/Red-Sock/rscli/plugins/project/interfaces"
)

func Do(p interfaces2.Project) error {
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
	interfaces2.Version
	do func(p interfaces2.Project) error
}

func GetLatestVersion() *Version {
	v := updates[len(updates)-1]
	return &v
}

var updates = []Version{
	{
		Version: v0_0_26_alpha.Version,
		do:      v0_0_26_alpha.Do,
	},
}
