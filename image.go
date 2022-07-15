/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2022 Piotr Grabowski
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package gomicsv

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/fauu/gomicsv/util"
)

const savedImagePNGQuality = 9

var interpolations = []gdk.InterpType{gdk.INTERP_NEAREST, gdk.INTERP_TILES, gdk.INTERP_BILINEAR, gdk.INTERP_HYPER}

func (app *App) pixbufLoaded() bool {
	if app.Config.DoublePage && !app.shouldForceSinglePage() {
		return app.S.PixbufL != nil && app.S.PixbufR != nil
	}
	return app.S.PixbufL != nil
}

func (app *App) pixbufSize() (w, h int) {
	if !app.pixbufLoaded() {
		return 0, 0
	}

	s := &app.S

	if app.Config.DoublePage && !app.shouldForceSinglePage() {
		return s.PixbufL.GetWidth() + s.PixbufR.GetWidth(), util.Max(s.PixbufL.GetHeight(), s.PixbufR.GetHeight())
	}
	return s.PixbufL.GetWidth(), s.PixbufL.GetHeight()
}

func (app *App) updateStatus() {
	s := &app.S

	if !app.pixbufLoaded() {
		return
	}

	zoom := int(100 * app.S.Scale)

	lenStr := "?"
	if s.Archive.Len() != nil {
		lenStr = fmt.Sprint(*s.Archive.Len())
	}

	markedStr := ""
	if app.currentPageIsJumpmarked() {
		markedStr = " MARKED "
	}

	var msg, title string
	if app.Config.DoublePage && !app.shouldForceSinglePage() {
		leftPath, _ := s.Archive.Name(s.ArchivePos)
		left := filepath.Base(leftPath)
		rightPath, _ := s.Archive.Name(s.ArchivePos + 1)
		right := filepath.Base(rightPath)

		leftIndex := s.ArchivePos + 1
		rightIndex := s.ArchivePos + 2

		leftw, lefth := s.PixbufL.GetWidth(), s.PixbufL.GetHeight()
		rightw, righth := s.PixbufR.GetWidth(), s.PixbufR.GetHeight()

		if app.Config.MangaMode {
			left, right = right, left
			leftIndex, rightIndex = rightIndex, leftIndex
			leftw, rightw = rightw, leftw
		}
		msg = fmt.Sprintf("%d+%d / %s %s |   %dx%d - %dx%d (%d%%)   |   %s   |   %s - %s", leftIndex, rightIndex, lenStr, markedStr, leftw, lefth, rightw, righth, zoom, s.ArchiveName, left, right)
		title = fmt.Sprintf("[%d+%d / %s] %s", leftIndex, rightIndex, lenStr, s.ArchiveName)
	} else {
		imgPath, _ := s.Archive.Name(s.ArchivePos)
		w, h := s.PixbufL.GetWidth(), s.PixbufL.GetHeight()
		msg = fmt.Sprintf("%d / %s %s  |   %dx%d (%d%%)   |   %s   |   %s", s.ArchivePos+1, lenStr, markedStr, w, h, zoom, s.ArchiveName, imgPath)
		title = fmt.Sprintf("[%d / %s] %s", s.ArchivePos+1, lenStr, s.ArchiveName)
	}
	app.setStatus(msg)

	app.W.MainWindow.SetTitle(title)
}

func (app *App) getImageAreaInnerSize() (width, height int) {
	alloc := app.W.ScrolledWindow.GetAllocation()
	return alloc.GetWidth() - 4, alloc.GetHeight() - 4 // 2u of padding from the ScrolledWindow and 2u from the Viewport
}

func (app *App) getScaledSize() (scale float64) {
	if !app.pixbufLoaded() {
		return
	}

	scrw, scrh := app.getImageAreaInnerSize()

	w, h := app.pixbufSize()
	switch app.Config.ZoomMode {
	case FitToWidth:
		needscale := (app.Config.Enlarge && w < scrw) || (app.Config.Shrink && w > scrw)
		if needscale {
			return float64(scrw) / float64(w)
		}
	case FitToHalfWidth:
		needscale := (app.Config.Enlarge && w/2 < scrw) || (app.Config.Shrink && w/2 > scrw)
		if needscale {
			return float64(scrw) / float64(w/2)
		}
	case FitToHeight:
		return float64(scrh) / float64(h)
	case BestFit:
		needscale := (app.Config.Enlarge && (w < scrw && h < scrh)) || (app.Config.Shrink && (w > scrw || h > scrh))
		if needscale {
			fw, _ := util.Fit(w, h, scrw, scrh)
			return float64(fw) / float64(w)
		}
	}
	return 1
}

func (app *App) shouldForceSinglePage() bool {
	if app.S.PixbufR == nil {
		return true
	}
	return app.Config.OneWide && (app.S.PixbufL.GetWidth() > app.S.PixbufL.GetHeight() || app.S.PixbufR.GetWidth() > app.S.PixbufR.GetHeight())
}

