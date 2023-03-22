package project

import (
	"github.com/Red-Sock/rscli/plugins/project/processor"
	"os"
)

func tidyProject(_ []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return processor.Tidy(wd)
}
