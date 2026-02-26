//go:build linux

package actions

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os/exec"
)

func copyImage(img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil { return err }
	if err := runPipe("xclip", []string{"-selection", "clipboard", "-t", "image/png", "-i"}, buf.Bytes()); err == nil {
		return nil
	}
	if err := runPipe("xsel", []string{"--clipboard", "--input", "--mime-type", "image/png"}, buf.Bytes()); err == nil {
		return nil
	}
	return fmt.Errorf("clipboard tool not found: install xclip or xsel")
}

func runPipe(name string, args []string, data []byte) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}