func (app *App) blit() {
	if !app.pixbufLoaded() {
		return
	}

	app.S.Scale = app.getScaledSize()

	// Check whether the scale of the left image is different from the old one?

	if app.Config.DoublePage && !app.shouldForceSinglePage() {
		left := app.S.PixbufL
		right := app.S.PixbufR

		if app.Config.MangaMode {
			left, right = right, left
		}

		if err := app.doBlit(app.W.ImageL, left, app.S.Scale); err != nil {
			app.showError(err.Error())
			return
		}

		if err := app.doBlit(app.W.ImageR, right, app.S.Scale); err != nil {
			app.showError(err.Error())
			return
		}
	} else {
		app.W.ImageR.Clear()
		if err := app.doBlit(app.W.ImageL, app.S.PixbufL, app.S.Scale); err != nil {
			app.showError(err.Error())
			return
		}
	}

	if app.S.Scale != 1 || app.Config.HFlip || app.Config.VFlip {
		util.GC()
	}
}

func (app *App) doBlit(image *gtk.Image, pixbuf *gdk.Pixbuf, scale float64) (err error) {
	image.Clear()

	if app.Config.HFlip {
		pixbuf, err = pixbuf.Flip(true)
		if err != nil {
			return err
		}
	}

	if app.Config.VFlip {
		pixbuf, err = pixbuf.Flip(false)
		if err != nil {
			return err
		}
	}

	if scale != 1 {
		w, h := pixbuf.GetWidth(), pixbuf.GetHeight()
		pixbuf, err = pixbuf.ScaleSimple(int(float64(w)*scale), int(float64(h)*scale), interpolations[app.Config.Interpolation])
		if err != nil {
			return err
		}
	}

	image.SetFromPixbuf(pixbuf)

	return nil
}

const saveImageErrorMsgTpl = "Couldn't save image: %v"

func (app *App) saveImage(path string) {
	if app.S.PixbufR == nil {
		if app.S.PixbufL == nil {
			return
		} else {
			if err := app.S.PixbufL.SavePNG(path, savedImagePNGQuality); err != nil {
				app.showError(fmt.Sprintf(saveImageErrorMsgTpl, err))
				return
			}
		}
	}

	// We know we're in double page mode
	stichedPixbuf, err := app.getStichedPixbuf()
	if err != nil {
		app.showError(fmt.Sprintf(saveImageErrorMsgTpl, err))
		return
	}
	if err := stichedPixbuf.SavePNG(path, savedImagePNGQuality); err != nil {
		app.showError(fmt.Sprintf(saveImageErrorMsgTpl, err))
		return
	}

	app.notificationShow(fmt.Sprintf("Saved to %s", path), ShortNotification)
}

const copyImageSuccessMsg = "Copied image to clipboard"

func (app *App) copyImageToClipboard() {
	clipboard, err := gtk.ClipboardGet(gdk.GdkAtomIntern("CLIPBOARD", true))
	if err != nil {
		log.Panicf("getting clipboard: %v", err)
	}

	if app.S.PixbufR == nil {
		if app.S.PixbufL != nil {
			clipboard.SetImage(app.S.PixbufL)
			app.notificationShow(copyImageSuccessMsg, ShortNotification)
		}
		return
	}

	// We know we're in double page mode
	stichedPixbuf, err := app.getStichedPixbuf()
	if err != nil {
		log.Printf("Error stiching images: %v", err)
		app.showError("Couldn't copy image to clipboard")
		return
	}

	clipboard.SetImage(stichedPixbuf)
	app.notificationShow(copyImageSuccessMsg, ShortNotification)
}

// getStichedPixbuf creates a Pixbuf combining the left and right image Pixbufs
func (app *App) getStichedPixbuf() (*gdk.Pixbuf, error) {
	l, r := app.S.PixbufL, app.S.PixbufR
	if app.Config.MangaMode {
		l, r = r, l
	}

	// ASSUMPTION: Both Pixbufs have the same parameters except for size
	mergedPixbuf, err := gdk.PixbufNew(
		l.GetColorspace(),
		l.GetHasAlpha(),
		l.GetBitsPerSample(),
		l.GetWidth()+r.GetWidth(),
		util.Max(l.GetHeight(), r.GetHeight()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating pixbuf: %v", err)
	}

	mergedPixbuf.Fill(app.Config.BackgroundColor.ToPixelInt())

	// POLISH: If heights are unequal, center the smaller image vertically, just as it appears in the program
	pixbufCopyPixels(l, mergedPixbuf, 0)
	pixbufCopyPixels(r, mergedPixbuf, l.GetWidth())

	return mergedPixbuf, nil
}
