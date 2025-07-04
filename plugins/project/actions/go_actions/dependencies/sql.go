package dependencies

import (
	"fmt"
	"path"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/renamer"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

type sqlConn struct {
	Cfg *rscliconfig.RsCliConfig
}

func (sc sqlConn) GetFolderName() string {
	return "sqldb"
}

func (sc sqlConn) applySqlConnectionFile(proj Project) error {
	if len(sc.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	fileName := path.Join(
		sc.Cfg.Env.PathsToClients[0],
		sc.GetFolderName(),
		patterns.SqlConnFile.Name)

	sqlConnFile := patterns.SqlConnFile.CopyWithNewName(fileName)

	renamer.ReplaceProjectName(proj.GetName(), sqlConnFile)
	proj.GetFolder().Add(sqlConnFile)

	return nil
}

func (sc sqlConn) applySqlDriver(proj Project, driverName, driverImportPath string) {
	fileName := path.Join(
		sc.Cfg.Env.PathsToClients[0],
		sc.GetFolderName(),
		driverName+".go")

	sqlDriverFile := &folder.Folder{
		Name: fileName,
		Content: []byte(
			fmt.Sprintf(`package %s

import %s
`, sc.GetFolderName(), driverImportPath)),
	}

	proj.GetFolder().Add(sqlDriverFile)

}
