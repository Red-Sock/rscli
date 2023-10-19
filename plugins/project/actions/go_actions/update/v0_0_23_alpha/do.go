package v0_0_23_alpha

import (
	"errors"

	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
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
			projpatterns.InternalFolder,
			projpatterns.ClientsFolder,
			projpatterns.PostgresFolder,
			projpatterns.ConnFileName,
		)
		if connFile != nil {
			connFile.Content = projpatterns.PgConnFile.Content
		}
	}

	{
		// add new way handling tx
		pgFolder := p.GetFolder().GetByPath(projpatterns.InternalFolder, projpatterns.ClientsFolder, projpatterns.PostgresFolder)
		if pgFolder != nil {
			pgFolder.Inner = []*folder.Folder{
				projpatterns.PgConnFile.Copy(),
				projpatterns.PgTxFile.Copy(),
			}
		}
	}

	return p.GetFolder().Build()
}
