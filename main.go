package main

import (
	"github.com/Red-Sock/rscli/internal"
	"github.com/Red-Sock/rscli/pkg/colors"
	"os"
)

func main() {
	err := internal.Run(os.Args[1:])
	if err != nil {
		_, _ = os.Stdout.WriteString(colors.TerminalColor(colors.ColorRed) + err.Error())
		defer os.Exit(1)
	}
}
