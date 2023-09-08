package v0_0_24_alpha

import (
	"errors"

	interfaces2 "github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

var Version = interfaces2.Version{
	Major:      0,
	Minor:      0,
	Negligible: 24,
	Additional: interfaces2.TagVersionAlpha,
}

func Do(p interfaces2.Project) (err error) {
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
	}

	{
	}

	return p.GetFolder().Build()
}
