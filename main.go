package main

import (
	"github.com/Red-Sock/rscli/internal"
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/pkg/colors"
)

func main() {
	err := internal.Run()
	if err != nil {
		println(colors.TerminalColor(colors.ColorRed) + err.Error())
	}
	println(randomizer.GoodGoodBuy())
}
