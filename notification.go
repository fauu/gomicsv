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

import "github.com/gotk3/gotk3/glib"

var (
	notificationCloseAfterS = map[NotificationLength]uint{
		ShortNotification: 2,
		LongNotification:  8,
	}
	notificationCloseSourceHandle *glib.SourceHandle
)

type NotificationLength = int

const (
	ShortNotification NotificationLength = iota
	LongNotification
)

func (app *App) notificationInit() {
	app.W.NotificationCloseButton.Connect("clicked", func() {
		app.W.NotificationRevealer.SetRevealChild(false)
	})
}

func (app *App) notificationShow(text string, length NotificationLength) {
	app.W.NotificationLabel.SetText(text)
	app.W.NotificationRevealer.SetRevealChild(true)

	// Cancel previous close timeout
	if notificationCloseSourceHandle != nil {
		glib.SourceRemove(*notificationCloseSourceHandle)
	}

	// Set a timeout to automatically close the notificaiton
	handle := glib.TimeoutSecondsAdd(notificationCloseAfterS[length], func() {
		app.notificationHide()
	})
	notificationCloseSourceHandle = &handle
}

func (app *App) notificationHide() {
	app.W.NotificationRevealer.SetRevealChild(false)
	notificationCloseSourceHandle = nil
}
