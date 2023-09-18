package tidy

import (
	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/patterns"
)

func DataSources(p interfaces.Project) error {
	err := applyDatasourceFolders(p)
	if err != nil {
		return errors.Wrap(err, "errors preparing data source folders")
	}

	return nil
}

func applyDatasourceFolders(p interfaces.Project) error {
	cfg := p.GetConfig()

	clientsFolder, err := cfg.GetDataSourceFolders()
	if err != nil {
		return errors.Wrap(err, "error obtaining clients folders from config")
	}
	if clientsFolder == nil {
		return nil
	}

	clientsFolderSrc := p.GetFolder().GetByPath(patterns.InternalFolder, patterns.ClientsFolder)
	if clientsFolderSrc != nil {
		clientsFolderSrc.Inner = nil
	}

	if len(clientsFolder.Inner) != 0 {
		p.GetFolder().GetByPath(patterns.InternalFolder).Add(clientsFolder)
	}

	err = p.GetFolder().Build()
	if err != nil {
		return errors.Wrap(err, "error building project after added clients")
	}

	return nil
}