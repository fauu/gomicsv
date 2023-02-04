/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2023 Piotr Grabowski
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
	"sort"
	"strings"
	"time"

	"github.com/fauu/gomicsv/pagecache"
	"github.com/fauu/gomicsv/util"
	"github.com/gotk3/gotk3/glib"
)

const (
	pageCacheTrimInterval       = 3 * time.Minute
	preloadedPageKeepAtLeastFor = 7 * time.Minute
)

func (app *App) setPage(n int) {
	if !app.archiveIsLoaded() {
		return
	}

	if n < 0 {
		n = 0
	}

	if app.S.Archive.Len() != nil && n >= *app.S.Archive.Len() {
		n = *app.S.Archive.Len() - 1
	}

	if n == app.S.ArchivePos {
		return
	}

	isPrev := false
	if n == app.S.ArchivePos-1 || (app.Config.DoublePage && n == app.S.ArchivePos-2) {
		isPrev = true
	}

	app.doSetPage(n)

	var scrollFunc func()
	if isPrev {
		scrollFunc = app.scrollToEnd
	} else {
		scrollFunc = app.scrollToStart
	}
	glib.TimeoutAdd(0, scrollFunc)
}

func (app *App) doSetPage(n int) {
	if !app.archiveIsLoaded() {
		return
	}

	app.jumpmarksHandleSetPage(n - 1)

	var err error
	app.S.PixbufL, err = app.S.Archive.Load(n, app.Config.EmbeddedOrientation, app.Config.NPreload)
	if err != nil {
		return
	}

	app.S.ArchivePos = n

	app.S.PixbufR = nil
	if app.Config.DoublePage && (app.S.Archive.Len() == nil || *app.S.Archive.Len() > n+1) {
		app.S.PixbufR, err = app.S.Archive.Load(n+1, app.Config.EmbeddedOrientation, app.Config.NPreload)
		if err != nil {
			app.showError(err.Error())
			return
		}
	}

	util.GC()

	app.blit()
	app.updateStatus()
}

func (app *App) trimPageCache() {
	if app.S.PageCache == nil {
		return
	}

	// Mark pages as no longer needed for preload unless they're recently fetched or near the current page
	var preloadDontKeep []int
	nearKeepStart := app.S.ArchivePos - app.Config.NPreload
	nearKeepEnd := app.S.ArchivePos + app.Config.NPreload
	for i, entry := range app.S.PageCache.Pages {
		if time.Since(entry.Time) > preloadedPageKeepAtLeastFor && (i < nearKeepStart || i > nearKeepEnd) {
			preloadDontKeep = append(preloadDontKeep, i)
		}
	}
	if len(preloadDontKeep) > 0 {
		app.S.PageCache.DontKeepSlice(preloadDontKeep, pagecache.KeepReasonPreload)
	}

	// Trim unneeded pages
	app.S.PageCache.Trim()

	// Log cache status
	var indices []int
	for i := range app.S.PageCache.Pages {
		indices = append(indices, i+1)
	}
	sort.Ints(indices)
	var strIndices []string
	for _, index := range indices {
		strIndices = append(strIndices, fmt.Sprint(index))
	}

	var remainingStr string
	if len(strIndices) == 0 {
		remainingStr = "-"
	} else {
		remainingStr = strings.Join(strIndices, ", ")
	}
	log.Printf("Trimmed page cache. Remaining pages: %s", remainingStr)
}
