package init

import (
	"fmt"
)

const (
	UtilityName = "init"
)

func GetHelpMessage() string {
	return fmt.Sprintf(`
Usage: rscli init [target] [options]

Targets:
    project - Opens project constructor or simply creates project based on local|global configuration
Options:
    [go, vue, react] - You can specify what predefined global pattern to use for new project

`)
}
