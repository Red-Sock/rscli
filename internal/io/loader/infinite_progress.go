package loader

import (
	"sync"
	"time"

	"github.com/Red-Sock/rscli/internal/io/colors"
)

type InfiniteLoader struct {
	Name      string
	AnimChars []string

	animSymb int

	doneM  *sync.Mutex
	ticker *time.Ticker
	isDone bool

	progressUpdC chan string
}

func NewInfiniteLoader(name string, animationSymbs []string) *InfiniteLoader {

	il := &InfiniteLoader{
		Name:         name,
		AnimChars:    animationSymbs,
		progressUpdC: make(chan string),
		doneM:        &sync.Mutex{},
		ticker:       time.NewTicker(time.Second / 8),
	}

	go func() {

		idx := 0
		for {
			select {
			case <-il.ticker.C:

				il.doneM.Lock()
				if il.isDone {
					return
				}
				il.progressUpdC <- il.AnimChars[idx]
				il.doneM.Unlock()

				idx++
				if idx >= len(il.AnimChars) {
					idx = 0
				}
			}

		}
	}()

	return il
}

func (p *InfiniteLoader) Done(success progressStatus) {
	p.doneM.Lock()
	defer p.doneM.Unlock()

	if p.isDone {
		return
	}

	p.isDone = true
	p.ticker.Stop()

	switch success {
	case DoneSuccessful:
		p.progressUpdC <- colors.TerminalColor(colors.ColorGreen) + "*"
	case DoneFailed:
		p.progressUpdC <- colors.TerminalColor(colors.ColorRed) + "X"
	case DoneNotAccessed:
		p.progressUpdC <- colors.TerminalColor(colors.ColorYellow) + "?"
	}

	close(p.progressUpdC)
}

func (p *InfiniteLoader) GetName() string {
	return p.Name
}

func (p *InfiniteLoader) GetProgressChan() <-chan string {
	return p.progressUpdC
}
