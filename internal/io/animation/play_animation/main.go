package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/rscli/internal/io/animation"
)

func main() {
	println("Processing Step132")

	frames := []string{
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
	}

	a, err := animation.New(
		animation.WithWriter(os.Stdout, os.Stdin),
		animation.WithStrFrames(frames...),
		animation.WithDuration(time.Second*5),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	a.Play()

}
