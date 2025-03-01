/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2025 Piotr Grabowski
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
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const kamiteLongRightClickDelay = 375

func (app *App) imageAreaInit() {
	app.W.ScrolledWindow.SetEvents(app.W.ScrolledWindow.GetEvents() | int(gdk.BUTTON_PRESS_MASK) | int(gdk.BUTTON_RELEASE_MASK))

	app.W.ScrolledWindow.Connect("scroll-event", func(self *gtk.ScrolledWindow, event *gdk.Event) {
		se := &gdk.EventScroll{Event: event}
		app.scroll(se.DeltaX(), se.DeltaY())
	})

	app.W.ScrolledWindow.Connect("button-press-event", func(self *gtk.ScrolledWindow, event *gdk.Event) bool {
		if app.S.UITemporarilyRevealed {
			app.S.UITemporarilyRevealed = false
			app.toggleHideUI(true)
			return true
		}

		be := &gdk.EventButton{Event: event}
		switch be.Button() {
		case 1:
			if (int)(be.X()) < self.GetAllocatedWidth()/2 {
				app.pageLeft()
			} else {
				app.pageRight()
			}
		case 3:
			if app.Config.KamiteEnabled {
				app.S.KamiteRightClickActionPending = true
				glib.TimeoutAdd(kamiteLongRightClickDelay, func() bool {
					if app.S.KamiteRightClickActionPending {
						app.kamiteRecognizeManualBlock()
						app.S.KamiteRightClickActionPending = false
					}
					return false
				})
			}
			// Empty
		case 2:
			app.dragScrollStart(be.X(), be.Y(), self.GetVAdjustment().GetValue(), self.GetHAdjustment().GetValue())
		}
		return true
	})

	app.W.ScrolledWindow.Connect("button-release-event", func(_ *gtk.ScrolledWindow, event *gdk.Event) bool {
		be := &gdk.EventButton{Event: event}
		switch be.Button() {
		case 1:
			// Empty
		case 3:
			if app.Config.KamiteEnabled {
				if app.S.KamiteRightClickActionPending {
					app.kamiteRecognizeImageUnderCursorBlock()
					app.S.KamiteRightClickActionPending = false
				}
			}
		case 2:
			app.dragScrollEnd()
		}
		return true
	})

	app.W.ScrolledWindow.Connect("motion-notify-event", func(sw *gtk.ScrolledWindow, event *gdk.Event) bool {
		if app.S.DragScroll.InProgress {
			app.dragScrollUpdate(sw, &gdk.EventButton{Event: event})
		}
		return false // Let it be handled for MainWindow
	})

}

func (app *App) handleImageAreaResize() {
	app.blit()
	app.updateStatus()
}
