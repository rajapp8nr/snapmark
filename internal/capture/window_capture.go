package capture

import (
	"image"

	"github.com/kbinani/screenshot"
)

func Window(win WindowInfo) (image.Image, error) {
	return screenshot.CaptureRect(win.Bounds)
}
