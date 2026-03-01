//go:build linux

package capture

import (
	"fmt"
	"image"

	"github.com/kbinani/screenshot"
)

func Window(win WindowInfo) (image.Image, error) {
	if isWayland() {
		return nil, fmt.Errorf("window capture is not supported on Wayland; use Select Region")
	}
	return screenshot.CaptureRect(win.Bounds)
}
