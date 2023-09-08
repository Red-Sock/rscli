package loader

import (
	"context"
	"sync"

	"github.com/morikuni/aec"

	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/colors"
)

// RunSeqLoader - allows to run sequential progress loader
func RunSeqLoader(ctx context.Context, io io.IO, progresses <-chan Progress) (done func() chan struct{}) {
	io.Print(aec.Hide.String())

	doneC := make(chan struct{})
	wg := &sync.WaitGroup{}
	io.Println()

	startDrawFunc := func(p Progress) {
		wg.Add(1)
		defer wg.Done()
		defer io.Print(aec.Show.String())
		defer io.Println()

		io.Print("_" + p.GetName())

		for {
			select {
			case v, ok := <-p.GetProgressChan():
				if !ok {
					return
				}

				io.Print(
					aec.Column(0).String() +
						v + colors.TerminalColor(colors.ColorDefault))

			}
		}
	}

	go func() {
		defer close(doneC)

		p, ok := <-progresses
		if !ok {
			return
		}

		if p == nil {
			io.PrintlnColored(colors.ColorRed, "Failed to start")
			return
		}
		startDrawFunc(p)
		for {
			select {
			case newP, newOk := <-progresses:
				if !newOk {
					return
				}
				// DONE SHOULD BE CALLED FROM THIS FUNCTION CALLER.
				// pass DoneNotAccessed in order to stop previously running loader.
				// Passing new value should not override existing value.
				p.Done(DoneNotAccessed)

				// Need to wait after passed done signal and before it is rendered
				wg.Wait()

				p = newP
				startDrawFunc(p)
			case <-ctx.Done():
				return
			}
		}
	}()

	return func() chan struct{} {
		return doneC
	}
}
