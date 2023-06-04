package shared_ui

import (
	composit_label "github.com/Red-Sock/rscli-uikit/basic/composit-label"
	"github.com/Red-Sock/rscli-uikit/basic/label"
	"github.com/Red-Sock/rscli-uikit/utils/common"
	"github.com/nsf/termbox-go"
)

const (
	Header = `
========================
    RedSock CLI tool    
========================
`
)

func GetHeaderFromLabel(lbl *label.Label) *composit_label.ComplexLabel {
	cl := composit_label.New(
		composit_label.Position(common.NewRelativePositioning(common.NewFillSpacePositioning(), common.NewFillSpacePositioning(), 0.4, 0.4)),
	)
	cl.AddText(composit_label.TextPart{
		R:  []rune(Header),
		Fg: termbox.ColorRed,
		Bg: termbox.ColorBlack,
	})
	cl.AddLabel(lbl)

	return cl
}

func GetHeaderFromText(str string) *composit_label.ComplexLabel {
	cl := composit_label.New(composit_label.SetBounds(30, 40))
	cl.AddText(composit_label.TextPart{
		R:  []rune(Header),
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
