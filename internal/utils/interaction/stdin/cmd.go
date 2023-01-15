package stdin

import (
	"io"
	"os"
)

func Cmd() io.Reader {
	return os.Stdin
}
