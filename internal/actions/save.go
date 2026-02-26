package actions

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func SaveDialog(w fyne.Window, img image.Image) {
	d := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if uc == nil {
			return
		}
		defer uc.Close()
		ext := strings.ToLower(filepath.Ext(uc.URI().Path()))
		var saveErr error
		switch ext {
		case ".png":
			saveErr = png.Encode(uc, img)
		case ".jpg", ".jpeg":
			saveErr = jpeg.Encode(uc, img, &jpeg.Options{Quality: 92})
		default:
			saveErr = fmt.Errorf("unsupported format: use .png or .jpg")
		}
		if saveErr != nil {
			dialog.ShowError(saveErr, w)
		}
	}, w)
	d.SetFileName("snapmark.png")
	d.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	d.Show()
	_ = os.ErrClosed
}
