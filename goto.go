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

	"github.com/fauu/gomicsv/util"
	"github.com/gotk3/gotk3/gtk"
)

func (app *App) goToDialogInit() {
	_, err := app.W.GoToDialog.AddButton("_Cancel", gtk.RESPONSE_CANCEL)
	checkDialogAddButtonErr(err)
	goButton, err := app.W.GoToDialog.AddButton("_Go", gtk.RESPONSE_ACCEPT)
	checkDialogAddButtonErr(err)

	app.W.GoToDialog.SetDefault(goButton)

	app.W.GoToSpinButton.Connect("value-changed", func() {
		app.W.GoToScrollbar.SetValue(app.W.GoToSpinButton.GetValue())
		app.goToDialogUpdateThumbnail()
	})

	app.W.GoToScrollbar.Connect("value-changed", func() {
		app.W.GoToSpinButton.SetValue(app.W.GoToScrollbar.GetValue())
	})
}

func (app *App) goToDialogUpdateThumbnail() {
	n := int(app.W.GoToSpinButton.GetValue() - 1)
	pixbuf, err := app.S.Archive.Load(n, app.Config.EmbeddedOrientation, 0)
	if err != nil {
		log.Printf("Error getting thumbnail: %v", err)
		return
	}

	w, h := util.Fit(pixbuf.GetWidth(), pixbuf.GetHeight(), 128, 128)

	scaled, err := pixbuf.ScaleSimple(w, h, interpolations[app.Config.Interpolation])
	if err != nil {
		log.Printf("Error scaling thumbnail: %v", err)
		return
	}

	app.S.GoToThumbPixbuf = scaled
	app.W.GoToThumbnailImage.SetFromPixbuf(scaled)

	util.GC()
}

func (app *App) goToDialogRun() {
	if !app.archiveIsLoaded() {
		return
	}

	if app.S.Archive.Len() == nil {
		app.goToDialogUpdateSpinButton(9999)
		app.W.GoToScrollbar.Hide()
	} else {
		app.goToDialogUpdateSpinButton(*app.S.Archive.Len())
		app.W.GoToScrollbar.SetRange(1, float64(*app.S.Archive.Len()))
		app.W.GoToScrollbar.SetValue(float64(app.S.ArchivePos) + 1)
		app.W.GoToScrollbar.SetIncrements(1, float64(*app.S.Archive.Len()))
		app.W.GoToScrollbar.Show()
	}

	app.W.GoToSpinButton.GrabFocus()

	app.goToDialogUpdateThumbnail()

	res := gtk.ResponseType(app.W.GoToDialog.Run())
	app.W.GoToDialog.Hide()
	if res == gtk.RESPONSE_ACCEPT {
		app.setPage(int(app.W.GoToSpinButton.GetValue()) - 1)

		app.W.GoToThumbnailImage.Clear()
		app.S.GoToThumbPixbuf = nil
		util.GC()
	}
}

func (app *App) goToDialogUpdateSpinButton(archiveLen int) {
	app.W.GoToSpinButton.SetRange(1, float64(archiveLen))
	app.W.GoToSpinButton.SetValue(float64(app.S.ArchivePos) + 1)
	app.W.GoToSpinButton.SetIncrements(1, float64(app.Config.NSkip))
}
