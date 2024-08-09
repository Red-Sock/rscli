package dockerfile_actions

import (
	"bytes"

	"github.com/Red-Sock/rscli/plugins/project/proj_interfaces"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

type DockerFileTidyAction struct {
}

func (a DockerFileTidyAction) Do(p proj_interfaces.Project) error {
	dockerFile := p.GetFolder().GetByPath(projpatterns.Dockerfile.Name)
	if dockerFile == nil {
		dockerFile = projpatterns.Dockerfile.Copy()
		p.GetFolder().Add(dockerFile)
	}

	servers := p.GetConfig().Servers

	if len(servers) != 0 {
		exposeStart := getExposeStartIdx(dockerFile.Content)
		exposeEnd := 0

		if exposeStart != -1 {
			exposeEnd = exposeStart + bytes.IndexByte(dockerFile.Content[exposeStart:], '\n')
		} else {
			exposeStart = len(dockerFile.Content)
			exposeEnd = len(dockerFile.Content)
		}

		ports := make([][]byte, 0, len(p.GetConfig().Servers))
		for _, s := range servers {
			ports = append(ports, []byte(s.GetPortStr()))
		}
		exposeSequence := append([]byte(`EXPOSE `), bytes.Join(ports, []byte(" "))...)

		secondPart := make([]byte, len(dockerFile.Content[exposeEnd:]))
		copy(secondPart, dockerFile.Content[exposeEnd:])

		dockerFile.Content = append(dockerFile.Content[:exposeStart], exposeSequence...)
		dockerFile.Content = append(dockerFile.Content, secondPart...)
	}

	return nil
}

func (a DockerFileTidyAction) NameInAction() string {
	return "Tidy dockerfile"
}

func getExposeStartIdx(content []byte) int {
	return bytes.Index(content, []byte("EXPOSE"))
}