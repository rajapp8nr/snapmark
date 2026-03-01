//go:build linux

package capture

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"os/exec"
	"strconv"
	"strings"
)

func ListWindows() ([]WindowInfo, error) {
	if isWayland() {
		return nil, fmt.Errorf("window capture list is not supported on Wayland; use Select Region")
	}
	cmd := exec.Command("bash", "-lc", "wmctrl -lG")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("wmctrl required for window listing on linux: %w", err)
	}
	windows := make([]WindowInfo, 0)
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}
		x, _ := strconv.Atoi(fields[2])
		y, _ := strconv.Atoi(fields[3])
		w, _ := strconv.Atoi(fields[4])
		h, _ := strconv.Atoi(fields[5])
		title := strings.Join(fields[7:], " ")
		if title == "" || w <= 0 || h <= 0 {
			continue
		}
		windows = append(windows, WindowInfo{Title: title, Bounds: image.Rect(x, y, x+w, y+h), Handle: fields[0]})
	}
	return windows, nil
}
