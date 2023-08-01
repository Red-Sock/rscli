package help

import (
	"fmt"
	"strings"

	initUtil "github.com/Red-Sock/rscli/internal/cli-utils/init"
)

const (
	UtilityName = "help"
)

func GetHelpMessage() string {
	return `
Usage: rscli [utility] [command] [options]

Basic utilities:
` +
		fmt.Sprintf(`%s: %s`, initUtil.UtilityName, shiftInOrder(initUtil.GetHelpMessage(), 1))
}

func shiftInOrder(in string, tabsAmount int) string {
	out := strings.Split(in, "\n")
	tabs := strings.Repeat("    ", tabsAmount)
	for idx := range out {
		out[idx] = tabs + out[idx]
	}
	return strings.Join(out, "\n")
}
