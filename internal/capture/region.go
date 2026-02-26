package capture

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
)

type regionSelector struct {
	widget.BaseWidget
	start, end fyne.Position
	dragging   bool
	onDone     func(image.Rectangle)
}

func newRegionSelector(onDone func(image.Rectangle)) *regionSelector {
	r := &regionSelector{onDone: onDone}
	r.ExtendBaseWidget(r)
	return r
}

func (r *regionSelector) Dragged(e *fyne.DragEvent) {
	if !r.dragging {
		r.start = e.Position
		r.dragging = true
	}
	r.end = e.Position
	r.Refresh()
}

func (r *regionSelector) DragEnd() {
	if !r.dragging {
		return
	}
	r.dragging = false
	rect := normalizeRect(r.start, r.end)
	r.onDone(rect)
}

func (r *regionSelector) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 96})
	selection := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 32})
	border := canvas.NewRectangle(color.NRGBA{R: 255, G: 64, B: 64, A: 220})
	border.StrokeColor = color.NRGBA{R: 255, G: 64, B: 64, A: 255}
	border.StrokeWidth = 2

	objs := []fyne.CanvasObject{bg, selection, border}
	return &regionRenderer{selector: r, bg: bg, selection: selection, border: border, objects: objs}
}

type regionRenderer struct {
	selector                *regionSelector
	bg, selection, border   *canvas.Rectangle
	objects                 []fyne.CanvasObject
}

func (r *regionRenderer) Layout(sz fyne.Size) {
	r.bg.Resize(sz)
	rect := normalizeRect(r.selector.start, r.selector.end)
	if rect.Empty() {
		r.selection.Hide()
		r.border.Hide()
		return
	}
	r.selection.Show()
	r.border.Show()
	p := fyne.NewPos(float32(rect.Min.X), float32(rect.Min.Y))
	s := fyne.NewSize(float32(rect.Dx()), float32(rect.Dy()))
	r.selection.Move(p)
	r.selection.Resize(s)
	r.border.Move(p)
	r.border.Resize(s)
}
func (r *regionRenderer) MinSize() fyne.Size { return fyne.NewSize(100, 100) }
func (r *regionRenderer) Refresh()            { r.Layout(r.selector.Size()) }
func (r *regionRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *regionRenderer) Destroy() {}

func normalizeRect(a, b fyne.Position) image.Rectangle {
	x1, y1 := int(a.X), int(a.Y)
	x2, y2 := int(b.X), int(b.Y)
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	return image.Rect(x1, y1, x2, y2)
}

func Region(app fyne.App) (image.Image, error) {
	w := app.NewWindow("Select Region")
	w.SetPadded(false)
	w.SetFullScreen(true)

	var (
		once sync.Once
		out  image.Image
		err  error
		done = make(chan struct{})
	)

	selector := newRegionSelector(func(rect image.Rectangle) {
		once.Do(func() {
			if rect.Dx() < 2 || rect.Dy() < 2 {
				err = fmt.Errorf("selection too small")
			} else {
				out, err = screenshot.CaptureRect(rect)
			}
			close(done)
			w.Close()
		})
	})

	w.SetContent(container.NewMax(selector))
	w.Show()
	<-done
	return out, err
}
