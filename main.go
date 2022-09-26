package main

import (
	"github.com/Red-Sock/rscli/internal"
	"os"
)

func main() {
	internal.Run(os.Args[1:])
}
