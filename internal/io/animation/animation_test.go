package animation

import (
	"fmt"
	"testing"
	"time"
)

func Test_PlayAnimation(t *testing.T) {
	a := Animation{
		Frames: []string{
			`
###############
#             #
#   ▄         #
#             #
###############`,
			`
###############
#             #
#       ▄     #
#             #
###############`,
			`
###############
#             #
#         ▄   #
#             #
###############`,
			`
###############
#             #
#             #
#   ▄         #
###############`,
			`
###############
#             #
#             #
#         ▄   #
###############`,
		},
		IO: nil,
	}

	start := time.Now()
	frameIndex := 0

	stepName := "performing magic"
	duration := time.Second * 5
	// Animation loop
	for time.Since(start) < duration {
		// Clear the terminal
		fmt.Print("\033[H\033[2J")
		// Display the current frame
		fmt.Printf("%s\n%s\n", stepName, a.Frames[frameIndex])

		// Advance to the next frame
		frameIndex = (frameIndex + 1) % len(a.Frames)

		// Wait before rendering the next frame
		time.Sleep(200 * time.Millisecond)
	}

	// Clear the terminal after animation
	fmt.Print("\033[H\033[2J")

	// Display the step completed with a green box and a checkmark
	fmt.Printf("%s: \033[42;37m ✔ \033[0m\n", stepName)
}
