package v0_0_9_alpha

import (
	"bytes"
	"github.com/Red-Sock/rscli/plugins/project/processor"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

func Do(p processor.Project) {
	mkFile := p.GetFolder().GetByPath(patterns.RsCliMkFileName)
	bytes.Replace(mkFile, [])
}
