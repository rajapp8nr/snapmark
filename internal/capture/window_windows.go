//go:build windows

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
	ps := `$p=Get-Process | Where-Object {$_.MainWindowTitle -ne ''}; foreach($x in $p){$x.MainWindowHandle.ToString() + '|' + $x.MainWindowTitle}`
	out, err := exec.Command("powershell", "-NoProfile", "-Command", ps).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list windows: %w", err)
	}
	wins := make([]WindowInfo, 0)
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		parts := strings.SplitN(s.Text(), "|", 2)
		if len(parts) != 2 {
			continue
		}
		h := strings.TrimSpace(parts[0])
		t := strings.TrimSpace(parts[1])
		if t == "" {
			continue
		}
		r, err := windowRect(h)
		if err != nil || r.Dx() <= 0 || r.Dy() <= 0 {
			continue
		}
		wins = append(wins, WindowInfo{Title: t, Bounds: r, Handle: h})
	}
	return wins, nil
}

func windowRect(handle string) (image.Rectangle, error) {
	ps := fmt.Sprintf(`Add-Type @"
using System;
using System.Runtime.InteropServices;
public class W {
[DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr hWnd, out RECT rect);
public struct RECT { public int Left; public int Top; public int Right; public int Bottom; }
}
"@; $r=New-Object W+RECT; [W]::GetWindowRect([IntPtr]::new(%s), [ref]$r) | Out-Null; Write-Output ($r.Left.ToString()+','+$r.Top+','+$r.Right+','+$r.Bottom)`, handle)
	out, err := exec.Command("powershell", "-NoProfile", "-Command", ps).Output()
	if err != nil {
		return image.Rectangle{}, err
	}
	parts := strings.Split(strings.TrimSpace(string(out)), ",")
	if len(parts) != 4 {
		return image.Rectangle{}, fmt.Errorf("invalid rect")
	}
	l, _ := strconv.Atoi(parts[0])
	t, _ := strconv.Atoi(parts[1])
	r, _ := strconv.Atoi(parts[2])
	b, _ := strconv.Atoi(parts[3])
	return image.Rect(l, t, r, b), nil
}
