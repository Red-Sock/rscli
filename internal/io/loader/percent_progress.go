package loader

import (
	"github.com/Red-Sock/rscli/internal/io/colors"
)

type percentLoader struct {
	Name      string
	AnimChars []string

	progress int

	progressUpdC chan string
}

func NewPercentLoader(name string, animationSymbs []string) *percentLoader {
	return &percentLoader{
		Name:         name,
		AnimChars:    animationSymbs,
		progressUpdC: make(chan string),
	}
}

func (p *percentLoader) UpdateProgress(newProgress int) {
	p.progress = newProgress

	if p.progress > 0 && p.progress <= 100 {
		p.progressUpdC <- p.GetLoaderSymb()
	}
}

func (p *percentLoader) GetLoaderSymb() string {
	if p.progress == 0 {
		return colors.TerminalColor(colors.ColorRed) + "X"
	}

	if p.progress >= 100 {
		return colors.TerminalColor(colors.ColorGreen) + p.AnimChars[len(p.AnimChars)-1]
	}

	idx := int(float32(p.progress%100*len(p.AnimChars)) / 100)
	return colors.TerminalColor(colors.ColorYellow) + p.AnimChars[idx]
}

func (p *percentLoader) GetName() string {
	return p.Name
}

func (p *percentLoader) GetProgressChan() <-chan string {
	return p.progressUpdC
}
