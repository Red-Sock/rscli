package dependencies

import (
	errors "github.com/Red-Sock/trace-errors"
)

type Rest struct {
	dependencyBase
}

func (r Rest) GetFolderName() string {
	if r.Name != "" {
		return r.Name
	}

	return "rest"
}

func (r Rest) AppendToProject(proj Project) error {
	err := r.applyFolder(proj, r.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error applying rest folder")
	}

	applyServerFolder(proj)
	return nil
}

func (r Rest) applyFolder(proj Project, defaultApiName string) error {
	ok, err := containsDependencyFolder(r.Cfg.Env.PathToServers, proj.GetFolder(), r.GetFolderName())
	if err != nil {
		return errors.Wrap(err, "error searching dependencies")
	}

	if ok {
		return nil
	}
	//serverF := projpatterns.RestServFile.CopyWithNewName(
	//	path.Join(r.Cfg.Env.PathToServers[0], defaultApiName, projpatterns.RestServFile.Name))

	//serverF.Content = renamer.ReplaceProjectNameFull(serverF.Content, proj.GetName())
	//
	//renamer2.ReplaceProjectName(proj.GetName(), serverF)
	//
	//proj.GetFolder().Add(
	//	serverF,
	//	projpatterns.RestServHandlerVersionExampleFile.CopyWithNewName(
	//		path.Join(r.Cfg.Env.PathToServers[0], defaultApiName,
	//			projpatterns.RestServHandlerVersionExampleFile.Name)),
	//)

	return nil
}
