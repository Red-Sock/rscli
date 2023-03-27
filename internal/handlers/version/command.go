package version

import (
	"bytes"

	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

const Command = "version"

type Handler struct{}

func (h *Handler) Do(args []string) error {
	version := patterns.RscliMK[bytes.IndexByte(patterns.RscliMK, '=')+1 : bytes.IndexByte(patterns.RscliMK, '\n')]
	println("Current project version is " + string(version))

	return nil
}
