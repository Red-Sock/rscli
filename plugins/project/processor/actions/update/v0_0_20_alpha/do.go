package v0_0_20_alpha

import (
	"errors"

	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

var Version = interfaces.Version{
	Major:      0,
	Minor:      0,
	Negligible: 20,
	Additional: interfaces.TagVersionAlpha,
}

func Do(p interfaces.Project) (err error) {
	defer func() {
		if err != nil {
			return
		}

		updErr := Version.UpdateProjectVersion(p)
		if updErr == nil {
			return
		}

		if err == nil {
			err = updErr
			return
		}

		err = errors.Join(err, updErr)
	}()

	{
		// TODO
	}

	return p.GetFolder().Build()
}
