//go:build windows

package actions

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"os/exec"
)

func copyImage(img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil { return err }
	tmp, err := os.CreateTemp("", "snapmark-*.png")
	if err != nil { return err }
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(buf.Bytes()); err != nil { return err }
	_ = tmp.Close()
	ps := `Add-Type -AssemblyName System.Windows.Forms; Add-Type -AssemblyName System.Drawing; $img=[System.Drawing.Image]::FromFile('` + tmp.Name() + `'); [System.Windows.Forms.Clipboard]::SetImage($img)`
	return exec.Command("powershell", "-NoProfile", "-Command", ps).Run()
}
