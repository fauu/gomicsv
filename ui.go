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
	"fmt"
	"log"
	"time"

	"github.com/fauu/gomicsv/pixbuf"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func (app *App) uiInit() {
	if err := app.loadWidgets(); err != nil {
		log.Panicf("loading widgets from the definition file: %v", err)
	}

	app.S.Cursor.init()

	app.notificationInit()

	app.menuInit()

	app.preferencesInit()

	app.toolbarInit()

	app.imageAreaInit()

	app.W.MainWindow.SetApplication(app.S.GTKApplication)
	app.W.MainWindow.SetDefaultSize(app.Config.WindowWidth, app.Config.WindowHeight)
	app.W.MainWindow.SetIcon(pixbuf.MustLoad(iconImg))

	var prevW, prevH int
	app.W.MainWindow.Connect("size-allocate", func(_ glib.IObject, allocationPtr uintptr) {
		alloc := gdk.WrapRectangle(allocationPtr)
		w, h := alloc.GetWidth(), alloc.GetHeight()
		if w == prevW && h == prevH {
			return
		}
		prevW, prevH = w, h
		app.handleImageAreaResize()
	})

	app.W.MainWindow.Connect("motion-notify-event", func(_ *gtk.ApplicationWindow, _ *gdk.Event) bool {
		app.S.Cursor.LastMoved = time.Now()
		return true
	})
	glib.TimeoutAdd(250, app.updateCursorVisibility)

	app.W.MainWindow.Connect("key-press-event", func(_ *gtk.ApplicationWindow, event *gdk.Event) {
		ke := &gdk.EventKey{Event: event}
		shift := ke.State()&uint(gdk.SHIFT_MASK) != 0
		ctrl := ke.State()&uint(gdk.CONTROL_MASK) != 0
		app.handleKeyPress(ke.KeyVal(), shift, ctrl)
	})

	app.W.MainWindow.Connect("delete-event", app.quit)

	app.syncWidgetsToConfig()

	app.W.MainWindow.ShowAll()
}

func (app *App) syncWidgetsToConfig() {
	app.W.MenuItemEnlarge.SetActive(app.Config.Enlarge)
	app.W.MenuItemShrink.SetActive(app.Config.Shrink)
	app.W.MenuItemHFlip.SetActive(app.Config.HFlip)
	app.W.MenuItemVFlip.SetActive(app.Config.VFlip)
	app.W.MenuItemRandom.SetActive(app.Config.Random)
	app.W.MenuItemSeamless.SetActive(app.Config.Seamless)
	app.W.MenuItemDoublePage.SetActive(app.Config.DoublePage)
	app.W.MenuItemMangaMode.SetActive(app.Config.MangaMode)

	switch app.Config.ZoomMode {
	case FitToWidth:
		app.W.MenuItemFitToWidth.SetActive(true)
	case FitToHalfWidth:
		app.W.MenuItemFitToHalfWidth.SetActive(true)
	case FitToHeight:
		app.W.MenuItemFitToHeight.SetActive(true)
	case BestFit:
		app.W.MenuItemBestFit.SetActive(true)
	default:
		app.W.MenuItemOriginal.SetActive(true)
	}

	rgba := app.Config.BackgroundColor.ToGdkRGBA()
	app.W.BackgroundColorButton.SetRGBA(&rgba)

	app.W.InterpolationComboBoxText.SetActive(app.Config.Interpolation)
	app.W.OneWideCheckButton.SetActive(app.Config.OneWide)
	app.W.RememberRecentCheckButton.SetActive(app.Config.RememberRecent)
	app.W.RememberPositionCheckButton.SetActive(app.Config.RememberPosition)
	app.W.RememberPositionHTTPCheckButton.SetActive(app.Config.RememberPositionHTTP)
	app.W.RememberPositionHTTPCheckButton.SetSensitive(app.Config.RememberPosition && app.Config.RememberPositionHTTP)
	app.W.EmbeddedOrientationCheckButton.SetActive(app.Config.EmbeddedOrientation)
	app.W.HideIdleCursorCheckButton.SetActive(app.Config.HideIdleCursor)
	app.W.KamiteEnabledCheckButton.SetActive(app.Config.KamiteEnabled)
	app.W.KamitePortContainer.SetSensitive(app.Config.KamiteEnabled)
	app.W.KamitePortEntry.SetText(fmt.Sprint(app.Config.KamitePort))
}

func (app *App) toggleHideUI(hideUI bool) {
	if hideUI {
		app.W.Menubar.Hide()
		app.W.Toolbar.Hide()
		app.W.Statusbar.Hide()
	} else {
		app.W.Menubar.Show()
		app.W.Toolbar.Show()
		app.W.Statusbar.Show()
	}
}

func (app *App) toggleFullscreen(fullscreen bool) {
	if fullscreen {
		app.W.MainWindow.Fullscreen()
	} else {
		app.W.MainWindow.Unfullscreen()
	}
}
