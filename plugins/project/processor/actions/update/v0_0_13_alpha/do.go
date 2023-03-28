package v0_0_13_alpha

import (
	"errors"

	"github.com/Red-Sock/rscli/plugins/project/processor/actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var Version = interfaces.Version{
	Major:      0,
	Minor:      0,
	Negligible: 13,
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

	connFile := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.ClientsFolder, patterns.PostgresFolder, patterns.ConnFile)
	connFile.Content = patterns.PgConn
	renamer.ReplaceProjectName(p.GetName(), connFile)

	return p.GetFolder().Build()
}
