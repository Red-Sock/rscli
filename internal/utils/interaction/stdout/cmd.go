package stdout

import (
	"io"
	"os"
)

func Cmd() io.Writer {
	return os.Stdout
}
