package makefile

import (
	"runtime"

	"go.redsock.ru/rerrors"

	"github.com/Red-Sock/rscli/internal/cmd"
)

var (
	ErrUnsupportedOS = rerrors.New("unsuported OS")
)

const bin = "make"

func Exists() bool {
	_, err := cmd.Execute(cmd.Request{
		Tool: bin,
		Args: []string{"--help"},
	})
	return err == nil
}

func Install() error {
	switch runtime.GOOS {
	case "darwin":
		if !Exists() {
			return installMacOS()
		}
	case "linux":
		if !Exists() {
			return installLinux()
		}
	default:
		return rerrors.Wrap(ErrUnsupportedOS, runtime.GOOS)
	}

	return nil
}

func Run(wd, mkFilePath string, command string) (string, error) {
	req := cmd.Request{
		Tool:    bin,
		Args:    []string{"-f", mkFilePath, command},
		WorkDir: wd,
	}

	msg, err := cmd.Execute(req)
	if err != nil {
		return "", rerrors.Wrap(err, "error running command:"+command+". "+err.Error())
	}

	return msg, nil
}

func installMacOS() error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    "brew",
		Args:    []string{"install", "make"},
		WorkDir: "",
	})
	if err != nil {
		return err
	}

	return nil
}

func installLinux() error {
	return nil
}
