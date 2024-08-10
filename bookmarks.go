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
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gtk"

	"github.com/fauu/gomicsv/util"
)

var bookmarkMenuItems []*gtk.MenuItem

type Bookmark struct {
	Path       string
	Page       uint
	TotalPages *uint
	Added      time.Time
}

func (app *App) addBookmark() {
	defer app.rebuildBookmarksMenu()

	for i := range app.Config.Bookmarks {
		b := &app.Config.Bookmarks[i]
		if b.Path == app.S.ArchivePath {
			b.Page = uint(app.S.ArchivePos + 1)
			b.TotalPages = util.IntPtrToUintPtr(app.S.Archive.Len())
			b.Added = time.Now()
			return
		}
	}

	app.Config.Bookmarks = append(app.Config.Bookmarks, Bookmark{
		Path:       app.S.ArchivePath,
		TotalPages: util.IntPtrToUintPtr(app.S.Archive.Len()),
		Page:       uint(app.S.ArchivePos + 1),
		Added:      time.Now(),
	})
}

func (app *App) rebuildBookmarksMenu() {
	for i := range bookmarkMenuItems {
		app.W.MenuBookmarks.Remove(bookmarkMenuItems[i])
		bookmarkMenuItems[i].Destroy()
	}
	bookmarkMenuItems = nil
	util.GC()

	for i := range app.Config.Bookmarks {
		bookmark := &app.Config.Bookmarks[i]
		base := filepath.Base(bookmark.Path)
		totalPages := "?"
		if bookmark.TotalPages != nil {
			totalPages = fmt.Sprint(*bookmark.TotalPages)
		}
		label := fmt.Sprintf("%s (%d/%s)", base, bookmark.Page, totalPages)
		bookmarkMenuItem, err := gtk.MenuItemNewWithLabel(label)
		if err != nil {
			app.showError(err.Error())
			return
		}
		bookmarkMenuItem.Connect("activate", func() {
			if app.S.ArchivePath != bookmark.Path {
				app.loadArchiveFromPath(bookmark.Path)
			}
			app.setPage(int(bookmark.Page) - 1)
		})
		bookmarkMenuItems = append(bookmarkMenuItems, bookmarkMenuItem)
		app.W.MenuBookmarks.Append(bookmarkMenuItem)
	}
	app.W.MenuBookmarks.ShowAll()
}
