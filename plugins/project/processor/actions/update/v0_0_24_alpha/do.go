package v0_0_24_alpha

import (
	"errors"

	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var Version = interfaces.Version{
	Major:      0,
	Minor:      0,
	Negligible: 24,
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
		// add new way handling tx
		pgFolder := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.ClientsFolder, patterns.PostgresFolder)
		if pgFolder != nil {
			pgFolder.Inner = []*folder.Folder{
				{
					Name:    patterns.ConnFile,
					Content: patterns.PgConnFile,
				},
				{
					Name:    patterns.PgTxFileName,
					Content: patterns.PgTxFile,
				},
			}
		}
	}

	return p.GetFolder().Build()
}
