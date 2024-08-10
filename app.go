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
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/fauu/gomicsv/archive"
	"github.com/fauu/gomicsv/imgdiff"
	"github.com/fauu/gomicsv/pagecache"
	"github.com/fauu/gomicsv/util"
)

const (
	AppName        = "gomicsv"
	AppNameDisplay = "Gomics-v"
	AppID          = "com.github.fauu.gomicsv"
)

type App struct {
	S      State
	W      Widgets
	Config Config
}

type State struct {
	BuildInfo                     BuildInfo
	GTKApplication                *gtk.Application
	Archive                       archive.Archive
	ArchivePos                    int
	ArchivePath                   string
	PixbufL, PixbufR              *gdk.Pixbuf
	GoToThumbPixbuf               *gdk.Pixbuf
	Scale                         float64
	PageCache                     *pagecache.PageCache
	ConfigDirPath                 string
	UserDataDirPath               string
	ReadLaterDirPath              string
	ImageHashes                   map[int]imgdiff.Hash
	Jumpmarks                     Jumpmarks
	Cursor                        CursorsState
	DragScroll                    DragScroll
	SmartScrollInProgress         bool
	KamiteRightClickActionPending bool
	RecentManager                 *gtk.RecentManager
	BackgroundColorCssProvider    *gtk.CssProvider
	PageCacheTrimTimeoutHandle    *glib.SourceHandle
	UITemporarilyRevealed         bool
}

//go:embed about.jpg
var aboutImg []byte

//go:embed icon.png
var iconImg []byte

//go:embed gomicsv.ui
var uiDef string

type AppStartupParams struct {
	Referer string
}

func (app *App) Init(nonFlagArgs []string, startupParams AppStartupParams, buildInfo BuildInfo) *gtk.Application {
	application, err := gtk.ApplicationNew(AppID, glib.APPLICATION_HANDLES_OPEN)
	if err != nil {
		log.Panicf("creating GTK Aplication: %v", err)
	}
	glib.SetPrgname(AppID)

	app.S.BuildInfo = buildInfo

	application.Connect("startup", func(self *gtk.Application) {
		app.ensureDirs()

		app.loadConfig()

		app.S.RecentManager, err = gtk.RecentManagerGetDefault()
		if err != nil {
			log.Panicf("getting default RecentManager: %v", err)
		}

		app.uiInit()

		app.syncStateToConfig()
	})

	// Runs after startup if no files were provided
	application.Connect("activate", func(_ *gtk.Application) {
		// Do nothing
	})

	// Runs after startup if some files were provided
	application.Connect("open", func(_ *gtk.Application, filesPtr unsafe.Pointer, _count int, _ string) {
		// Get first file's path and load it, ignoring the rest
		offset := uintptr(0) // FUTURE: To get the i-th file pointer, set this to `i * unsafe.Sizeof(uintptr(0))`
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(filesPtr) + offset))
		file := &glib.File{
			Object: &glib.Object{
				GObject: glib.ToGObject(*ptr),
			},
		}
		path := file.GetPath()
		// `path` is empty when an URL is passed. In such case, try to get the URL directly from `args`
		if path == "" && len(nonFlagArgs) >= 2 && util.IsLikelyHTTPURL(nonFlagArgs[1]) {
			app.loadArchiveFromURL(nonFlagArgs[1], startupParams.Referer)
		} else {
			app.loadArchiveFromPath(path)
		}
	})

	app.S.GTKApplication = application

	return application
}

func (app *App) ensureDirs() {
	configPath, err := getConfigLocation(AppName)
	if err != nil {
		log.Panicf("getting config location: %v", err)
	}
	userDataPath, err := getUserDataLocation(AppName)
	if err != nil {
		log.Panicf("getting user data location: %v", err)
	}

	app.S.ConfigDirPath = configPath
	app.S.UserDataDirPath = userDataPath
	app.S.ReadLaterDirPath = filepath.Join(userDataPath, ReadLaterDir)

	if err := os.MkdirAll(app.S.ConfigDirPath, 0755); err != nil {
		log.Panicf("creting config directory: %v", err)
	}

	if err := os.MkdirAll(app.S.UserDataDirPath, 0755); err != nil {
		log.Panicf("creating user data directory: %v", err)
	}

	if err := os.MkdirAll(app.S.ReadLaterDirPath, 0755); err != nil {
		log.Panicf("creating read later directory: %v", err)
	}
}

func (app *App) syncStateToConfig() {
	app.setFullscreen(app.Config.Fullscreen)
	app.setHideUI(app.Config.HideUI)
	app.setZoomMode(app.Config.ZoomMode)
	app.setDoublePage(app.Config.DoublePage)
	app.setMangaMode(app.Config.MangaMode)
	app.setBackgroundColor(app.Config.BackgroundColor)
}

func (app *App) handleKeyPress(key uint, shift bool, ctrl bool) {
	switch key {
	case gdk.KEY_v:
		if ctrl {
			app.maybeLoadArchiveFromClipboardURL()
		}
	case gdk.KEY_Down:
		if ctrl {
			app.nextArchive()
		} else if shift {
			app.scroll(0, 1)
		} else {
			app.skipForward()
		}
	case gdk.KEY_Up:
		if ctrl {
			app.previousArchive()
		} else if shift {
			app.scroll(0, -1)
		} else {
			app.skipBackward()
		}
	case gdk.KEY_Right:
		if ctrl {
			app.nextScene()
		} else if shift {
			app.scroll(1, 0)
		} else {
			app.nextPage()
		}
	case gdk.KEY_Left:
		if ctrl {
			app.previousScene()
		} else if shift {
			app.scroll(-1, 0)
		} else {
			app.previousPage()
		}
	case gdk.KEY_space:
		if ctrl {
			app.W.MenuItemPreviousPage.Activate()
		} else {
			app.W.MenuItemNextPage.Activate()
		}
	case gdk.KEY_F11:
		app.W.MenuItemFullscreen.Activate()
	case gdk.KEY_KP_Home:
		app.W.MenuItemFirstPage.Activate()
	case gdk.KEY_KP_End:
		app.W.MenuItemLastPage.Activate()
	case gdk.KEY_KP_Page_Up:
		if ctrl {
			app.W.MenuItemPreviousArchive.Activate()
		} else {
			app.W.MenuItemPreviousPage.Activate()
		}
	case gdk.KEY_KP_Next:
		if ctrl {
			app.W.MenuItemNextArchive.Activate()
		} else {
			app.W.MenuItemNextPage.Activate()
		}
	case gdk.KEY_Alt_L:
		if app.Config.HideUI {
			app.S.UITemporarilyRevealed = !app.S.UITemporarilyRevealed
			app.toggleHideUI(!app.S.UITemporarilyRevealed)
		}
	}
}

func (app *App) setStatus(msg string) {
	contextID := app.W.Statusbar.GetContextId("main")
	app.W.Statusbar.Push(contextID, msg)
}

func (app *App) showError(msg string) {
	app.notificationShow(fmt.Sprintf("Error: %s", msg), LongNotification)
}

func (app *App) quit() {
	app.Config.WindowWidth, app.Config.WindowHeight = app.W.MainWindow.GetSize()

	app.saveConfig()
	app.maybeSaveReadingPosition()

	app.S.GTKApplication.Quit()
}
