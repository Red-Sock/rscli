package io

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/io/colors"
)

type IO interface {
	Println(in ...string)
	Print(in string)
	PrintlnColored(color colors.Color, in ...string)
	PrintColored(color colors.Color, in string)

	Error(in string)

	GetInput() (string, error)
	GetInputOneOf(options []string) string
}

type StdIO struct{}

func (p StdIO) Println(in ...string) {
	for idx := range in {
		fmt.Print(in[idx])
	}
	fmt.Print("\n")
}
func (p StdIO) Print(in string) {
	fmt.Print(in)
}
func (p StdIO) PrintlnColored(color colors.Color, in ...string) {
	fmt.Print(colors.TerminalColor(color))
	p.Println(in...)
	fmt.Print(colors.TerminalColor(colors.ColorDefault))
}
func (p StdIO) PrintColored(color colors.Color, in string) {
	p.Print(colors.TerminalColor(color))
	p.Print(in)
	p.Print(colors.TerminalColor(colors.ColorDefault))
}
func (p StdIO) Error(in string) {
	p.Println("")
	_, _ = os.Stderr.WriteString(in)
}
func (p StdIO) GetInput() (string, error) {
	out, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return out, errors.Wrap(err, "error reading user input")
	}

	return out[:len(out)-1], nil
}
func (p StdIO) GetInputOneOf(options []string) string {
	panic("not implemented")
}
