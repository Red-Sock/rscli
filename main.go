package main

import (
	"github.com/Red-Sock/rscli/internal"
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/pkg/colors"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"os"
)

func main() {
	err := internal.Run(os.Args[1:])
	if err != nil {
		_, _ = os.Stdout.WriteString(colors.TerminalColor(colors.ColorRed) + err.Error())
		os.Exit(1)
	}
	_, _ = os.Stdout.WriteString(help.Header + randomizer.GoodGoodBuy())
}
