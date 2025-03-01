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
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/fauu/gomicsv/archive"
	"github.com/fauu/gomicsv/imgdiff"
)

func (app *App) randomPage() {
	if !app.archiveIsLoaded() {
		return
	}
	if app.S.Archive.Len() == nil {
		return
	}

	app.setPage(rand.Int() % *app.S.Archive.Len())
}

func (app *App) pageLeft() {
	if app.Config.MangaMode {
		app.nextPage()
	} else {
		app.previousPage()
	}
}

func (app *App) pageRight() {
	if app.Config.MangaMode {
		app.previousPage()
	} else {
		app.nextPage()
	}
}

func (app *App) skipLeft() {
	if app.Config.MangaMode {
		app.skipForward()
	} else {
		app.skipBackward()
	}
}

func (app *App) skipRight() {
	if app.Config.MangaMode {
		app.skipBackward()
	} else {
		app.skipForward()
	}
}

func (app *App) leftmostPage() {
	if app.Config.MangaMode {
		app.lastPage()
	} else {
		app.firstPage()
	}
}

func (app *App) rightmostPage() {
	if app.Config.MangaMode {
		app.firstPage()
	} else {
		app.lastPage()
	}
}

func (app *App) archiveLeft() {
	if app.Config.MangaMode {
		app.nextArchive()
	} else {
		app.previousArchive()
	}
}

func (app *App) archiveRight() {
	if app.Config.MangaMode {
		app.previousArchive()
	} else {
		app.nextArchive()
	}
}

func (app *App) previousPage() {
	if !app.archiveIsLoaded() {
		if app.Config.Seamless {
			app.previousArchive()
		}
		return
	}

	if app.Config.Random && app.S.Archive.Len() != nil {
		app.randomPage()
		return
	}

	n := 1
	if app.Config.DoublePage && app.S.ArchivePos > 1 {
		n = 2
	}

	if app.Config.Seamless && app.S.ArchivePos+1 <= n {
		app.previousArchive()
		return
	}

	app.setPage(app.S.ArchivePos - n)

	if app.Config.DoublePage &&
		app.shouldForceSinglePage() &&
		app.S.Archive.Len() != nil &&
		*app.S.Archive.Len()-app.S.ArchivePos > 1 {

		app.nextPage()
	}
}

func (app *App) nextPage() {
	if !app.archiveIsLoaded() {
		if app.Config.Seamless {
			app.nextArchive()
		}
		return
	}

	if app.Config.Random && app.S.Archive.Len() != nil {
		app.randomPage()
		return
	}

	n := 1
	if app.Config.DoublePage &&
		!app.shouldForceSinglePage() &&
		app.S.Archive.Len() != nil &&
		*app.S.Archive.Len() > app.S.ArchivePos+2 {

		n = 2
	}

	if app.Config.Seamless &&
		app.S.Archive.Len() != nil &&
		*app.S.Archive.Len()-app.S.ArchivePos <= n {

		app.nextArchive()
		return
	}

	app.setPage(app.S.ArchivePos + n)
}

func (app *App) firstPage() {
	if !app.archiveIsLoaded() {
		return
	}
	app.setPage(0)
}

func (app *App) lastPage() {
	if !app.archiveIsLoaded() {
		return
	}

	if app.S.Archive.Len() == nil {
		return
	}

	offset := -1
	if app.Config.DoublePage && *app.S.Archive.Len() >= 2 {
		offset = -2
	}
	app.setPage(*app.S.Archive.Len() + offset)
}

func (app *App) imageHash(n int) (imgdiff.Hash, bool) {
	if hash, ok := app.S.ImageHashes[n]; ok {
		return hash, true
	}

	pixbuf, err := app.S.Archive.Load(n, app.Config.EmbeddedOrientation, 0)
	if err != nil {
		app.showError(err.Error())
		return 0, false
	}

	return imgdiff.DHash(pixbuf), true
}

func (app *App) skipForward() {
	app.setPage(app.S.ArchivePos + app.Config.NSkip)
}

func (app *App) skipBackward() {
	app.setPage(app.S.ArchivePos - app.Config.NSkip)
}

func (app *App) sceneLeft() {
	if app.Config.MangaMode {
		app.nextScene()
	} else {
		app.previousScene()
	}
}

func (app *App) sceneRight() {
	if app.Config.MangaMode {
		app.previousScene()
	} else {
		app.nextScene()
	}
}

