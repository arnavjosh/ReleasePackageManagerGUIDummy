package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"rpmg/consts"
	"rpmg/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type inputs struct {
	NewManifestButton *utils.ColoredButton
	ManifestPath      *utils.FileChooserEntry
	Platform          *widget.SelectEntry
	Application       *widget.SelectEntry
	VersionFrom       *widget.Entry
	BuildFrom         *widget.Entry
	VersionTo         *widget.Entry
	BuildTo           *widget.Entry
	RunBtn            *widget.Button
	AutofillBtn       *widget.Button
	devAPI            *utils.ColoredButton
	stageAPI          *utils.ColoredButton
	prodAPI           *utils.ColoredButton
	SelectedAPI       string
}

//note for Dr. Bezaire: originally this struct was written into the companies code, I am rewriting it here so that the code doesn't break in the dummy version

type VersionInfo struct {
	Version string
	Build   string
}

func newInputs(w fyne.Window, state *appState) *inputs {
	c := &inputs{
		Platform:    widget.NewSelectEntry([]string{"egps", "ehub", "e3d", "ltp"}),
		Application: widget.NewSelectEntry([]string{"gmed", "e3d"}),
		VersionFrom: widget.NewEntry(),
		BuildFrom:   widget.NewEntry(),
		VersionTo:   widget.NewEntry(),
		BuildTo:     widget.NewEntry(),
	}
	c.Platform.SetPlaceHolder("e.g. egps")
	c.Application.SetPlaceHolder("e.g. gmed")
	c.VersionFrom.SetPlaceHolder("e.g. 1.0.0")
	c.BuildFrom.SetPlaceHolder("e.g. 123")
	c.VersionTo.SetPlaceHolder("e.g. 1.1.0")
	c.BuildTo.SetPlaceHolder("e.g. 124")
	state.forceBox.SetChecked(false)

	// set new
	defaultColor := utils.NewButtonColors(consts.BRIGHT_BLUE, consts.BLACK, consts.TRANSPARENT)
	hoverColor := utils.NewButtonColors(consts.SLATE_BLUE, consts.BLACK, consts.TRANSPARENT)
	selectedColor := defaultColor
	disabledColor := defaultColor

	c.ManifestPath = utils.NewFileChooserEntry(w, "/..")
	c.ManifestPath.Callback = func(string) {
		c.validateManifestPath(w, state)
	}

	c.NewManifestButton = utils.NewColoredButton("New", defaultColor, hoverColor, selectedColor, disabledColor, func() {
		c.ManifestPath.SetPath("")
		emptyVersionInfo := VersionInfo{Version: "", Build: ""}
		c.setAllText("", "", emptyVersionInfo, emptyVersionInfo)
		c.enableDisableAll(true)
	})

	c.NewManifestButton.SetSize(14)
	c.NewManifestButton.SetBold(true)

	c.AutofillBtn = widget.NewButton("Autofill", func() {
		from := VersionInfo{Version: "6.0", Build: "rc11b"}
		to := VersionInfo{Version: "6.1", Build: "rc2a"}
		c.setAllText("egps", "gmed", from, to)
		c.devAPI.Refresh()
		c.stageAPI.Refresh()
		c.prodAPI.Refresh()
		state.artifacts.ArtifactPath.SetPath(`..\example\artifacts`)
		state.setManifestPath("")
	})

	defaultColor = utils.NewButtonColors(consts.NAVY_BLUE, consts.LIGHT_BLUE, consts.TRANSPARENT)
	hoverColor = utils.NewButtonColors(consts.SLATE_BLUE, consts.LIGHT_BLUE, consts.TRANSPARENT)
	selectedColor = utils.NewButtonColors(consts.SLATE_BLUE, consts.WHITE, consts.BRIGHT_BLUE)
	disabledColor = defaultColor

	c.devAPI = utils.NewColoredButton("dev", defaultColor, hoverColor, disabledColor, selectedColor, func() {
		c.setSelectedAPI("dev")
	})

	c.stageAPI = utils.NewColoredButton("stage", defaultColor, hoverColor, disabledColor, selectedColor, func() {
		c.setSelectedAPI("stage")
	})

	c.prodAPI = utils.NewColoredButton("prod", defaultColor, hoverColor, disabledColor, selectedColor, func() {
		c.setSelectedAPI("prod")
	})

	return c
}

