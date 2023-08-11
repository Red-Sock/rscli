package spinner

import (
	"sync"
	"time"

	"github.com/Red-Sock/rscli/internal/stdio"
)

func RhombSpinner() []string {
	return []string{"⬖", "⬘", "⬗", "⬙"}
}

// Start - continuously prints symbols from Spinner.SpinSet within Spinner.Timeout between each print
func Start(std stdio.IO, timeout time.Duration, spinSet []string) (stopFunc func(finalState string)) {
	t := time.NewTicker(timeout)

	stopChan := make(chan string)

	wg := &sync.WaitGroup{}

	stopFunc = func(finalState string) {
		wg.Add(1)
		stopChan <- finalState
		wg.Wait()
	}

	go func() {
		var idx int
		for {
			select {
			case <-t.C:
				std.Print(spinSet[idx] + "\010")

				idx++
				if idx >= len(spinSet) {
					idx = 0
				}
			case finalState := <-stopChan:
				std.Print(finalState + "\010")
				wg.Done()
				return
			}
		}
	}()

	return stopFunc
}
