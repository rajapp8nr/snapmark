//go:build darwin

package capture

import (
	"encoding/json"
	"fmt"
	"image"
	"os/exec"
)

type cgWindow struct {
	Name   string `json:"kCGWindowName"`
	Bounds struct {
		X      float64 `json:"X"`
		Y      float64 `json:"Y"`
		Width  float64 `json:"Width"`
		Height float64 `json:"Height"`
	} `json:"kCGWindowBounds"`
}

func ListWindows() ([]WindowInfo, error) {
	script := `osascript -l JavaScript <<'JXA'
ObjC.import('CoreGraphics');
const opts = $.kCGWindowListOptionOnScreenOnly | $.kCGWindowListExcludeDesktopElements;
const list = $.CGWindowListCopyWindowInfo(opts, $.kCGNullWindowID);
console.log(JSON.stringify(ObjC.deepUnwrap(list)));
JXA`
	out, err := exec.Command("bash", "-lc", script).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query windows: %w", err)
	}
	var raw []cgWindow
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, err
	}
	res := make([]WindowInfo, 0)
	for _, w := range raw {
		if w.Name == "" || w.Bounds.Width < 40 || w.Bounds.Height < 40 {
			continue
		}
		r := image.Rect(int(w.Bounds.X), int(w.Bounds.Y), int(w.Bounds.X+w.Bounds.Width), int(w.Bounds.Y+w.Bounds.Height))
		res = append(res, WindowInfo{Title: w.Name, Bounds: r})
	}
	return res, nil
}
