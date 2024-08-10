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
	"log"
	"time"

	"github.com/gotk3/gotk3/gdk"
)

type CursorsState struct {
	Current      Cursor
	Cache        CursorCache
	LastMoved    time.Time
	Visible      bool
	ForceVisible bool
}

type Cursor int

const (
	cursorNone Cursor = iota
	cursorDefault
	cursorGrabbing
)

type CursorCache = map[Cursor]*gdk.Cursor

func (cursor *CursorsState) init() {
	cursor.Cache = CursorCache{
		cursorNone:     loadCursor("none"),
		cursorDefault:  loadCursor("default"),
		cursorGrabbing: loadCursor("grabbing"),
	}
	cursor.Current = cursorDefault
}

func (cursor *CursorsState) reset() {
	cursor.LastMoved = time.Now()
	cursor.Visible = false
	cursor.ForceVisible = false
}

func (app *App) setCursor(cursor Cursor) {
	win, err := app.W.ImageViewport.GetWindow()
	if err != nil {
		log.Panicf("getting Image Viewport's Window: %v", err)
	}
	if cursor != cursorNone {
		app.S.Cursor.Current = cursor
	}
	win.SetCursor(app.S.Cursor.Cache[cursor])
}

func (app *App) hideCursor() {
	app.setCursor(cursorNone)
	app.S.Cursor.Visible = false
}

func (app *App) showCursor() {
	app.setCursor(app.S.Cursor.Current)
	app.S.Cursor.Visible = true
}

func (app *App) updateCursorVisibility() bool {
	shouldBeHidden := false
	if !app.S.DragScroll.InProgress && app.Config.HideIdleCursor && !app.S.Cursor.ForceVisible {
		shouldBeHidden = time.Since(app.S.Cursor.LastMoved).Seconds() > 1
	}

	if shouldBeHidden && app.S.Cursor.Visible {
		app.hideCursor()
	} else if !shouldBeHidden && !app.S.Cursor.Visible {
		app.showCursor()
	}

	return true
}

func loadCursor(name string) *gdk.Cursor {
	disp, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Panicf("getting default display: %v", err)
	}
	c, err := gdk.CursorNewFromName(disp, name)
	if err != nil {
		log.Panicf("creating cursor: %v", err)
	}
	return c
}