func (c *inputs) validateManifestPath(w fyne.Window, state *appState) {
	manifestPath := c.ManifestPath.Path
	// gets the file name from the path
	manifestFile := filepath.Base(manifestPath)

	platform, application, from, to, err := ParseManifestName(manifestFile)
	if err != nil {
		dialog.ShowError(err, w)
	}
	c.setAllText(platform, application, from, to)
	c.enableDisableAll(false)
	state.setManifestPath(manifestPath)
}

// dummy version of the function (original is in the companies code)
func ParseManifestName(manifestFile string) (string, string, VersionInfo, VersionInfo, error) {
	return "egps", "gmed", VersionInfo{Version: "1.0.0", Build: "123"}, VersionInfo{Version: "1.1.0", Build: "124"}, nil
}

func (c *inputs) setSelectedAPI(api string) {
	c.SelectedAPI = api

	c.devAPI.SetSelected(false)
	c.stageAPI.SetSelected(false)
	c.prodAPI.SetSelected(false)

	// Set importance based on selected API

	switch api {
	case "dev":
		c.devAPI.SetSelected(true)
	case "stage":
		c.stageAPI.SetSelected(true)
	case "prod":
		c.prodAPI.SetSelected(true)
	}
	c.devAPI.Refresh()
	c.stageAPI.Refresh()
	c.prodAPI.Refresh()
}

func (c *inputs) buildInputUI() fyne.CanvasObject {
	packageHeader := canvas.NewText("Package", color.White)
	packageHeader.TextStyle = fyne.TextStyle{Bold: true}
	packageHeader.Alignment = fyne.TextAlignLeading
	packageHeader.TextSize = 18

	header := container.NewBorder(nil, nil, nil, c.NewManifestButton, packageHeader)

	manifestPickFull := container.NewVBox(
		utils.NewLightBlueLabel("Manifest Path"),
		c.ManifestPath,
	)

	// --- Top Row: Platform and Application side by side
	platformAppRow := container.NewGridWithColumns(2,
		container.NewVBox(
			utils.NewLightBlueLabel("Platform"),
			c.Platform,
		),
		container.NewVBox(
			utils.NewLightBlueLabel("Application"),
			c.Application,
		),
	)

	// --- From Section
	fromSection := container.NewVBox(
		utils.NewLightBlueLabel("From"),
		container.NewGridWithColumns(2,
			container.NewVBox(utils.NewLightBlueLabel("Version"), c.VersionFrom),
			container.NewVBox(utils.NewLightBlueLabel("Build"), c.BuildFrom),
		),
	)

	toSection := container.NewVBox(
		utils.NewLightBlueLabel("To"),
		container.NewGridWithColumns(2,
			container.NewVBox(utils.NewLightBlueLabel("Version"), c.VersionTo),
			container.NewVBox(utils.NewLightBlueLabel("Build"), c.BuildTo),
		),
	)

	// wraps from and to sections in rounded background
	fromSectionWrapped := utils.NewRoundedBackground(fromSection, consts.LIGHT_BLUE, 16)
	toSectionWrapped := utils.NewRoundedBackground(toSection, consts.LIGHT_BLUE, 16)

	upperBlock := container.NewVBox(
		platformAppRow,
		layout.NewSpacer(),
		container.NewGridWithColumns(2,
			fromSectionWrapped,
			toSectionWrapped,
		),
	)

	wrappedBlock := utils.NewRoundedBackground(upperBlock, consts.LIGHT_BLUE, 16)

	buttonsContainer := container.NewGridWithColumns(3,
		container.NewVBox(layout.NewSpacer(), c.devAPI, layout.NewSpacer()),
		container.NewVBox(layout.NewSpacer(), c.stageAPI, layout.NewSpacer()),
		container.NewVBox(layout.NewSpacer(), c.prodAPI, layout.NewSpacer()),
	)

	buttonsContainerWrapped := utils.NewRoundedBackgroundWithHug(buttonsContainer, consts.LIGHT_BLUE, 4)

	apiBlock := container.NewVBox(
		utils.NewLightBlueLabel("API Server"),
		buttonsContainerWrapped,
	)

	buttonRow := container.NewHBox(layout.NewSpacer(), layout.NewSpacer())

	return container.NewVBox(
		header,
		manifestPickFull,
		wrappedBlock,
		apiBlock,
		layout.NewSpacer(),
		buttonRow,
	)
}

