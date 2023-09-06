package scripts

import (
	"io/fs"
	"os"
	"path"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/patterns"
	"github.com/Red-Sock/rscli/internal/config"
)

func createEnvFolders(wd string, cfg *config.RsCliConfig) (err error) {
	err = prepareEnvDirBasic(wd)
	if err != nil {
		return err
	}

	var projects []string
	projects, err = ListProjects(wd, cfg)
	if err != nil {
		return errors.Wrap(err, "error obtaining projects list")
	}

	err = createEnvFoldersForProjects(wd, projects)
	if err != nil {
		return errors.Wrap(err, "error creating env folders for projects")
	}

	return nil
}
func prepareEnvDirBasic(wd string) (err error) {
	envDir := path.Join(wd, patterns.EnvDir)

	_, err = os.ReadDir(envDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = os.Mkdir(envDir, 0755)
			if err != nil {
				return errors.Wrap(err, "error making dir")
			}
		} else {
			return errors.Wrap(err, "error reading env directory: "+envDir)
		}
	}

	for _, f := range []patterns.File{patterns.EnvFile, patterns.DockerComposeFile, selectMakefile()} {
		err = createFileIfNotExists(path.Join(envDir, f.Name), f.Content)
		if err != nil {
			return errors.Wrap(err, "error writing "+f.Name+" file")
		}
	}
	return nil
}

func createEnvFoldersForProjects(wd string, projectsNames []string) error {
	envDir := path.Join(wd, patterns.EnvDir)
	for _, name := range projectsNames {

		err := os.Mkdir(path.Join(envDir, name), os.ModePerm)
		if err != nil {
			if !errors.Is(err, os.ErrExist) {
				return errors.Wrap(err, "error creating folder for project "+name)
			}
		}
	}

	return nil
}
