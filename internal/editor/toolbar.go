package editor

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewToolbar(e *Editor) fyne.CanvasObject {
	toolSelect := widget.NewSelect([]string{"Rectangle", "Ellipse", "Arrow", "Text", "Pixelate"}, func(v string) {
		switch v {
		case "Rectangle":
			e.state.Active = ToolRectangle
		case "Ellipse":
			e.state.Active = ToolEllipse
		case "Arrow":
			e.state.Active = ToolArrow
		case "Text":
			e.state.Active = ToolText
		case "Pixelate":
			e.state.Active = ToolPixelate
		}
	})
	toolSelect.SetSelected("Rectangle")

	stroke := widget.NewSelect([]string{"1", "2", "3", "4", "6", "8"}, func(v string) {
		n, _ := strconv.Atoi(v)
		e.state.Stroke = n
	})
	stroke.SetSelected("3")

	font := widget.NewSelect([]string{"12", "16", "18", "24", "32"}, func(v string) {
		n, _ := strconv.Atoi(v)
		e.state.FontSize = n
	})
	font.SetSelected("18")

	colorBtn := widget.NewButton("Colour", func() {
		dialog.ShowColorPicker("Pick colour", "", func(c color.Color) { e.state.Color = c }, e.win)
	})

	undoBtn := widget.NewButton("Undo (Ctrl+Z)", func() { e.undo() })
	saveBtn := widget.NewButton("Save As", func() { e.save() })
	copyBtn := widget.NewButton("Copy", func() { e.copyClipboard() })

	return container.NewVBox(toolSelect, stroke, font, colorBtn, undoBtn, saveBtn, copyBtn)
}
