/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2022 Piotr Grabowski
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

func (app *App) toolbarInit() {
	app.W.ButtonNextPage.Connect("clicked", app.nextPage)
	app.W.ButtonPreviousPage.Connect("clicked", app.previousPage)
	app.W.ButtonFirstPage.Connect("clicked", app.firstPage)
	app.W.ButtonLastPage.Connect("clicked", app.lastPage)
	app.W.ButtonNextArchive.Connect("clicked", app.nextArchive)
	app.W.ButtonPreviousArchive.Connect("clicked", app.previousArchive)
	app.W.ButtonSkipForward.Connect("clicked", app.skipForward)
	app.W.ButtonSkipBackward.Connect("clicked", app.skipBackward)
}
