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

import (
	"reflect"

	"github.com/gotk3/gotk3/gtk"
)

type Widgets struct {
	MainWindow                         *gtk.ApplicationWindow `build:"MainWindow"`
	MainContainer                      *gtk.Box               `build:"MainContainer"`
	Menubar                            *gtk.MenuBar           `build:"Menubar"`
	ScrolledWindow                     *gtk.ScrolledWindow    `build:"ScrolledWindow"`
	ImageViewport                      *gtk.Viewport          `build:"ImageViewport"`
	ImageBox                           *gtk.Box               `build:"ImageBox"`
	ImageL                             *gtk.Image             `build:"ImageL"`
	ImageR                             *gtk.Image             `build:"ImageR"`
	NotificationRevealer               *gtk.Revealer          `build:"NotificationRevealer"`
	NotificationLabel                  *gtk.Label             `build:"NotificationLabel"`
	NotificationCloseButton            *gtk.Button            `build:"NotificationCloseButton"`
	MenuAbout                          *gtk.Menu              `build:"MenuAbout"`
	AboutDialog                        *gtk.AboutDialog       `build:"AboutDialog"`
	MenuFile                           *gtk.Menu              `build:"MenuFile"`
	MenuEdit                           *gtk.Menu              `build:"MenuEdit"`
	MenuView                           *gtk.Menu              `build:"MenuView"`
	MenuNavigation                     *gtk.Menu              `build:"MenuNavigation"`
	MenuBookmarks                      *gtk.Menu              `build:"MenuBookmarks"`
	MenuJumpmarks                      *gtk.Menu              `build:"MenuJumpmarks"`
	Statusbar                          *gtk.Statusbar         `build:"Statusbar"`
	MenuItemOpen                       *gtk.MenuItem          `build:"MenuItemOpen"`
	MenuItemOpenURL                    *gtk.MenuItem          `build:"MenuItemOpenURL"`
	MenuItemClose                      *gtk.MenuItem          `build:"MenuItemClose"`
	MenuItemQuit                       *gtk.MenuItem          `build:"MenuItemQuit"`
	MenuItemSaveImage                  *gtk.MenuItem          `build:"MenuItemSaveImage"`
	ArchiveFileChooserDialog           *gtk.FileChooserDialog `build:"ArchiveFileChooserDialog"`
	SaveImageFileChooserDialog         *gtk.FileChooserDialog `build:"SaveImageFileChooserDialog"`
	OpenURLDialog                      *gtk.Dialog            `build:"OpenURLDialog"`
	OpenURLDialogURLEntry              *gtk.Entry             `build:"OpenURLDialogURLEntry"`
	OpenURLDialogExplanationLabel      *gtk.Label             `build:"OpenURLDialogExplanationLabel"`
	OpenURLDialogRefererEntry          *gtk.Entry             `build:"OpenURLDialogRefererEntry"`
	Toolbar                            *gtk.Toolbar           `build:"Toolbar"`
	ButtonNextPage                     *gtk.ToolButton        `build:"ButtonNextPage"`
	ButtonPreviousPage                 *gtk.ToolButton        `build:"ButtonPreviousPage"`
	ButtonLastPage                     *gtk.ToolButton        `build:"ButtonLastPage"`
	ButtonFirstPage                    *gtk.ToolButton        `build:"ButtonFirstPage"`
	ButtonNextArchive                  *gtk.ToolButton        `build:"ButtonNextArchive"`
	ButtonPreviousArchive              *gtk.ToolButton        `build:"ButtonPreviousArchive"`
	ButtonSkipForward                  *gtk.ToolButton        `build:"ButtonSkipForward"`
	ButtonSkipBackward                 *gtk.ToolButton        `build:"ButtonSkipBackward"`
	MenuItemNextPage                   *gtk.MenuItem          `build:"MenuItemNextPage"`
	MenuItemPreviousPage               *gtk.MenuItem          `build:"MenuItemPreviousPage"`
	MenuItemLastPage                   *gtk.MenuItem          `build:"MenuItemLastPage"`
	MenuItemFirstPage                  *gtk.MenuItem          `build:"MenuItemFirstPage"`
	MenuItemNextArchive                *gtk.MenuItem          `build:"MenuItemNextArchive"`
	MenuItemPreviousArchive            *gtk.MenuItem          `build:"MenuItemPreviousArchive"`
	MenuItemSkipForward                *gtk.MenuItem          `build:"MenuItemSkipForward"`
	MenuItemSkipBackward               *gtk.MenuItem          `build:"MenuItemSkipBackward"`
	MenuItemEnlarge                    *gtk.CheckMenuItem     `build:"MenuItemEnlarge"`
	MenuItemShrink                     *gtk.CheckMenuItem     `build:"MenuItemShrink"`
	MenuItemFullscreen                 *gtk.CheckMenuItem     `build:"MenuItemFullscreen"`
	MenuItemHideUI                     *gtk.CheckMenuItem     `build:"MenuItemHideUI"`
	MenuItemSeamless                   *gtk.CheckMenuItem     `build:"MenuItemSeamless"`
	MenuItemRandom                     *gtk.CheckMenuItem     `build:"MenuItemRandom"`
	MenuItemCopyImageToClipboard       *gtk.MenuItem          `build:"MenuItemCopyImageToClipboard"`
	MenuItemPreferences                *gtk.MenuItem          `build:"MenuItemPreferences"`
	MenuItemHFlip                      *gtk.CheckMenuItem     `build:"MenuItemHFlip"`
	MenuItemVFlip                      *gtk.CheckMenuItem     `build:"MenuItemVFlip"`
	MenuItemMangaMode                  *gtk.CheckMenuItem     `build:"MenuItemMangaMode"`
	MenuItemDoublePage                 *gtk.CheckMenuItem     `build:"MenuItemDoublePage"`
	MenuItemGoTo                       *gtk.MenuItem          `build:"MenuItemGoTo"`
	MenuItemBestFit                    *gtk.RadioMenuItem     `build:"MenuItemBestFit"`
	MenuItemOriginal                   *gtk.RadioMenuItem     `build:"MenuItemOriginal"`
	MenuItemFitToWidth                 *gtk.RadioMenuItem     `build:"MenuItemFitToWidth"`
	MenuItemFitToHeight                *gtk.RadioMenuItem     `build:"MenuItemFitToHeight"`
	MenuItemAbout                      *gtk.MenuItem          `build:"MenuItemAbout"`
	GoToThumbnailImage                 *gtk.Image             `build:"GoToThumbnailImage"`
	GoToDialog                         *gtk.Dialog            `build:"GoToDialog"`
	GoToSpinButton                     *gtk.SpinButton        `build:"GoToSpinButton"`
	GoToScrollbar                      *gtk.Scrollbar         `build:"GoToScrollbar"`
	PreferencesDialog                  *gtk.Dialog            `build:"PreferencesDialog"`
	BackgroundColorButton              *gtk.ColorButton       `build:"BackgroundColorButton"`
	PagesToSkipSpinButton              *gtk.SpinButton        `build:"PagesToSkipSpinButton"`
	InterpolationComboBoxText          *gtk.ComboBoxText      `build:"InterpolationComboBoxText"`
	RememberRecentCheckButton          *gtk.CheckButton       `build:"RememberRecentCheckButton"`
	RememberPositionCheckButton        *gtk.CheckButton       `build:"RememberPositionCheckButton"`
	RememberPositionHTTPCheckButton    *gtk.CheckButton       `build:"RememberPositionHTTPCheckButton"`
	OneWideCheckButton                 *gtk.CheckButton       `build:"OneWideCheckButton"`
	SmartScrollCheckButton             *gtk.CheckButton       `build:"SmartScrollCheckButton"`
	EmbeddedOrientationCheckButton     *gtk.CheckButton       `build:"EmbeddedOrientationCheckButton"`
	HideIdleCursorCheckButton          *gtk.CheckButton       `build:"HideIdleCursorCheckButton"`
	KamiteEnabledCheckButton           *gtk.CheckButton       `build:"KamiteEnabledCheckButton"`
	KamitePortContainer                *gtk.Box               `build:"KamitePortContainer"`
	KamitePortEntry                    *gtk.Entry             `build:"KamitePortEntry"`
	MenuItemAddBookmark                *gtk.MenuItem          `build:"AddBookmarkMenuItem"`
	MenuItemToggleJumpmark             *gtk.MenuItem          `build:"ToggleJumpmarkMenuItem"`
	MenuItemCycleJumpmarksBackward     *gtk.MenuItem          `build:"CycleJumpmarksBackwardMenuItem"`
	MenuItemCycleJumpmarksForward      *gtk.MenuItem          `build:"CycleJumpmarksForwardMenuItem"`
	MenuItemJumpmarksReturnFromCycling *gtk.MenuItem          `build:"JumpmarksReturnFromCyclingMenuItem"`
	RecentChooserMenu                  *gtk.RecentChooserMenu `build:"RecentChooserMenu"`
}

// loadWidgets fills the Widgets struct based on the glade UI definition file
func (app *App) loadWidgets() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	builder, err := gtk.BuilderNew()
	if err != nil {
		return err
	}

	if err = builder.AddFromString(uiDef); err != nil {
		return err
	}

	widgets := &Widgets{}

	widgetsStruct := reflect.ValueOf(widgets).Elem()

	for i := 0; i < widgetsStruct.NumField(); i++ {
		field := widgetsStruct.Field(i)
		widget := widgetsStruct.Type().Field(i).Tag.Get("build")
		if widget == "" {
			continue
		}

		obj, err := builder.GetObject(widget)
		if err != nil {
			return err
		}

		w := reflect.ValueOf(obj).Convert(field.Type())
		field.Set(w)
	}

	app.W = *widgets

	return nil
}
