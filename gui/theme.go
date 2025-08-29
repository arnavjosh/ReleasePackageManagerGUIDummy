package main

import (
	"image/color"

	"rpmg/consts"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct{}

func (customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return consts.NAVY_BLUE
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNameMenuBackground:
		return consts.NAVY_BLUE
	case theme.ColorNameHeaderBackground:
		return consts.NAVY_BLUE
	case theme.ColorNameOverlayBackground:
		return consts.NAVY_BLUE
	case theme.ColorNameInputBackground:
		return consts.NAVY_BLUE
	case theme.ColorNameInputBorder:
		return consts.LIGHT_BLUE
	case theme.ColorNameButton:
		return color.Transparent
	case theme.ColorNameHover:
		return consts.SLATE_BLUE
	case theme.ColorNameFocus:
		return consts.SLATE_BLUE
	case theme.ColorNameDisabledButton:
		return consts.NAVY_BLUE
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 79, G: 123, B: 148, A: 255} // #4F7B94
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (customTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameText:
		return 16
	default:
		return theme.DefaultTheme().Size(name)
	}
}
