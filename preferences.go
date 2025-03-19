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
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func (app *App) preferencesInit() {
	_, err := app.W.PreferencesDialog.AddButton("_OK", gtk.RESPONSE_ACCEPT)
	checkDialogAddButtonErr(err)

	app.W.BackgroundColorButton.Connect("color-set", func(self *glib.Object) {
		chooser := &gtk.ColorChooser{
			Object: self,
		}
		app.setBackgroundColor(NewColorFromGdkRGBA(chooser.GetRGBA()))
	})

	app.W.PagesToSkipSpinButton.SetRange(1, 100)
	app.W.PagesToSkipSpinButton.SetIncrements(1, 10)
	app.W.PagesToSkipSpinButton.SetValue(float64(app.Config.NSkip))
	app.W.PagesToSkipSpinButton.Connect("value-changed", func(self *gtk.SpinButton) {
		app.Config.NSkip = int(self.GetValue())
	})

	app.W.InterpolationComboBoxText.Connect("changed", func(self *gtk.ComboBoxText) {
		app.setInterpolation(self.GetActive())
	})

	app.W.SmartScrollCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setSmartScroll(self.GetActive())
	})

	app.W.MangaModeReverseNavigationCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setMangaModeReverseNavigation(self.GetActive())
	})

	app.W.RememberRecentCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setRememberRecent(self.GetActive())
	})

	app.W.RememberPositionCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setRememberPosition(self.GetActive())
		app.W.RememberPositionHTTPCheckButton.SetSensitive(self.GetActive())
		if !self.GetActive() {
			app.W.RememberPositionHTTPCheckButton.SetActive(false)
		}
	})

	app.W.RememberPositionHTTPCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setRememberPositionHTTP(self.GetActive())
	})

	app.W.OneWideCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setOneWide(self.GetActive())
	})

	app.W.EmbeddedOrientationCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setEmbeddedOrientation(self.GetActive())
	})

	app.W.HideIdleCursorCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setHideIdleCursor(self.GetActive())
	})

	app.W.KamiteEnabledCheckButton.Connect("toggled", func(self *gtk.CheckButton) {
		app.setKamiteEnabled(self.GetActive())
		app.W.KamitePortContainer.SetSensitive(self.GetActive())
	})

	app.W.KamitePortEntry.Connect("insert-text", func(self *gtk.Entry, text string) {
		// Allow only digits
		for _, c := range text {
			if c < '0' || c > '9' {
				self.StopEmission("insert-text")
				return
			}
		}
	})
	app.W.KamitePortEntry.Connect("changed", func(self *gtk.Entry) {
		text, err := self.GetText()
		if err != nil {
			log.Panicf("getting Kamite port entry text: %v", err)
		}
		if text == "" {
			// POLISH: Restore default value
			return
		}
		port, err := strconv.Atoi(text)
		if err != nil {
			log.Panicf("parsing Kamite port entry as int: %v", err)
		}
		app.setKamitePort(port)
	})
}

func (app *App) preferencesDialogRun() {
	app.S.Cursor.ForceVisible = true
	res := gtk.ResponseType(app.W.PreferencesDialog.Run())
	app.W.PreferencesDialog.Hide()
	if res == gtk.RESPONSE_ACCEPT {
		app.saveConfig()
	}
	app.S.Cursor.ForceVisible = false
}
