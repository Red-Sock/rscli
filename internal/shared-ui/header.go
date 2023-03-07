package shared_ui

import (
	composit_label "github.com/Red-Sock/rscli-uikit/basic/composit-label"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli/pkg/service/help"
	"github.com/nsf/termbox-go"
)

func GetHeaderFromLabel(lbl *label.Label) *composit_label.ComplexLabel {
	cl := composit_label.New()
	cl.AddText(composit_label.TextPart{
		R:  []rune(help.Header),
		Fg: termbox.ColorRed,
		Bg: termbox.ColorBlack,
	})
	cl.AddLabel(lbl)

	return cl
}

func GetHeaderFromText(str string) *composit_label.ComplexLabel {
	cl := composit_label.New()
	cl.AddText(composit_label.TextPart{
		R:  []rune(help.Header),
		Fg: termbox.ColorRed,
		Bg: termbox.ColorBlack,
	})
	cl.AddText(composit_label.TextPart{
		R:  []rune(str),
		Fg: termbox.ColorYellow,
		Bg: termbox.ColorBlack,
	})

	return cl
}
