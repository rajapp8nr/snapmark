//go:build darwin

package actions

import (
	"bytes"
	"image"
	"image/png"
	"os/exec"
)

func copyImage(img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil { return err }
	cmd := exec.Command("osascript", "-e", `set the clipboard to (read (POSIX file \"/dev/stdin\") as «class PNGf»)`)
	cmd.Stdin = bytes.NewReader(buf.Bytes())
	return cmd.Run()
}
