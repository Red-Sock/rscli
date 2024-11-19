package processor

import (
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
)

func WithIo(io io.IO) opt {
	return func(p *Processor) {
		p.IO = io
	}
}

func WithWd(wd string) opt {
	return func(p *Processor) {
		p.WD = wd
	}
}

func WithConfig(cfg *config.RsCliConfig) opt {
	return func(p *Processor) {
		p.RscliConfig = cfg
	}
}
