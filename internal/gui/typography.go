package gui

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewTitleLabel(text string) *widget.Label {
	l := widget.NewLabelWithStyle(text, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	l.Wrapping = fyne.TextTruncate
	return l
}

func NewSubheaderLabel(text string) *widget.Label {
	l := widget.NewLabelWithStyle(text, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	l.Wrapping = fyne.TextWrapWord
	return l
}

func NewSectionLabel(text string) *widget.Label {
	l := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return l
}

func NewQuestionTitle(text string) *widget.Label {
	l := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return l
}

func NewMutedText(text string) *widget.Label {
	l := widget.NewLabel(text)
	l.Wrapping = fyne.TextWrapWord
	return l
}
