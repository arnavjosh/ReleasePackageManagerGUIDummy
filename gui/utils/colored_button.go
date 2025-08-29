package utils

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ButtonColors struct {
	Fill    color.Color
	Text    color.Color
	Outline color.Color
}

type ColoredButton struct {
	widget.BaseWidget
	Label          string
	DefaultColors  ButtonColors
	HoverColors    ButtonColors
	DisabledColors ButtonColors
	SelectedColors ButtonColors

	minSize fyne.Size

	OnTapped func()

	isHovered  bool
	isSelected bool
	disabled   bool

	bold bool
	size float32

	bg     *canvas.Rectangle
	label  *canvas.Text
	border *canvas.Rectangle
}

func (b *ColoredButton) SetBold(bold bool) {
	b.bold = bold
}

func (b *ColoredButton) SetSize(size float32) {
	b.size = size
}

func NewButtonColors(fill, text, outline color.Color) ButtonColors {
	return ButtonColors{
		Fill:    fill,
		Text:    text,
		Outline: outline,
	}
}

func NewColoredButton(label string, defaultColor, hoverColor, disabledColor, selectedColor ButtonColors, tapped func()) *ColoredButton {
	btn := &ColoredButton{
		Label:          label,
		DefaultColors:  defaultColor,
		HoverColors:    hoverColor,
		DisabledColors: disabledColor,
		SelectedColors: selectedColor,
		OnTapped:       tapped,
	}
	btn.size = 16
	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *ColoredButton) CreateRenderer() fyne.WidgetRenderer {
	b.bg = canvas.NewRectangle(b.DefaultColors.Fill)
	b.bg.CornerRadius = 4

	b.label = canvas.NewText(b.Label, b.DefaultColors.Text)
	b.label.Alignment = fyne.TextAlignCenter
	b.label.TextStyle.Bold = b.bold
	b.label.TextSize = b.size

	b.border = canvas.NewRectangle(b.DefaultColors.Outline)
	b.border.CornerRadius = 4
	b.border.StrokeWidth = 3
	b.border.FillColor = color.Transparent

	paddedLabel := NewPaddingContainer(b.label, 2, 15)

	container := container.NewStack(b.bg, b.border, paddedLabel)
	return &coloredButtonRenderer{
		ColoredButton: b,
		container:     container,
	}
}

func (r *coloredButtonRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)

	// The last object in the container stack is the padded label container
	paddedLabel := r.container.Objects[len(r.container.Objects)-1]

	min := paddedLabel.MinSize()
	paddedLabel.Move(fyne.NewPos(
		(size.Width-min.Width)/2,
		(size.Height-min.Height)/2,
	))
}

type coloredButtonRenderer struct {
	ColoredButton *ColoredButton
	container     *fyne.Container
}

func (r *coloredButtonRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (b *ColoredButton) MinSize() fyne.Size {
	if !b.minSize.IsZero() {
		return b.minSize
	}
	return b.BaseWidget.MinSize()
}

func (r *coloredButtonRenderer) Refresh() {
	canvas.Refresh(r.ColoredButton)
}

func (r *coloredButtonRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *coloredButtonRenderer) Objects() []fyne.CanvasObject {
	return r.container.Objects
}

func (r *coloredButtonRenderer) Destroy() {
	// This function is required to satisfy the interface but has no behavior for this widget.
}

func (b *ColoredButton) Tapped(*fyne.PointEvent) {
	if b.disabled {
		return
	}

	// Apply 'pressed' state instantly
	pressedFill := darkenColor(b.bg.FillColor, 0.85)
	pressedText := darkenColor(b.label.Color, 0.7)

	b.bg.FillColor = pressedFill
	b.label.Color = pressedText
	b.bg.Refresh()
	b.label.Refresh()

	go func() {
		time.Sleep(100 * time.Millisecond)

		fyne.Do(func() {
			// Instead of restoring original colors, call Refresh() to update colors based on current state
			b.Refresh()
		})
	}()

	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// The math works like this:
// c.RGBA() returns each channel (r, g, b, a) as uint32 values in the range [0, 65535].
// Shifting right by 8 bits (`r >> 8`) converts these to 8-bit values in [0, 255].
// Each of R, G, and B is multiplied by factor between 0.0 and 1.0 to make the color darker (closer to 0 = black).
// Opacity channel (A) isn't changed, so transparency is preserved
func darkenColor(c color.Color, factor float64) color.Color {
	r, g, b, a := c.RGBA()
	return color.NRGBA{
		R: uint8(float64(r>>8) * factor),
		G: uint8(float64(g>>8) * factor),
		B: uint8(float64(b>>8) * factor),
		A: uint8(a >> 8),
	}
}

func (b *ColoredButton) MouseIn(*desktop.MouseEvent) {
	if b.disabled || b.isSelected {
		return
	}
	b.isHovered = true
	b.Refresh()
}

func (b *ColoredButton) MouseOut() {
	if b.disabled || b.isSelected {
		return
	}
	b.isHovered = false
	b.Refresh()
}

func (b *ColoredButton) MouseMoved(*desktop.MouseEvent) {
	// This function is required to satisfy the interface but has no behavior for this widget.
}

func (b *ColoredButton) SetSelected(selected bool) {
	b.isSelected = selected
	b.isHovered = false
	b.Refresh()
}

func (b *ColoredButton) SetDisabled(disabled bool) {
	b.disabled = disabled
	b.Refresh()
}

func (b *ColoredButton) Refresh() {
	b.BaseWidget.Refresh()

	var colors ButtonColors = b.DefaultColors
	switch {
	case b.disabled:
		b.label.TextStyle.Bold = b.bold
		colors = b.DisabledColors
	case b.isSelected:
		b.label.TextStyle.Bold = true
		colors = b.SelectedColors
	case b.isHovered:
		b.label.TextStyle.Bold = b.bold
		colors = b.HoverColors
	default:
		b.label.TextStyle.Bold = b.bold
		colors = b.DefaultColors
	}

	b.bg.FillColor = colors.Fill
	b.label.Color = colors.Text
	b.border.StrokeColor = colors.Outline
	b.border.FillColor = color.Transparent

	b.bg.Refresh()
	b.label.Refresh()
	b.border.Refresh()
}

func (b *ColoredButton) SetMinSize(size fyne.Size) {
	b.minSize = size
	b.Refresh()
}
