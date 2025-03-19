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
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/fauu/gomicsv/archive"
	"github.com/fauu/gomicsv/imgdiff"
	"github.com/fauu/gomicsv/pagecache"
	"github.com/fauu/gomicsv/util"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func (app *App) loadArchiveFromURL(url string, httpReferer string) {
	app.doLoadArchive(url, true, httpReferer)
}

func (app *App) loadArchiveFromPath(path string) {
	app.doLoadArchive(path, false, "")
}

func (app *App) doLoadArchive(path string, assumeHTTPURL bool, httpReferer string) {
	if strings.TrimSpace(path) == "" {
		return
	}

	if assumeHTTPURL && !util.IsLikelyHTTPURL(path) {
		// For cases when a non-fully qualified URL is provided
		path = "https://" + path
	}

	if !assumeHTTPURL && !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting current working directory: %v", err)
			return
		}
		path = filepath.Join(wd, path)
	}

	if app.archiveIsLoaded() {
		app.archiveClose()
	}

	app.S.ImageHashes = make(map[int]imgdiff.Hash)

	app.S.ArchivePath = path

	cache := pagecache.NewPageCache()
	app.S.PageCache = &cache
	interval := uint(pageCacheTrimInterval.Seconds())
	handle := glib.TimeoutSecondsAdd(interval, app.trimPageCache)
	app.S.PageCacheTrimTimeoutHandle = &handle

	var err error
	if app.S.Archive, err = archive.NewArchive(path, &cache, httpReferer); err != nil {
		app.showError(fmt.Sprintf("Couldn't open %s: %v", path, err))
		return
	}

	app.archiveHandleLenKnowledge(app.S.Archive.Len() != nil)

	app.W.ButtonRightArchive.SetSensitive(!assumeHTTPURL)
	app.W.ButtonLeftArchive.SetSensitive(!assumeHTTPURL)

	app.W.MenuItemCopyImageToClipboard.SetSensitive(true)

	if !assumeHTTPURL {
		err := os.Chdir(filepath.Dir(app.S.ArchivePath))
		if err != nil {
			log.Printf("Could not chdir into archive path: %v", err)
			return
		}

		if app.Config.RememberRecent {
			u := &url.URL{Path: path, Scheme: "file"}
			ok := app.S.RecentManager.AddItem(u.String())
			if !ok {
				log.Printf("Could not add %s as a recent item", path)
			}
		}
	}

	startPage := 0
	if (!assumeHTTPURL && app.Config.RememberPosition) || (assumeHTTPURL && app.Config.RememberPositionHTTP) {
		savedArchivePos, err := app.loadReadingPosition(path)
		if err == nil {
			startPage = savedArchivePos
		} else if !os.IsNotExist(err) {
			log.Printf("Error loading reading position: %v", err)
		}
	}
	app.doSetPage(startPage)
}

func (app *App) archiveIsLoaded() bool {
	return app.S.ArchivePath != ""
}

func (app *App) archiveGetBaseName() string {
	var name string
	if app.S.Archive.Kind() == archive.HTTPKind {
		name = app.S.ArchivePath
	} else {
		name = filepath.Base(app.S.ArchivePath)
		if ext := filepath.Ext(name); len(ext) > 1 {
			name = strings.TrimSuffix(name, ext)
		}
	}
	return name
}

func (app *App) archiveClose() {
	if !app.archiveIsLoaded() {
		return
	}

	app.maybeSaveReadingPosition()

	app.S.Archive.Close()

	app.S.ArchivePath = ""
	app.S.ArchivePos = 0

	app.S.PageCache = nil

	app.S.ImageHashes = nil

	app.clearJumpmarks()

	// Cancel page cache trim timeout
	if app.S.PageCacheTrimTimeoutHandle != nil {
		glib.SourceRemove(*app.S.PageCacheTrimTimeoutHandle)
		app.S.PageCacheTrimTimeoutHandle = nil
	}

	app.W.ImageL.Clear()
	app.W.ImageR.Clear()
	app.S.PixbufL = nil
	app.S.PixbufR = nil
	app.S.Cursor.reset()
	app.W.MenuItemCopyImageToClipboard.SetSensitive(false)
	app.setStatus("")
	app.W.MainWindow.SetTitle(AppNameDisplay)

	util.GC()
}

func (app *App) archiveHandleLenKnowledge(known bool) {
	var lastPageButton *gtk.ToolButton
	if app.isNavigationRightToLeft() {
		lastPageButton = app.W.ButtonLeftmostPage
	} else {
		lastPageButton = app.W.ButtonRightmostPage
	}
	lastPageButton.SetSensitive(known)
	app.W.MenuItemLastPage.SetSensitive(known)
	app.W.MenuItemRandom.SetSensitive(known)
}

func (app *App) maybeLoadArchiveFromClipboardURL() {
	clipboard, err := gtk.ClipboardGet(gdk.GdkAtomIntern("CLIPBOARD", true))
	if err != nil {
		log.Panicf("getting clipboard: %v", err)
	} else {
		text, err := clipboard.WaitForText()
		if err != nil {
			log.Panicf("getting clipboard text: %v", err)
		} else {
			if util.IsLikelyHTTPURL(text) {
				app.loadArchiveFromURL(text, "")
			}
		}
	}
}
