package editor

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/snapmark/snapmark/internal/actions"
)

type Editor struct {
	app        fyne.App
	win        fyne.Window
	base       *image.RGBA
	current    *image.RGBA
	view       *canvas.Image
	state      ToolState
	ops        []Operation
	history    [][]Operation
	arrowStart *fyne.Position
	dragStart  *fyne.Position
	dragLast   fyne.Position
}

func New(a fyne.App, src image.Image) *Editor {
	b := cloneToRGBA(src)
	e := &Editor{app: a, win: a.NewWindow("SnapMark Editor"), base: b, current: cloneToRGBA(src)}
	e.state = ToolState{Active: ToolRectangle, Color: color.NRGBA{R: 255, A: 255}, Stroke: 3, FontSize: 18}
	e.view = canvas.NewImageFromImage(e.current)
	e.view.FillMode = canvas.ImageFillOriginal

	drawArea := newDrawArea(e)
	content := container.NewBorder(nil, nil, NewToolbar(e), nil, drawArea)
	e.win.SetContent(content)
	e.win.Resize(fyne.NewSize(float32(e.current.Bounds().Dx())+180, float32(e.current.Bounds().Dy())+80))
	e.win.Canvas().AddShortcut(&fyne.ShortcutUndo{}, func(fyne.Shortcut) {
		e.undo()
	})
	return e
}

func (e *Editor) Show() { e.win.Show() }

func (e *Editor) pushHistory() {
	cp := make([]Operation, len(e.ops))
	copy(cp, e.ops)
	e.history = append(e.history, cp)
	if len(e.history) > 10 {
		e.history = e.history[len(e.history)-10:]
	}
}

func (e *Editor) undo() {
	if len(e.history) == 0 {
		return
	}
	e.ops = e.history[len(e.history)-1]
	e.history = e.history[:len(e.history)-1]
	e.render()
}

func (e *Editor) addOp(op Operation) {
	e.pushHistory()
	e.ops = append(e.ops, op)
	e.render()
}

func (e *Editor) render() {
	e.current = e.flatten(false)
	e.view.Image = e.current
	e.view.Refresh()
}

func (e *Editor) save() {
	final := e.flatten(true)
	actions.SaveDialog(e.win, final)
}

func (e *Editor) copyClipboard() {
	if err := actions.CopyImage(e.current); err != nil {
		dialog.ShowError(err, e.win)
		return
	}
	dialog.ShowInformation("Clipboard", "Copied image to clipboard", e.win)
}

func (e *Editor) flatten(bakePixelate bool) *image.RGBA {
	out := cloneToRGBA(e.base)
	for _, op := range e.ops {
		if _, ok := op.(PixelateOp); ok {
			continue
		}
		op.Apply(out)
	}
	for _, op := range e.ops {
		pix, ok := op.(PixelateOp)
		if !ok {
			continue
		}
		if bakePixelate {
			pix.ApplyOverlay(out, out)
		} else {
			pix.ApplyOverlay(out, e.base)
		}
	}
	return out
}

func (e *Editor) renderPixelatePreview(start, end fyne.Position) {
	preview := e.flatten(false)
	x1, y1, x2, y2 := norm(int(start.X), int(start.Y), int(end.X), int(end.Y))
	drawDashedRect(preview, x1, y1, x2, y2)
	e.current = preview
	e.view.Image = e.current
	e.view.Refresh()
}

func drawDashedRect(img *image.RGBA, x1, y1, x2, y2 int) {
	x1, y1, x2, y2 = norm(x1, y1, x2, y2)
	if x2-x1 < 2 || y2-y1 < 2 {
		return
	}
	dash := 6
	for x := x1; x < x2; x++ {
		on := ((x - x1) / dash % 2) == 0
		if on {
			setPixel(img, x, y1, color.White)
			setPixel(img, x, y2-1, color.White)
		} else {
			setPixel(img, x, y1, color.Black)
			setPixel(img, x, y2-1, color.Black)
		}
	}
	for y := y1; y < y2; y++ {
		on := ((y - y1) / dash % 2) == 0
		if on {
			setPixel(img, x1, y, color.White)
			setPixel(img, x2-1, y, color.White)
		} else {
			setPixel(img, x1, y, color.Black)
			setPixel(img, x2-1, y, color.Black)
		}
	}
}

type drawArea struct {
	widget.BaseWidget
	e *Editor
}

func newDrawArea(e *Editor) *drawArea {
	d := &drawArea{e: e}
	d.ExtendBaseWidget(d)
	return d
}
func (d *drawArea) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewMax(d.e.view))
}

func (d *drawArea) Dragged(ev *fyne.DragEvent) {
	if d.e.state.Active != ToolRectangle && d.e.state.Active != ToolEllipse && d.e.state.Active != ToolPixelate {
		return
	}
	if d.e.dragStart == nil {
		p := ev.Position
		d.e.dragStart = &p
	}
	d.e.dragLast = ev.Position
	if d.e.state.Active == ToolPixelate {
		d.e.renderPixelatePreview(*d.e.dragStart, d.e.dragLast)
	}
}

func (d *drawArea) DragEnd() {
	if d.e.dragStart == nil {
		return
	}
	s := *d.e.dragStart
	e := d.e.dragLast
	d.e.dragStart = nil
	if d.e.state.Active == ToolRectangle {
		d.e.addOp(RectOp{X1: int(s.X), Y1: int(s.Y), X2: int(e.X), Y2: int(e.Y), Color: d.e.state.Color, Stroke: d.e.state.Stroke})
	} else if d.e.state.Active == ToolEllipse {
		d.e.addOp(EllipseOp{X1: int(s.X), Y1: int(s.Y), X2: int(e.X), Y2: int(e.Y), Color: d.e.state.Color, Stroke: d.e.state.Stroke})
	} else if d.e.state.Active == ToolPixelate {
		d.e.addOp(PixelateOp{X1: int(s.X), Y1: int(s.Y), X2: int(e.X), Y2: int(e.Y)})
	}
}

func (d *drawArea) Tapped(ev *fyne.PointEvent) {
	s := d.e.state
	switch s.Active {
	case ToolArrow:
		if d.e.arrowStart == nil {
			p := ev.Position
			d.e.arrowStart = &p
			return
		}
		st := *d.e.arrowStart
		d.e.arrowStart = nil
		d.e.addOp(ArrowOp{X1: int(st.X), Y1: int(st.Y), X2: int(ev.Position.X), Y2: int(ev.Position.Y), Color: s.Color, Stroke: s.Stroke})
	case ToolText:
		entry := widget.NewEntry()
		dialog.NewCustomConfirm("Text", "Place", "Cancel", entry, func(ok bool) {
			if ok && entry.Text != "" {
				d.e.addOp(TextOp{X: int(ev.Position.X), Y: int(ev.Position.Y), Text: entry.Text, Color: s.Color, FontSize: s.FontSize})
			}
		}, d.e.win).Show()
	}
}
func (d *drawArea) TappedSecondary(*fyne.PointEvent) {}

func cloneToRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}
