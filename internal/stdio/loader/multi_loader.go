package loader

import (
	"context"
	"sync"

	"github.com/morikuni/aec"

	"github.com/Red-Sock/rscli/internal/stdio"
	"github.com/Red-Sock/rscli/pkg/colors"
)

type progressStatus int

const (
	InProgress progressStatus = iota
	// DoneSuccessful - Green. Successful. Good
	DoneSuccessful
	// DoneFailed - Red. Failed. Bad
	DoneFailed
	// DoneNotAccessed - Yellow. Not accessed during execution. Ok, not that bad
	DoneNotAccessed
)

type Progress interface {
	GetName() string
	Done(isSuccess progressStatus)
	GetProgressChan() <-chan string
}

// RunMultiLoader - allows to run multiple progresses and show state of each one simultaneously
// argument:
// ctx - to cancel run
// io - target to write to
// progresses - Progress interface objects
// return argument:
// done - function for done channel. Call it before exiting function from which RunMultiLoader was called
// in order to wait for all progresses to be completed and shown states
func RunMultiLoader(ctx context.Context, io stdio.IO, progresses []Progress) (done func() chan struct{}) {
	io.Print(aec.Hide.String())

	doneLoaders := &sync.WaitGroup{}

	doneLoaders.Add(len(progresses))

	printLock := &sync.Mutex{}

	for idx := range progresses {
		io.Println("_ ", progresses[idx].GetName())

		go func(uidx uint) {
			defer doneLoaders.Done()

			for {
				select {
				case v, ok := <-progresses[uidx].GetProgressChan():
					if !ok {
						return
					}
					printLock.Lock()

					io.Print(
						aec.Up(uint(len(progresses))-uidx).String() +
							aec.Column(1).String() +
							v + colors.TerminalColor(colors.ColorDefault) +
							aec.Down(uint(len(progresses))-uidx).String())

					printLock.Unlock()
				case <-ctx.Done():
				}
			}

		}(uint(idx))
	}

	doneC := make(chan struct{})
	done = func() chan struct{} {
		return doneC
	}

	go func() {
		doneLoaders.Wait()
		close(doneC)

		io.Print(aec.Show.String())
		io.Print(aec.Column(0).String())
	}()

	return done
}
