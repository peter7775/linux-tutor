package gui

import (
	_ "embed"
	"image/color"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed fonts/regular.otf
var regularTTF []byte

//go:embed fonts/bold.otf
var boldTTF []byte

//go:embed fonts/italic.otf
var italicTTF []byte

//go:embed fonts/bolditalic.otf
var boldItalicTTF []byte

type lemonTheme struct{}

var _ fyne.Theme = (*lemonTheme)(nil)

func (l *lemonTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground, theme.ColorNameMenuBackground, theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x09, G: 0x09, B: 0x09, A: 0xff}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xd6, G: 0xd2, B: 0x3a, A: 0xff}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0xc9, G: 0xc3, B: 0x4a, A: 0xff}
	case theme.ColorNameButton, theme.ColorNameHover:
		return color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xff}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x5a, G: 0x56, B: 0x00, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x8f, G: 0x8a, B: 0x1f, A: 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x4b, G: 0x4b, B: 0x4b, A: 0xff}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (l *lemonTheme) Font(style fyne.TextStyle) fyne.Resource {
	switch {
	case style.Bold && style.Italic:
		return fyne.NewStaticResource("bolditalic.otf", boldItalicTTF)
	case style.Bold:
		return fyne.NewStaticResource("bold.otf", boldTTF)
	case style.Italic:
		return fyne.NewStaticResource("italic.otf", italicTTF)
	default:
		return fyne.NewStaticResource("regular.otf", regularTTF)
	}
}

func (l *lemonTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (l *lemonTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 6
	case theme.SizeNameText:
		return 15
	case theme.SizeNameHeadingText:
		return 34
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameInlineIcon:
		return 18
	case theme.SizeNameScrollBar:
		return 12
	default:
		return theme.DefaultTheme().Size(name)
	}
}
