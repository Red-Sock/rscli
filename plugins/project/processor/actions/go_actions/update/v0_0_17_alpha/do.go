package v0_0_17_alpha

import (
	"errors"

	"github.com/Red-Sock/rscli/pkg/cmd"
	"github.com/Red-Sock/rscli/pkg/folder"
	"github.com/Red-Sock/rscli/plugins/project/processor/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var Version = interfaces.Version{
	Major:      0,
	Minor:      0,
	Negligible: 17,
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
	p.GetFolder().AddWithPath([]string{patterns.InternalFolder, patterns.UtilsFolder, patterns.CloserFolder}, &folder.Folder{
		Name:    patterns.CloserFile,
		Content: patterns.UtilsCloserFile,
	})

	connFile := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.ClientsFolder, patterns.PostgresFolder, patterns.ConnFile)
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
