package main

import (
	"flag"
	"image/color"

	"rpmg/assets"
	"rpmg/consts"
	"rpmg/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type formState struct {
	signContent   fyne.CanvasObject
	formContainer *fyne.Container
}

func main() {
	// Add command line flags
	var (
		awsCred               = flag.String("aws-cred", "", "Path to AWS credentials file")
		apiConfigPathOverride = flag.String("api", "", "Path to API configuration file")
		privateKeyPath        = flag.String("private", "", "Path to private key file")
		//schemaDir             = flag.String("schema", rpm.AbsSchemaDir, "Path to schema directory")
	)
	flag.Parse()

	a := app.NewWithID("com.globusmedical.releasepackagemanager")
	a.Settings().SetTheme(customTheme{})
	w := a.NewWindow("Release Package Manager")
	w.CenterOnScreen()

	state := newAppState(w)
	// Pass the flags to the app state
	state.configPaths = &configPaths{
		AWSCred:               *awsCred,
		ApiConfigPathOverride: *apiConfigPathOverride,
		PrivateKeyPath:        *privateKeyPath,
		SchemaDir:             "",
	}

	state.enableStagePublish(false)
	forms := &formState{}

	// Prepares tab content
	forms.signContent = state.buildUI()

	// main container (objects are stacked)
	forms.formContainer = container.NewStack(forms.signContent)

	tabsButtons := container.NewHBox(
		state.forceBox,
		layout.NewSpacer(),
		utils.NewRoundedBackgroundWithHug(state.btnSign, consts.LIGHT_BLUE, 4),
		utils.MakeArrow(),
		utils.NewRoundedBackgroundWithHug(state.btnStage, consts.LIGHT_BLUE, 4),
		utils.MakeArrow(),
		utils.NewRoundedBackgroundWithHug(state.btnPublish, consts.LIGHT_BLUE, 4),
	)

	// left side with forms
	leftSide := container.NewBorder(nil, tabsButtons, nil, nil, forms.formContainer)
	// pads left side
	leftSidePadded := container.NewPadded(leftSide)

	rightSide := state.artifacts.buildArtifactPathUI()
	rightSidePadded := container.NewPadded(rightSide)

	// https://pkg.go.dev/fyne.io/fyne/v2@v2.6.1/container#Split.SetOffset
	mainContent := container.NewHSplit(leftSidePadded, rightSidePadded)
	mainContent.Offset = 0.5

	headerBox := makeHeader("Release Package Manager")

	w.SetContent(container.NewBorder(
		headerBox, nil, nil, nil, mainContent))

	// temporarily not fullscreened so that graphic rendering at different sizes is easy to test
	// w.SetFullScreen(true)
	// Set the initial window size to fully render examples
	w.Resize(fyne.NewSize(1000, 600))

	w.ShowAndRun()
}

func makeHeader(text string) fyne.CanvasObject {
	header := canvas.NewText(text, color.White)
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.TextSize = 24 // Adjust size as needed
	header.Alignment = fyne.TextAlignLeading

	line := canvas.NewLine(consts.LIGHT_BLUE)
	line.StrokeWidth = 1
	line.Resize(fyne.NewSize(200, 1))

	iconRes := fyne.NewStaticResource("quitIcon.png", assets.QuitIconPNG)
	var quitBtn *widget.Button
	if iconRes != nil {
		quitBtn = widget.NewButtonWithIcon("", iconRes, func() {
			fyne.CurrentApp().Quit()
		})
		quitBtn.Importance = widget.LowImportance
		quitBtn.Resize(fyne.NewSize(32, 32))
	} else {
		// fallback text button if icon not found
		quitBtn = widget.NewButton("Quit", func() {
			fyne.CurrentApp().Quit()
		})
	}
	headerBar := container.NewBorder(nil, nil, nil, quitBtn, header)

	withoutLine := container.NewPadded(headerBar)
	withLine := container.NewVBox(
		withoutLine,
		line,
	)
	return withLine
}
