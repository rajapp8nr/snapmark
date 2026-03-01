//go:build linux

package capture

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"strings"
)

func fullScreenWayland() (image.Image, error) {
	cmd := exec.Command("grim", "-")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("wayland capture failed (install grim): %w", err)
	}
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func regionWayland() (image.Image, error) {
	slurpOut, err := exec.Command("slurp").Output()
	if err != nil {
		return nil, fmt.Errorf("region selector failed (install slurp): %w", err)
	}
	geom := strings.TrimSpace(string(slurpOut))
	if geom == "" {
		return nil, fmt.Errorf("empty selection")
	}
	cmd := exec.Command("grim", "-g", geom, "-")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("wayland region capture failed (install grim): %w", err)
	}
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		return nil, err
	}
	return img, nil
}
