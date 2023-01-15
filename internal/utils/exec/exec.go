package exec

import (
	"os/exec"
)

func Run(cmd exec.Cmd, workDir string) error {
	cmd.Dir = workDir

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
