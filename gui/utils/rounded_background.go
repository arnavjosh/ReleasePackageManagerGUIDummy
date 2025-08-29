package utils

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RoundedBackground struct {
	widget.BaseWidget
	content      fyne.CanvasObject
	bgColor      color.Color
	cornerRadius float32
	rect         *canvas.Rectangle
	strokeWidth  float32
	hug          bool
}

func NewRoundedBackground(obj fyne.CanvasObject, bgColor color.Color, cornerRadius float32) fyne.CanvasObject {
	r := &RoundedBackground{
		content:      obj,
		bgColor:      bgColor,
		cornerRadius: cornerRadius,
		rect:         canvas.NewRectangle(bgColor),
		strokeWidth:  1,
	}
	r.rect.SetMinSize(obj.MinSize())
	r.rect.CornerRadius = cornerRadius
	r.rect.StrokeWidth = r.strokeWidth
	r.rect.StrokeColor = bgColor
	r.rect.FillColor = color.Transparent

	r.ExtendBaseWidget(r)
	r.hug = false
	return r
}

func NewRoundedBackgroundWithHug(obj fyne.CanvasObject, bgColor color.Color, cornerRadius float32) fyne.CanvasObject {
	r := NewRoundedBackground(obj, bgColor, cornerRadius).(*RoundedBackground)
	r.hug = true
	return r
}

func (r *RoundedBackground) CreateRenderer() fyne.WidgetRenderer {
	r.rect.CornerRadius = r.cornerRadius
	var objects []fyne.CanvasObject
	if r.hug {
		objects = []fyne.CanvasObject{r.content, r.rect}
	} else {
		objects = []fyne.CanvasObject{container.NewPadded(r.content), r.rect}
	}

	return &roundedBackgroundRenderer{
		bg:      r,
		objects: objects,
	}
}

type roundedBackgroundRenderer struct {
	bg      *RoundedBackground
	objects []fyne.CanvasObject
}

func (r *roundedBackgroundRenderer) Layout(size fyne.Size) {
	r.bg.rect.Resize(size)
	r.objects[0].Resize(size)
}

func (r *roundedBackgroundRenderer) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

func (r *roundedBackgroundRenderer) Refresh() {
	r.bg.rect.FillColor = color.Transparent
	r.bg.rect.StrokeColor = r.bg.bgColor
	r.bg.rect.StrokeWidth = r.bg.strokeWidth
	r.bg.rect.CornerRadius = r.bg.cornerRadius
	r.bg.rect.Refresh()
	r.objects[0].Refresh()
}

func (r *roundedBackgroundRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *roundedBackgroundRenderer) Destroy() {
	// No resources to clean up
}