func (app *App) nextScene() {
	if !app.archiveIsLoaded() {
		return
	}

	if app.S.Archive.Len() == nil {
		return
	}

	if app.S.PixbufL == nil {
		return
	}
	hash := imgdiff.DHash(app.S.PixbufL)

	dn := app.Config.SceneScanSkip
	if *app.S.Archive.Len()-1-app.S.ArchivePos <= dn {
		dn = 1
	}

	for n := app.S.ArchivePos + 1; n < *app.S.Archive.Len(); n += dn {
		h, ok := app.imageHash(n)
		if !ok {
			return
		}
		distance := float32(imgdiff.Distance(hash, h)) / 64

		if distance > app.Config.ImageDiffThres {
			if dn == 1 || n == app.S.ArchivePos+1 {
				app.doSetPage(n)
				return
			}

			// Did we go too fast?
			for l := n - 1; l >= app.S.ArchivePos+1; l-- {
				h, ok := app.imageHash(l)
				if !ok {
					return
				}
				d := float32(imgdiff.Distance(hash, h)) / 64
				if d <= app.Config.ImageDiffThres {
					app.doSetPage(l + 1)
					return
				}
			}
			return
		}
	}
}

func (app *App) previousScene() {
	if !app.archiveIsLoaded() {
		return
	}

	if app.S.PixbufL == nil {
		return
	}
	hash := imgdiff.DHash(app.S.PixbufL)

	dn := app.Config.SceneScanSkip
	if app.S.ArchivePos <= dn {
		dn = 1
	}

	for n := app.S.ArchivePos - 1; n >= 0; n -= dn {
		h, ok := app.imageHash(n)
		if !ok {
			return
		}
		distance := float32(imgdiff.Distance(hash, h)) / 64

		if distance > app.Config.ImageDiffThres {
			if dn == 1 || n == app.S.ArchivePos-1 {
				app.doSetPage(n)
				return
			}

			// Did we go too fast?
			for l := n + 1; l <= app.S.ArchivePos-1; l++ {
				h, ok := app.imageHash(l)
				if !ok {
					return
				}
				d := float32(imgdiff.Distance(hash, h)) / 64
				if d <= app.Config.ImageDiffThres {
					app.doSetPage(l - 1)
					return
				}
			}
			return
		}
	}
}

// TODO(fau): Distinguish a failiure from the "no next archive" condition and inform the user accordingly
func (app *App) nextArchive() bool {
	newName, err := app.archiveNameRelativeToCurrent(1)
	if err != nil {
		log.Printf("Error getting next archive: %v", err)
		return false
	}

	app.loadArchiveFromPath(newName)
	return true
}

// TODO(fau): Distinguish a failiure from the "no previous archive" condition and inform the user accordingly
func (app *App) previousArchive() bool {
	newName, err := app.archiveNameRelativeToCurrent(-1)
	if err != nil {
		log.Printf("Error getting previous archive: %v", err)
		return false
	}

	app.loadArchiveFromPath(newName)
	app.lastPage()
	return true
}

// currentArchiveIdx determines the index of the current archive in the directory. We need to do
// this every time, since the filesystem is mutable
func (app *App) currentArchiveIdx() (idx int, err error) {
	dir, name := filepath.Split(app.S.ArchivePath)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return
		}
	}
	arNames, err := archive.ListInDirectory(dir)
	if err != nil {
		return
	}

	idx = -1
	for i := 0; i < len(arNames); i++ {
		if arNames[i] == name {
			idx = i
		}
	}
	if idx == -1 {
		return 0, errors.New("could not find the current archive in the current directory. Deleted, perhaps?")
	}
	return
}

// archiveNameRelativeToCurrent gets the name of the archive in the current directory whose
// relative position with regards to the current archive is equal to relIdx
// TODO(utkan): Use inotify to avoid obtaining list from the scratch all the time
func (app *App) archiveNameRelativeToCurrent(relIdx int) (newName string, err error) {
	dir, _ := filepath.Split(app.S.ArchivePath)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return
		}
	}
	arNames, err := archive.ListInDirectory(dir)
	if err != nil {
		return
	}

	currIdx, err := app.currentArchiveIdx()
	if err != nil {
		return "", nil
	}

	idx := currIdx + relIdx
	if idx < 0 || idx >= len(arNames) {
		err = errors.New("no more archives in the directory")
		return
	}

	newName = filepath.Join(dir, arNames[idx])
	return
}
