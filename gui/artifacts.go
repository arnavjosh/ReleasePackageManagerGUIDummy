package main

import (
	"fmt"
	"image/color"

	"rpmg/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

type artifacts struct {
	window       fyne.Window
	ArtifactPath *utils.FileChooserEntry
	//ArtifactManager *utils.DraggableManager
	TriggerChange func()
}

func newArtifacts(w fyne.Window) *artifacts {
	arts := &artifacts{
		window:       w,
		ArtifactPath: utils.NewFileChooserEntry(w, "/.."),
	}

	arts.ArtifactPath.Spec = (utils.FileSpecDirectory)
	arts.ArtifactPath.Callback = (func(lastPath string) {
		// some add button stuff
	})
	return arts
}

func (arts *artifacts) buildArtifactPathUI() fyne.CanvasObject {
	header := canvas.NewText("Artifacts", color.White)
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.Alignment = fyne.TextAlignLeading
	header.TextSize = 18

	//arts.ArtifactManager = utils.NewDraggableManager(arts.window, arts.TriggerChange)

	return container.NewVBox(
		header,
		arts.ArtifactPath,
		//arts.ArtifactManager,
		//arts.makeAddArtifactButton(),
	)
}

func (arts *artifacts) getArtifactPath() string {
	return arts.ArtifactPath.Path
}

func (arts *artifacts) validateArtifactPath(w fyne.Window) bool {
	if arts.getArtifactPath() == "" {
		err := fmt.Errorf("artifact Path field is required")
		dialog.ShowError(err, w)
		return false
	}
	return true
}

/*
func (arts *artifacts) updateArtifactManager(manifestPath string) error {
	manifest, err := rpm.LoadManifest(manifestPath, arts.ArtifactPath.Path)
	if err != nil {
		return fmt.Errorf("error loading manifest for artifacts: %w", err)
	}
	arts.ArtifactManager.UpdateItems(&manifest.Artifacts)
	return nil
}
*/

/*
func (arts *artifacts) makeAddArtifactButton() *fyne.Container {
	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if !arts.validateArtifactPath(arts.window) {
			return
		}

		available, err := arts.getAvailableArtifacts()
		if err != nil {
			dialog.ShowError(err, arts.window)
			return
		}

		if len(available) == 0 {
			dialog.ShowInformation("No New Artifacts", "No new artifacts found to add.", arts.window)
			return
		}

		arts.showAddArtifactDialog(available)
	})

	return container.NewHBox(layout.NewSpacer(), addButton)
}

func (arts *artifacts) getAvailableArtifacts() ([]rpm.Artifact, error) {
	allArtifacts, err := rpm.ScanArtifacts(arts.getArtifactPath())
	if err != nil {
		return nil, fmt.Errorf("failed to scan artifacts: %w", err)
	}

	existing := make(map[string]bool)
	for _, a := range *arts.ArtifactManager.Artifacts {
		existing[a.ArtifactProperties.Name] = true
	}

	var available []rpm.Artifact
	for _, artifact := range allArtifacts {
		name := artifact.ArtifactProperties.Name
		if !existing[name] {
			available = append(available, artifact)
		}
	}

	return available, nil
}

func (arts *artifacts) showAddArtifactDialog(available []rpm.Artifact) {
	checkMap := make(map[string]*widget.Check)
	checkboxes := make([]fyne.CanvasObject, 0, len(available))

	for _, artifact := range available {
		name := artifact.ArtifactProperties.Name
		check := widget.NewCheck(name, nil)
		checkMap[name] = check
		checkboxes = append(checkboxes, check)
	}

	content := container.NewVBox(
		widget.NewLabel("Select artifacts to add:"),
		container.NewVBox(checkboxes...),
	)

	dialog.ShowCustomConfirm("Add Artifacts", "Add", "Cancel", content, func(confirmed bool) {
		if !confirmed {
			return
		}

		if arts.ArtifactManager.Artifacts == nil {
			empty := make([]rpm.Artifact, 0)
			arts.ArtifactManager.Artifacts = &empty
		}

		currentArtifacts := *arts.ArtifactManager.Artifacts
		updatedArtifacts := currentArtifacts

		addedAny := false
		for _, artifact := range available {
			name := artifact.ArtifactProperties.Name
			if checkMap[name].Checked {
				updatedArtifacts = append(updatedArtifacts, artifact)
				addedAny = true
			}
		}

		if addedAny {
			arts.ArtifactManager.UpdateItems(&updatedArtifacts)

			if arts.TriggerChange != nil {
				arts.TriggerChange()
			}
		}
	}, arts.window)
}
*/
