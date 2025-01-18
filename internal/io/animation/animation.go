package animation

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"go.redsock.ru/rerrors"
	"golang.org/x/term"
)

type Animation struct {
	frames   [][]byte
	duration time.Duration

	out *os.File
	in  *os.File

	width  int
	height int
}

func New(opts ...opt) (*Animation, error) {
	a := &Animation{
		frames: nil,
		out:    os.Stdout,
		in:     os.Stdin,
	}

	for _, o := range opts {
		err := o(a)
		if err != nil {
			return nil, rerrors.Wrap(err)
		}
	}

	return a, nil
}

func (a *Animation) Play() {
	startRow, startCol := a.getCursorPosition()

	// Print the step name and start the animation
	a.moveCursorTo(startRow, startCol)

	start := time.Now()
	frameIndex := 0

	// Animation loop
	for time.Since(start) < a.duration {
		// Move to the animation area
		a.moveCursorTo(startRow, startCol)

		// Print the current frame
		_, _ = a.out.Write(a.frames[frameIndex])

		// Advance to the next frame
		frameIndex = (frameIndex + 1) % len(a.frames)

		// Wait before rendering the next frame
		time.Sleep(200 * time.Millisecond)
	}

	// Clear the animation area and display the step completed
	a.clearAnimationArea(startRow, startCol)
	time.Sleep(10 * time.Second)
	a.moveCursorTo(startRow, startCol)
}

func (a *Animation) moveCursorTo(row, col int) {
	_, _ = a.out.WriteString(fmt.Sprintf("\033[%d;%dH", row, col))
}

func (a *Animation) clearAnimationArea(startRow, startCol int) {
	for i := 0; i < a.height; i++ {
		a.moveCursorTo(startRow+i, startCol)
		_, _ = a.out.WriteString(a.makeSpaces(a.width))
	}
	a.moveCursorTo(startRow, startCol)
}

// getCursorPosition queries and reads the current cursor position
func (a *Animation) getCursorPosition() (row, col int) {
	// Put the terminal in raw mode
	oldState, err := term.MakeRaw(int(a.in.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		os.Exit(1)
	}
	defer term.Restore(int(a.in.Fd()), oldState) // Ensure terminal is restored

	// Query the terminal for the current cursor position
	fmt.Print("\033[6n")

	// Read the response from the terminal
	reader := bufio.NewReader(a.in)
	response, _ := reader.ReadString('R')

	// Parse the response format "\033[row;colR"
	_, _ = fmt.Sscanf(response, "\033[%d;%dR", &row, &col)

	return row, col
}

func (a *Animation) makeSpaces(length int) string {
	return fmt.Sprintf("%*s", length, "")
}
