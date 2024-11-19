package init_new

import (
	"path"
)

func (p *Proc) collectOsPath(name string, args []string) (dirPath string) {
	if len(args) > 1 {
		return args[1]
	}

	return path.Join(p.WD, path.Base(name))
}