func (c *inputs) createManifest(dir string) (string, error) {
	values := map[string]string{
		"platform":  strings.TrimSpace(c.Platform.Text),
		"app":       strings.TrimSpace(c.Application.Text),
		"verFrom":   strings.TrimSpace(c.VersionFrom.Text),
		"buildFrom": strings.TrimSpace(c.BuildFrom.Text),
		"verTo":     strings.TrimSpace(c.VersionTo.Text),
		"buildTo":   strings.TrimSpace(c.BuildTo.Text),
	}
	for k, v := range values {
		if v == "" {
			return "", fmt.Errorf("field '%s' is required. If there is a pre-existing manifest, then enter its path. Otherwise, fill out all other fields", k)
		}
	}
	manifestPath := createPath(values, dir)
	return manifestPath, nil
}

func (c *inputs) validateSignInputs() (map[string]string, error) {
	values := map[string]string{
		"manifestPath": strings.TrimSpace(c.ManifestPath.Path),
		"apiServer":    strings.TrimSpace(c.SelectedAPI),
	}

	for k, v := range values {
		if v == "" {
			return nil, fmt.Errorf("sign: field '%s' is required", k)
		}
	}
	return values, nil
}

func createPath(inputs map[string]string, directory string) string {
	manifestFile := fmt.Sprintf("%s-%s-%s-%s-to-%s-%s.manifest.json",
		inputs["platform"], inputs["app"],
		inputs["verFrom"], inputs["buildFrom"],
		inputs["verTo"], inputs["buildTo"])
	manifestPath := filepath.Join(directory, manifestFile)
	return manifestPath
}

func (c *inputs) setAllText(platform, application string, from, to VersionInfo) {
	c.Platform.Text = platform
	c.Application.Text = application
	c.BuildFrom.Text = from.Build
	c.BuildTo.Text = to.Build
	c.VersionFrom.Text = from.Version
	c.VersionTo.Text = to.Version
	c.refreshAll()
}

func (c *inputs) enableDisableAll(enable bool) {
	if enable {
		c.Platform.Enable()
		c.Application.Enable()
		c.VersionFrom.Enable()
		c.BuildFrom.Enable()
		c.VersionTo.Enable()
		c.BuildTo.Enable()
	} else {
		c.Platform.Disable()
		c.Application.Disable()
		c.VersionFrom.Disable()
		c.BuildFrom.Disable()
		c.VersionTo.Disable()
		c.BuildTo.Disable()
	}
	c.refreshAll()
}

func (c *inputs) refreshAll() {
	c.Platform.Refresh()
	c.Application.Refresh()
	c.VersionFrom.Refresh()
	c.BuildFrom.Refresh()
	c.VersionTo.Refresh()
	c.BuildTo.Refresh()
	// Only refresh API buttons if their state changes (e.g., in setSelectedAPI)
}

func (c *inputs) checkAPISelected() error {
	if c.SelectedAPI == "" {
		return fmt.Errorf("please select an API server")
	}
	return nil
}
