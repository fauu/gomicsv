/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2024 Piotr Grabowski
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
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/fauu/gomicsv/util"
)

type DragScroll struct {
	InProgress          bool
	StartX              float64
	StartY              float64
	StartHAdjustmentVal float64
	StartVAdjustmentVal float64
}

func (app *App) scroll(dx, dy float64) {
	if !app.archiveIsLoaded() {
		return
	}

	imgw, imgh := app.getImageAreaInnerSize()

	vadj := app.W.ScrolledWindow.GetVAdjustment()
	hadj := app.W.ScrolledWindow.GetHAdjustment()

	vIncrement := vadj.GetMinimumIncrement()
	vVal := vadj.GetValue()
	vMin := vadj.GetLower()
	vMax := vadj.GetUpper() - float64(imgh) - 2

	hIncrement := hadj.GetMinimumIncrement()
	hVal := hadj.GetValue()
	hMin := hadj.GetLower()
	hMax := hadj.GetUpper() - float64(imgw) - 2

	if dy > 0 {
		if vVal >= vMax {
			if app.Config.SmartScroll {
				if app.S.SmartScrollInProgress {
					app.nextPage()
					app.S.SmartScrollInProgress = false
				} else {
					app.S.SmartScrollInProgress = true
				}
			}
		} else {
			vadj.SetValue(util.Clamp(vVal+vIncrement, vMin, vMax))
			app.W.ScrolledWindow.SetVAdjustment(vadj)
			app.S.SmartScrollInProgress = false
		}
	} else if dy < 0 {
		if vVal <= vMin {
			if app.Config.SmartScroll {
				if app.S.SmartScrollInProgress {
					app.previousPage()
					app.S.SmartScrollInProgress = false
				} else {
					app.S.SmartScrollInProgress = true
				}
			}
		} else {
			vadj.SetValue(util.Clamp(vVal-vIncrement, vMin, vMax))
			app.W.ScrolledWindow.SetVAdjustment(vadj)
			app.S.SmartScrollInProgress = false
		}
	}

	if dx > 0 && hVal < hMax {
		hadj.SetValue(util.Clamp(hVal+hIncrement, hMin, hMax))
		app.W.ScrolledWindow.SetHAdjustment(hadj)
	} else if dx < 0 && hVal > hMin {
		hadj.SetValue(util.Clamp(hVal-hIncrement, hMin, hMax))
		app.W.ScrolledWindow.SetHAdjustment(hadj)
	}
}

func (app *App) scrollToStart() {
	app.W.ScrolledWindow.SetVAdjustment(nil)          // Needed to prevent a bug where it scrolls back by itself
	app.W.ScrolledWindow.GetVAdjustment().SetValue(0) // Vertical: top

	var newHadj float64 = 0
	if app.Config.MangaMode {
		imgw, _ := app.getImageAreaInnerSize()
		newHadj = float64(imgw)
	}
	app.W.ScrolledWindow.SetHAdjustment(nil)
	app.W.ScrolledWindow.GetHAdjustment().SetValue(newHadj) // Horizontal: left (non-manga) or right (manga) edge
}

func (app *App) scrollToEnd() {
	imgw, imgh := app.getImageAreaInnerSize()

	app.W.ScrolledWindow.SetVAdjustment(nil)
	app.W.ScrolledWindow.GetVAdjustment().SetValue(float64(imgh)) // Vertical: bottom

	var newHadj float64 = 0
	if !app.Config.MangaMode {
		newHadj = float64(imgw)
	}
	app.W.ScrolledWindow.SetHAdjustment(nil)
	app.W.ScrolledWindow.GetHAdjustment().SetValue(newHadj) // Horizontal: left (manga) or right (non-manga) edge
}

func (app *App) dragScrollStart(x, y, vAdjustmentVal, hAdjustmentVal float64) {
	app.S.DragScroll.InProgress = true
	app.S.DragScroll.StartX = x
	app.S.DragScroll.StartY = y
	app.S.DragScroll.StartVAdjustmentVal = vAdjustmentVal
	app.S.DragScroll.StartHAdjustmentVal = hAdjustmentVal
	app.setCursor(cursorGrabbing)
}

func (app *App) dragScrollUpdate(sw *gtk.ScrolledWindow, event *gdk.EventButton) {
	deltaX := app.S.DragScroll.StartX - event.X()
	sw.GetHAdjustment().SetValue(app.S.DragScroll.StartHAdjustmentVal + deltaX)
	deltaY := app.S.DragScroll.StartY - event.Y()
	sw.GetVAdjustment().SetValue(app.S.DragScroll.StartVAdjustmentVal + deltaY)
}

func (app *App) dragScrollEnd() {
	app.S.DragScroll.InProgress = false
	app.setCursor(cursorDefault)
}
