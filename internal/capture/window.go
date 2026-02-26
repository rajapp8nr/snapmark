package capture

import "image"

type WindowInfo struct {
	Title  string
	Bounds image.Rectangle
	Handle string
}
