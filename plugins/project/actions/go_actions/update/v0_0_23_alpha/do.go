package v0_0_23_alpha

import (
	"errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

var Version = interfaces.Version{
	Major:      0,
	Minor:      0,
	Negligible: 23,
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
			patterns.ConnFileName,
		)
		if connFile != nil {
			connFile.Content = patterns.PgConnFile
		}
	}

	{
		// add new way handling tx
		pgFolder := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.ClientsFolder, patterns.PostgresFolder)
		if pgFolder != nil {
			pgFolder.Inner = []*folder.Folder{
				{
					Name:    patterns.ConnFileName,
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
