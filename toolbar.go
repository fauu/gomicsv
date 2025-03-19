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
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func (app *App) toolbarInit() {
	app.W.ButtonPageLeft.Connect("clicked", app.pageLeft)
	app.W.ButtonPageRight.Connect("clicked", app.pageRight)
	app.W.ButtonLeftmostPage.Connect("clicked", app.leftmostPage)
	app.W.ButtonRightmostPage.Connect("clicked", app.rightmostPage)
	app.W.ButtonLeftArchive.Connect("clicked", app.archiveLeft)
	app.W.ButtonRightArchive.Connect("clicked", app.archiveRight)
	app.W.ButtonSkipLeft.Connect("clicked", app.skipLeft)
	app.W.ButtonSkipRight.Connect("clicked", app.skipRight)
}

func swapToolButtonsText(a, b *gtk.ToolButton) {
	var err error

	tmp := a.GetLabel()
	a.SetLabel(b.GetLabel())
	b.SetLabel(tmp)

	aTooltipText, err := a.GetTooltipText()
	if err != nil {
		logSetTooltipTextError(err)
		return
	}
	bTooltipText, err := b.GetTooltipText()
	if err != nil {
		logSetTooltipTextError(err)
		return
	}

	// NOTE: SetTooltipText() fails silently
	err = a.SetProperty("tooltip-text", bTooltipText)
	if err != nil {
		logSetTooltipTextError(err)
		return
	}
	err = b.SetProperty("tooltip-text", aTooltipText)
	if err != nil {
		logSetTooltipTextError(err)
		return
	}
}

func logSetTooltipTextError(err error) {
	log.Printf("Error setting TooltipButton text: %v", err)
}

func (app *App) reverseMirrorNavigationButtonsText() {
	swapToolButtonsText(app.W.ButtonPageLeft, app.W.ButtonPageRight)
	swapToolButtonsText(app.W.ButtonSkipLeft, app.W.ButtonSkipRight)
	swapToolButtonsText(app.W.ButtonLeftmostPage, app.W.ButtonRightmostPage)
	swapToolButtonsText(app.W.ButtonLeftArchive, app.W.ButtonRightArchive)
	app.S.MirrorNavigationButtonsTextReversed = !app.S.MirrorNavigationButtonsTextReversed
}
