package cli

import (
	"testing"
)

func TestColoredText_Init(t *testing.T) {
	ct := NewColoredText("colored")
	t.Logf("here is a [%v] string", ct.SetFG(BrightWhite).Init().SetBG(BrightBlue).String())
}

func TestColoredText_Reset(t *testing.T) {
	ct := NewColoredText("colored")
	t.Logf("here is a [%v] string", ct.SetFG(BrightWhite).Reset().SetBG(BrightBlue).String())
}

func TestColoredText_SetBG(t *testing.T) {
	ct := NewColoredText("colored")
	t.Logf("here is a [%v] string", ct.SetBG(BrightBlue).String())
}

func TestColoredText_SetFG(t *testing.T) {
	ct := NewColoredText("colored")
	t.Logf("here is a [%v] string", ct.SetFG(BrightBlue).String())
}

func TestColoredText_String(t *testing.T) {
	ct := NewColoredText("colored")
	t.Logf("here is a [%v] string", ct.SetBG(Red).SetFG(BrightWhite).String())
}

func TestNewColoredText(t *testing.T) {
	ct := NewColoredText("colored")
	t.Logf("here is a [%v] string", ct.SetFG(BrightWhite).
		SetBG(BrightBlue).
		SetAttr(Bold, true).
		SetAttr(Bold, false).
		SetAttr(Underscore, true).
		SetAttr(Blink, true).
		SetAttr(ReverseVideo, true).
		SetAttr(Concealed, true).
		String())
	t.Logf("here is a [%v] string", NewColoredText("empty").
		String())
	t.Logf("end...")
}
