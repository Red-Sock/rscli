package processor

import "github.com/Red-Sock/rscli/plugins/project/processor/actions"

func Tidy(p *Project) error {
	return actions.Tidy(p)
}
