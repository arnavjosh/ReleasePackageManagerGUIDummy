package utils

import (
	"image/color"
	"os"

	"rpmg/assets"
	"rpmg/consts"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func MakeArrow() *canvas.Image {
	arrowResource := fyne.NewStaticResource("small-arrow.png", assets.SmallArrowPNG)
	if arrowResource == nil {
		return nil
	}
	arrow := canvas.NewImageFromResource(arrowResource)
	arrow.FillMode = canvas.ImageFillContain
	arrow.SetMinSize(fyne.NewSize(20, 20))
	return arrow
}

func MakeMoveableIcon() *canvas.Image {
	moveableResource := fyne.NewStaticResource("moveable.png", assets.MoveablePNG)
	if moveableResource == nil {
		return nil
	}
	moveable := canvas.NewImageFromResource(moveableResource)
	moveable.FillMode = canvas.ImageFillContain
	moveable.SetMinSize(fyne.NewSize(20, 20))
	return moveable
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func NewPaddingContainer(obj fyne.CanvasObject, paddingV, paddingH float32) fyne.CanvasObject {
	topPad := canvas.NewRectangle(color.Transparent)
	topPad.SetMinSize(fyne.NewSize(0, paddingV))

	bottomPad := canvas.NewRectangle(color.Transparent)
	bottomPad.SetMinSize(fyne.NewSize(0, paddingV))

	leftPad := canvas.NewRectangle(color.Transparent)
	leftPad.SetMinSize(fyne.NewSize(paddingH, 0))

	rightPad := canvas.NewRectangle(color.Transparent)
	rightPad.SetMinSize(fyne.NewSize(paddingH, 0))

	return container.NewVBox(
		topPad,
		container.NewHBox(
			leftPad,
			obj,
			rightPad,
		),
		bottomPad,
	)
}

func NewLightBlueLabel(text string) *canvas.Text {
	label := canvas.NewText(text, consts.LIGHT_BLUE)
	return label
}
