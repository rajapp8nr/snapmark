package capture

import (
	"fmt"
	"image"

	"github.com/kbinani/screenshot"
)

func FullScreen() (image.Image, error) {
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
