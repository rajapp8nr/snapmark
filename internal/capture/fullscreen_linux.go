//go:build linux

package capture

import (
	"fmt"
	"image"
	"os"

	"github.com/kbinani/screenshot"
)

func FullScreen() (image.Image, error) {
	if isWayland() {
		img, err := fullScreenWayland()
		if err == nil {
			return img, nil
		}
	}

	n := screenshot.NumActiveDisplays()
	if n == 0 {
		return nil, fmt.Errorf("no active displays")
	}
	var all image.Rectangle
	for i := 0; i < n; i++ {
		b := screenshot.GetDisplayBounds(i)
		if i == 0 {
			all = b
		} else {
			all = all.Union(b)
		}
	}
	return screenshot.CaptureRect(all)
}

func isWayland() bool {
	return os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland"
}
