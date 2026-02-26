package main

import (
	"fmt"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/snapmark/snapmark/internal/capture"
	"github.com/snapmark/snapmark/internal/editor"
)

func main() {
	a := app.NewWithID("com.snapmark.app")
	w := a.NewWindow("SnapMark")
	w.Resize(fyne.NewSize(420, 240))

	status := widget.NewLabel("Pick a capture mode")

	openEditor := func(img image.Image) {
		ed := editor.New(a, img)
		ed.Show()
	}

	fullBtn := widget.NewButton("Full Screen", func() {
		status.SetText("Capturing full screen in 2 seconds...")
		go func() {
			time.Sleep(2 * time.Second)
			img, err := capture.FullScreen()
			fyne.Do(func() {
				if err != nil {
					dialog.ShowError(err, w)
					status.SetText("Capture failed")
					return
				}
				status.SetText("Captured")
				openEditor(img)
			})
		}()
	})

	regionBtn := widget.NewButton("Select Region", func() {
		img, err := capture.Region(a)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		openEditor(img)
	})

	windowBtn := widget.NewButton("Select Window", func() {
		wins, err := capture.ListWindows()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if len(wins) == 0 {
			dialog.ShowInformation("No windows", "No capturable windows found.", w)
			return
		}

		options := make([]string, len(wins))
		idxByTitle := map[string]int{}
		for i, win := range wins {
			title := fmt.Sprintf("%s (%dx%d)", win.Title, win.Bounds.Dx(), win.Bounds.Dy())
			options[i] = title
			idxByTitle[title] = i
		}
		selected := options[0]
		selectWidget := widget.NewSelect(options, func(v string) { selected = v })
		selectWidget.SetSelected(selected)

		d := dialog.NewCustomConfirm("Select Window", "Capture", "Cancel", selectWidget, func(ok bool) {
			if !ok {
				return
			}
			idx := idxByTitle[selected]
			img, err := capture.Window(wins[idx])
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			openEditor(img)
		}, w)
		d.Show()
	})

	w.SetContent(container.NewVBox(
		widget.NewRichTextFromMarkdown("# SnapMark\nCross-platform screenshot annotation"),
		fullBtn,
		regionBtn,
		windowBtn,
		status,
	))

	w.ShowAndRun()
}
