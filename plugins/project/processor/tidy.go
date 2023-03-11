package processor

import "github.com/Red-Sock/rscli/plugins/project/processor/actions"

func Tidy(pathToProject string) error {
	p, err := LoadProject(pathToProject)
	if err != nil {
		return err
	}

	return actions.Tidy(p)
}
