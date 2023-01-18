package colors

import "github.com/nsf/termbox-go"

type Color uint64

const (
	ColorDefault Color = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

func UIColor(c Color) termbox.Attribute {
	return termbox.Attribute(c)
}

func TerminalColor(c Color) string {
	switch c {
	case ColorDefault:
		return "\033[0m"
	case ColorBlack:
		return "\033[30m"
	case ColorRed:
		return "\033[31m"
	case ColorGreen:
		return "\033[32m"
	case ColorYellow:
		return "\033[33m"
	case ColorBlue:
		return "\033[34m"
	case ColorMagenta:
		return "\033[35m"
	case ColorCyan:
		return "\033[36m"
	case ColorWhite:
		return "\033[39m"
	default:
		return "\033[0m"
	}
}
