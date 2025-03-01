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
	"sort"

	"github.com/gotk3/gotk3/gtk"

	"github.com/fauu/gomicsv/pagecache"
	"github.com/fauu/gomicsv/util"
)

type Jumpmarks struct {
	list  []int
	cycle JumpmarksCycle
}

type JumpmarksCycle struct {
	page       *int
	returnPage *int
	dontClear  bool
}

type JumpmarkCycleDirection int

const (
	cycleDirectionBackward JumpmarkCycleDirection = iota
	cycleDirectionForward
)

func (jumpmarks Jumpmarks) has(page int) bool {
	for _, mark := range jumpmarks.list {
		// ASSUMPTION: The slice is sorted
		if mark < page {
			continue
		} else if mark == page {
			return true
		} else {
			break
		}
	}
	return false
}

func (jumpmarks Jumpmarks) size() int {
	return len(jumpmarks.list)
}

func (jumpmarks *Jumpmarks) toggle(page int) bool {
	for i, mark := range jumpmarks.list {
		// ASSUMPTION: The slice is sorted
		if mark < page {
			continue
		} else if mark == page {
			// Is in the marked list. Remove by swapping with the last element
			jumpmarks.list[i] = jumpmarks.list[len(jumpmarks.list)-1]
			jumpmarks.list = jumpmarks.list[:len(jumpmarks.list)-1]
			sort.Ints(jumpmarks.list)
			return false
		} else {
			break
		}
	}
	// Isn't in the marked list. Add
	jumpmarks.list = append(jumpmarks.list, page)
	sort.Ints(jumpmarks.list)
	return true
}

func (app *App) clearJumpmarks() {
	// The following isn't needed until we use this function elsewhere than during archive closing, which we currently don't
	// for _, mark := range app.S.Jumpmarks.list {
	// 	app.S.PageCache.UnforbidRemoval(mark)
	// }
	app.S.Jumpmarks = Jumpmarks{}
}

var jumpmarkMenuItems []*gtk.MenuItem

func (app *App) currentPageIsJumpmarked() bool {
	return app.S.Jumpmarks.has(app.S.ArchivePos + 1)
}

func (app *App) toggleJumpmark() {
	page := app.S.ArchivePos + 1
	currentMarked := app.S.Jumpmarks.toggle(page)

	var prefix string
	if currentMarked {
		app.S.PageCache.Keep(app.S.ArchivePos, pagecache.KeepReasonJumpmark)
		prefix = "Marked"
	} else {
		app.S.PageCache.DontKeep(app.S.ArchivePos, pagecache.KeepReasonJumpmark)
		prefix = "Unmarked"
	}
	app.notificationShow(fmt.Sprintf("%s page %d", prefix, page), ShortNotification)

	app.updateJumpmarkToggleLabel(currentMarked)
	app.updateJumpmarkCycleMenuItems()
	app.rebuildJumpmarksMenuList()
	app.updateStatus()
}

func (app *App) cycleJumpmarks(direction JumpmarkCycleDirection) {
	if app.S.Jumpmarks.size() == 0 {
		return
	}
	nextPage := -1
	if app.S.Jumpmarks.cycle.page != nil {
		nextPage = *app.S.Jumpmarks.cycle.page
	} else {
		curr := app.S.ArchivePos + 1
		app.S.Jumpmarks.cycle.returnPage = &curr
	}
	if direction == cycleDirectionForward {
		nextPage++
	} else {
		nextPage--
	}
	if nextPage < 0 {
		nextPage = app.S.Jumpmarks.size() - 1
	} else if nextPage >= app.S.Jumpmarks.size() {
		nextPage = 0
	}
	app.S.Jumpmarks.cycle.page = &nextPage
	app.S.Jumpmarks.cycle.dontClear = true
	app.setPage(app.S.Jumpmarks.list[nextPage] - 1)
	app.S.Jumpmarks.cycle.dontClear = false
}

func (app *App) returnFromCyclingJumpmarks() {
	p := app.S.Jumpmarks.cycle.returnPage
	if p != nil {
		app.setPage(*p - 1)
	}
}

func (app *App) jumpmarksHandleSetPage(page int) {
	jumpmarks := &app.S.Jumpmarks
	if !jumpmarks.cycle.dontClear {
		jumpmarks.cycle.page = nil
		jumpmarks.cycle.returnPage = nil
	}
	app.updateJumpmarkToggleLabel(app.currentPageIsJumpmarked())
	app.updateJumpmarkCycleMenuItems()

	s := false
	if jumpmarks.cycle.returnPage != nil && page != *jumpmarks.cycle.returnPage {
		s = true
	}
	app.W.MenuItemJumpmarksReturnFromCycling.SetSensitive(s)
}

func (app *App) updateJumpmarkToggleLabel(currentMarked bool) {
	var w string
	if currentMarked {
		w = "Unmark"
	} else {
		w = "Mark"
	}
	label := fmt.Sprintf("%s current page", w)
	app.W.MenuItemToggleJumpmark.SetLabel(label)
}

func (app *App) updateJumpmarkCycleMenuItems() {
	s := false
	if app.S.Jumpmarks.size() > 0 {
		s = true
	}
	app.W.MenuItemCycleJumpmarksBackward.SetSensitive(s)
	app.W.MenuItemCycleJumpmarksForward.SetSensitive(s)
}

func (app *App) rebuildJumpmarksMenuList() {
	for _, item := range jumpmarkMenuItems {
		app.W.MenuJumpmarks.Remove(item)
		item.Destroy()
	}
	jumpmarkMenuItems = nil
	util.GC()

	for _, mark := range app.S.Jumpmarks.list {
		label := fmt.Sprintf("%d", mark)
		menuItem, err := gtk.MenuItemNewWithLabel(label)
		if err != nil {
			app.showError(err.Error())
			return
		}
		page := mark - 1 // Make a new variable so that the correct value gets passed to the callback
		menuItem.Connect("activate", func() {
			app.setPage(page)
		})
		jumpmarkMenuItems = append(jumpmarkMenuItems, menuItem)
		app.W.MenuJumpmarks.Append(menuItem)
		menuItem.Show()
	}
}
