package v0_0_23_alpha

import (
	"errors"

	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var Version = interfaces.Version{
	Major:      0,
	Minor:      0,
	Negligible: 21,
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
		connFile := p.GetFolder().GetByPath(
			patterns.InternalFolder,
			patterns.ClientsFolder,
			patterns.PostgresFolder,
			patterns.ConnFile,
		)
		if connFile != nil {
			connFile.Content = patterns.PgConnFile
		}
	}

	return p.GetFolder().Build()
}
