package utils

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FileChooserSpec int

const (
	FileSpecJSON FileChooserSpec = iota
	FileSpecDirectory
)

type FileChooserEntry struct {
	widget.BaseWidget
	label    *widget.Label
	Path     string
	window   fyne.Window
	Spec     FileChooserSpec
	icon     *widget.Button
	bg       *canvas.Rectangle
	Callback func(string)
	hovered  bool
}

func NewFileChooserEntry(window fyne.Window, placeholder string) *FileChooserEntry {
	f := &FileChooserEntry{
		label:  widget.NewLabel(placeholder),
		window: window,
		Spec:   FileSpecJSON,
	}
	f.label.Wrapping = fyne.TextTruncate
	f.label.Truncation = fyne.TextTruncateEllipsis
	f.BaseWidget.ExtendBaseWidget(f)
	return f
}
func (f *FileChooserEntry) CreateRenderer() fyne.WidgetRenderer {
	f.bg = canvas.NewRectangle(theme.InputBackgroundColor())
	f.bg.StrokeColor = theme.InputBorderColor()
	f.bg.StrokeWidth = 1
	f.bg.CornerRadius = 6

	f.icon = widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		f.showChooser()
	})
	f.icon.Importance = widget.LowImportance

	scroll := container.NewHScroll(f.label)
	scroll.SetMinSize(fyne.NewSize(100, f.label.MinSize().Height))

	content := container.NewBorder(nil, nil, nil, f.icon, scroll)
	main := container.NewMax(f.bg, content)

	return &fileChooserEntryRenderer{
		entry: f,
		bg:    f.bg,
		main:  main,
		objects: []fyne.CanvasObject{
			f.bg,
			main,
		},
	}
}

func (f *FileChooserEntry) Tapped(_ *fyne.PointEvent) {
	f.showChooser()
}

func (f *FileChooserEntry) TappedSecondary(_ *fyne.PointEvent) {
	// This function is required to satisfy the interface but has no behavior for this widget.
}

func (f *FileChooserEntry) MouseIn(*desktop.MouseEvent) {
	f.hovered = true
	f.bg.FillColor = theme.HoverColor()
	f.bg.Refresh()
	f.Refresh()
}

func (f *FileChooserEntry) MouseOut() {
	f.hovered = false
	f.bg.FillColor = theme.InputBackgroundColor()
	f.bg.Refresh()
	f.Refresh()
}

func (f *FileChooserEntry) MouseMoved(*desktop.MouseEvent) {
	// This function is required to satisfy the interface but has no behavior for this widget.
}

func (f *FileChooserEntry) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (f *FileChooserEntry) showChooser() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter path manually (optional)")
	entry.SetText(f.Path)

	var dlg *dialog.CustomDialog

	openButton := widget.NewButtonWithIcon("Open", theme.ConfirmIcon(), func() {
		if entry.Text != "" {
			lastPath := f.Path
			f.SetPath(entry.Text)
			if f.Callback != nil {
				dlg.Hide()
				f.Callback(lastPath)
			}
			dlg.Hide()
		}
	})

	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		dlg.Hide()
	})

	pickButton := widget.NewButton("Browse...", func() {
		switch f.Spec {
		case FileSpecDirectory:
			dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
				if uri != nil {
					path := uri.Path()
					entry.SetText(path)
				}
			}, f.window)

		case FileSpecJSON:
			var fileDlg *dialog.FileDialog
			filter := storage.NewExtensionFileFilter([]string{".json"})

			fileDlg = dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
				if r != nil {
					path := r.URI().Path()
					entry.SetText(path)
				}
				fileDlg.Hide()
			}, f.window)

			fileDlg.SetFilter(filter)
			fileDlg.Show()
		}
	})

	// Custom layout without the default cancel button
	content := container.NewVBox(
		entry,
		pickButton,
		layout.NewSpacer(),
		container.NewHBox(cancelButton, layout.NewSpacer(), openButton),
	)

	dlg = dialog.NewCustomWithoutButtons("Choose Path", content, f.window)
	dlg.Resize(fyne.NewSize(400, 180))
	dlg.Show()
}

// Renderer

type fileChooserEntryRenderer struct {
	entry   *FileChooserEntry
	bg      *canvas.Rectangle
	main    *fyne.Container
	objects []fyne.CanvasObject
}

func (r *fileChooserEntryRenderer) Layout(size fyne.Size) {
	r.main.Resize(size)
	r.bg.Resize(size)
}

func (r *fileChooserEntryRenderer) MinSize() fyne.Size {
	return r.main.MinSize()
}

func (r *fileChooserEntryRenderer) Refresh() {
	r.bg.FillColor = theme.InputBackgroundColor()
	if r.entry.hovered {
		r.bg.FillColor = theme.HoverColor()
	}
	r.bg.Refresh()
	r.main.Refresh()
}

func (r *fileChooserEntryRenderer) BackgroundColor() color.Color {
	return color.NRGBA{0, 0, 0, 0}
}

func (r *fileChooserEntryRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *fileChooserEntryRenderer) Destroy() {
	// This function is required to satisfy the interface but has no behavior for this widget.
}

func (f *FileChooserEntry) SetPath(path string) {
	f.Path = path
	f.label.SetText(path)
}
