package actions

import "image"

func CopyImage(img image.Image) error {
	return copyImage(img)
}
