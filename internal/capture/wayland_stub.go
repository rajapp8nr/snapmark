//go:build !linux

package capture

import (
	"fmt"
	"image"
)

func isWayland() bool { return false }

func fullScreenWayland() (image.Image, error) {
	return nil, fmt.Errorf("wayland capture not supported on this platform")
}

func regionWayland() (image.Image, error) {
	return nil, fmt.Errorf("wayland region capture not supported on this platform")
}
