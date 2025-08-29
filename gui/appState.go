package main

import (
	"fmt"
	"os"
	"path/filepath"

	"rpmg/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type configPaths struct {
	AWSCred               string
	ApiConfigPathOverride string
	PrivateKeyPath        string
	SchemaDir             string
}

type appState struct {
	window       fyne.Window
	inputs       *inputs
	artifacts    *artifacts
	manifestPath string
	forceBox     *widget.Check
	btnSign      *widget.Button
	btnStage     *widget.Button
	btnPublish   *widget.Button
	configPaths  *configPaths
}

func newAppState(w fyne.Window) *appState {
	state := &appState{
		window: w,
	}
	state.forceBox = widget.NewCheck("Force", func(checked bool) {
		// checking will be checked when the commands have to be run, not now
	})
	state.inputs = newInputs(w, state)
	state.artifacts = newArtifacts(w)

	state.btnPublish = widget.NewButton("Publish", func() {
		state.runPublishStep()
	})

	state.btnStage = widget.NewButton("Stage", func() {
		state.runStageStep()
	})

	state.btnSign = widget.NewButton("Sign", func() {
		state.runSignStep()
	})

	state.artifacts.TriggerChange = func() {
		state.enableStagePublish(false)
	}
	return state
}

func (state *appState) enableStagePublish(shouldEnable bool) {
	if shouldEnable {
		state.btnStage.Enable()
		state.btnPublish.Enable()
		return
	}
	state.btnStage.Disable()
	state.btnPublish.Disable()
}

func (state *appState) setManifestPath(path string) {
	state.manifestPath = path
	state.inputs.ManifestPath.SetPath(path)
	//state.artifacts.updateArtifactManager(path)
}

func (state *appState) changeManifestDir() {
	currentPath := state.manifestPath

	content := container.NewVBox(
		widget.NewLabel("Current Manifest Path:"),
		widget.NewLabelWithStyle(currentPath, fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		widget.NewLabel("\nDo you want to choose a new path or keep the current one?"),
	)

	dialog.ShowCustomConfirm(
		"Select Manifest Directory",
		"Choose New Path",
		"Keep Current",
		content,
		func(chooseNew bool) {
			if chooseNew {
				state.selectAndUpdateManifestDir(currentPath)
			} else {
				state.forceBox.SetChecked(false)
			}
		},
		state.window,
	)
}

func (state *appState) selectAndUpdateManifestDir(currentPath string) {
	dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, state.window)
			return
		}
		if uri != nil {

			chosenPath := filepath.Join(filepath.FromSlash(uri.Path()), filepath.Base(state.manifestPath))
			err := os.Rename(currentPath, chosenPath)
			if err != nil {
				dialog.ShowError(err, state.window)
				return
			}
			state.setManifestPath(chosenPath)
		}
	}, state.window).Show()
}

func (state *appState) runSignStep() {
	if !state.artifacts.validateArtifactPath(state.window) {
		return
	}
	if !utils.FileExists(state.inputs.ManifestPath.Path) {
		directory, err := os.Getwd()
		if err != nil {
			dialog.ShowError(err, state.window)
			return
		}
		path, err := state.inputs.createManifest(directory)
		if err != nil {
			dialog.ShowError(err, state.window)
			return
		}
		if !utils.FileExists(state.manifestPath) {
			state.manifestPath = path
			state.runCreateStep()
			return
		}
		state.setManifestPath(path)
		state.inputs.enableDisableAll(false)
	}

	err := state.inputs.checkAPISelected()
	if err != nil {
		dialog.ShowError(err, state.window)
		return
	}

	err = state.runExternalCommand("sign")
	if err != nil {
		dialog.ShowError(err, state.window)
		return
	}
	state.enableStagePublish(true)
	state.changeManifestDir()
	// fmt.Println(state.manifestPath) // Debug: print manifest path if needed
}

func (state *appState) runCreateStep() {
	err := state.runExternalCommand("create")
	if err != nil {
		dialog.ShowError(err, state.window)
		state.setManifestPath("")
		return
	}
	state.setManifestPath(state.manifestPath)
	state.inputs.enableDisableAll(false)
}

func (state *appState) runStageStep() {
	_, err := state.inputs.validateSignInputs()
	if err != nil {
		dialog.ShowError(err, state.window)
		return
	}
	err = state.runExternalCommand("stage")
	if err != nil {
		dialog.ShowError(err, state.window)
		return
	}
	state.forceBox.SetChecked(false)
}

func (state *appState) runPublishStep() {
	_, err := state.inputs.validateSignInputs()
	if err != nil {
		dialog.ShowError(err, state.window)
		return
	}
	err = state.runExternalCommand("publish")
	if err != nil {
		dialog.ShowError(err, state.window)
		return
	}
	state.forceBox.SetChecked(false)
}

func (state *appState) buildUI() fyne.CanvasObject {
	return state.inputs.buildInputUI()
}

func (state *appState) runExternalCommand(step string) error {
	w := state.window

	var message string
	var err error

	//artifactPath := state.artifacts.ArtifactPath.Path
	//config := rpm.NewRpmConfig(state.manifestPath, artifactPath, state.inputs.SelectedAPI)

	// Set configuration paths

	/*
		if state.configPaths.ApiConfigPathOverride != "" {
			config.ApiConfigPath = state.configPaths.ApiConfigPathOverride
		}

		config.AWSCred = state.configPaths.AWSCred
		config.PrivateKeyPath = state.configPaths.PrivateKeyPath
		config.SchemaDir = state.configPaths.SchemaDir
	*/

	//boxChecked := state.forceBox.Checked
	switch step {
	case "create":

		message, err = "The create step would have been run here", nil

	case "sign":

		/*
			for _, a := range *state.artifacts.ArtifactManager.Artifacts {
				fmt.Println("Artifact: ", a.ArtifactProperties.Name, ", Order Number:", a.InstallOrder, ", Install Method:", a.InstallMethod.Type)
			}

			var manifest rpm.Manifest

			manifest, err = rpm.Manifest{}.New(state.manifestPath, *state.artifacts.ArtifactManager.Artifacts)
			if err != nil {
				message = ""
				break
			}
			manifest.Save(state.manifestPath)
		*/
		message, err = "The sign step would have been run", nil

	case "stage":

		message, err = "The stage step would have been run", nil

	case "publish":

		message, err = "The publish step would have been run", nil

	default:
		return fmt.Errorf("unknown step %s", step)
	}

	if err != nil {
		return fmt.Errorf("error running %s step: %v", step, err)
	}

	if message != "" {
		dialog.ShowInformation("Success", fmt.Sprintf("%s step completed successfully!\n\nOutput:\n%s", step, message), w)
	}
	return nil
}
