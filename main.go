package main

import (
	"github.com/Red-Sock/rscli/internal"
	"github.com/nsf/termbox-go"
	"os"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	internal.Run(os.Args[1:])
}
