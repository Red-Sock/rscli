package v0_0_17_alpha

import (
	"errors"
	"path"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	interfaces2 "github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

var Version = interfaces2.Version{
	Major:      0,
	Minor:      0,
	Negligible: 17,
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
	p.GetFolder().Add(&folder.Folder{
		Name:    path.Join(patterns.InternalFolder, patterns.UtilsFolder, patterns.CloserFolder, patterns.CloserFile),
		Content: patterns.UtilsCloserFile,
	})

	connFile := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.ClientsFolder, patterns.PostgresFolder, patterns.ConnFileName)
	connFile.Content = patterns.PgConnFile

	go_actions.ReplaceProjectName(p.GetName(), connFile)

	err = p.GetFolder().Build()
	if err != nil {
		return err
	}

	_, err = cmd.Execute(cmd.Request{
		Tool:    "go",
		Args:    []string{"mod", "tidy"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return err
	}

	return p.GetFolder().Build()
}