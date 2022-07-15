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
	"fmt"
	"log"
	"net/url"
	"runtime"
	"strings"

	"github.com/fauu/gomicsv/pixbuf"
	"github.com/flytam/filenamify"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func (app *App) menuInit() {
	app.menuInitOpenDialog()
	app.menuInitOpenURLDialog()
	app.menuInitSaveImageDialog()

	app.W.MenuItemQuit.Connect("activate", app.quit)
	app.W.MenuItemClose.Connect("activate", app.archiveClose)
	app.W.MenuItemNextPage.Connect("activate", app.nextPage)
	app.W.MenuItemPreviousPage.Connect("activate", app.previousPage)
	app.W.MenuItemFirstPage.Connect("activate", app.firstPage)
	app.W.MenuItemLastPage.Connect("activate", app.lastPage)
	app.W.MenuItemNextArchive.Connect("activate", app.nextArchive)
	app.W.MenuItemPreviousArchive.Connect("activate", app.previousArchive)
	app.W.MenuItemSkipForward.Connect("activate", app.skipForward)
	app.W.MenuItemSkipBackward.Connect("activate", app.skipBackward)

	app.W.MenuItemEnlarge.Connect("toggled", func() {
		app.setEnlarge(app.W.MenuItemEnlarge.GetActive())
	})

	app.W.MenuItemShrink.Connect("toggled", func() {
		app.setShrink(app.W.MenuItemShrink.GetActive())
	})

	app.W.MenuItemFullscreen.Connect("toggled", func() {
		app.setFullscreen(app.W.MenuItemFullscreen.GetActive())
	})

	app.W.MenuItemHideUI.Connect("toggled", func() {
		app.setHideUI(app.W.MenuItemHideUI.GetActive())
	})

	app.W.MenuItemSeamless.Connect("toggled", func() {
		app.setSeamless(app.W.MenuItemSeamless.GetActive())
	})

	app.W.MenuItemRandom.Connect("toggled", func() {
		app.setRandom(app.W.MenuItemRandom.GetActive())
	})

	app.W.MenuItemHFlip.Connect("toggled", func() {
		app.setHFlip(app.W.MenuItemHFlip.GetActive())
	})

	app.W.MenuItemVFlip.Connect("toggled", func() {
		app.setVFlip(app.W.MenuItemVFlip.GetActive())
	})

	app.W.MenuItemMangaMode.Connect("toggled", func() {
		app.setMangaMode(app.W.MenuItemMangaMode.GetActive())
	})

	app.W.MenuItemDoublePage.Connect("toggled", func() {
		app.setDoublePage(app.W.MenuItemDoublePage.GetActive())
	})

	app.W.MenuItemOriginal.Connect("toggled", func() {
		if app.W.MenuItemOriginal.GetActive() {
			app.setZoomMode(Original)
		}
	})

	app.W.MenuItemBestFit.Connect("toggled", func() {
		if app.W.MenuItemBestFit.GetActive() {
			app.setZoomMode(BestFit)
		}
	})

	app.W.MenuItemFitToWidth.Connect("toggled", func() {
		if app.W.MenuItemFitToWidth.GetActive() {
			app.setZoomMode(FitToWidth)
		}
	})

	app.W.MenuItemFitToHalfWidth.Connect("toggled", func() {
		if app.W.MenuItemFitToHalfWidth.GetActive() {
			app.setZoomMode(FitToHalfWidth)
		}
	})

	app.W.MenuItemFitToHeight.Connect("toggled", func() {
		if app.W.MenuItemFitToHeight.GetActive() {
			app.setZoomMode(FitToHeight)
		}
	})

	app.W.MenuItemCopyImageToClipboard.Connect("activate", func() {
		app.copyImageToClipboard()
	})

	app.W.MenuItemAddBookmark.Connect("activate", app.addBookmark)

	app.W.MenuItemToggleJumpmark.Connect("activate", app.toggleJumpmark)

	app.W.MenuItemCycleJumpmarksBackward.Connect("activate", func() {
		app.cycleJumpmarks(cycleDirectionBackward)
	})

	app.W.MenuItemCycleJumpmarksForward.Connect("activate", func() {
		app.cycleJumpmarks(cycleDirectionForward)
	})

	app.W.MenuItemJumpmarksReturnFromCycling.Connect("activate", app.returnFromCyclingJumpmarks)

	app.W.MenuItemPreferences.Connect("activate", app.preferencesDialogRun)

	app.W.MenuItemAbout.Connect("activate", func() {
		app.S.Cursor.ForceVisible = true
		app.W.AboutDialog.Run()
		app.W.AboutDialog.Hide()
		app.S.Cursor.ForceVisible = false
	})

	app.W.AboutDialog.SetLogo(pixbuf.MustLoad(aboutImg))
	if len(app.S.BuildInfo.Version) >= 0 {
		versionStr := fmt.Sprintf("Version: %s (built: %s)\nCompiler version: %s", app.S.BuildInfo.Version, app.S.BuildInfo.Date, runtime.Version())
		app.W.AboutDialog.SetVersion(versionStr)
	}

	app.menuSetupAccels()

	app.rebuildBookmarksMenu()

	app.goToDialogInit()

	app.W.MenuItemGoTo.Connect("activate", app.goToDialogRun)

	app.W.RecentChooserMenu.Connect("item-activated", func() {
		uri := app.W.RecentChooserMenu.GetCurrentUri()
		u, err := url.Parse(uri)
		if err != nil {
			app.showError(err.Error())
			return
		}
		app.loadArchiveFromPath(u.Path)
	})
}

var FILE_CHOOSER_RESPONSE_ACCEPT gtk.ResponseType = 100

func (app *App) menuInitOpenDialog() {
	app.W.MenuItemOpen.Connect("activate", func() {
		res := gtk.ResponseType(app.W.ArchiveFileChooserDialog.Run())
		app.W.ArchiveFileChooserDialog.Hide()
		if res == FILE_CHOOSER_RESPONSE_ACCEPT {
			filename := app.W.ArchiveFileChooserDialog.GetFilename()
			if filename == "" {
				var err error
				filename, err = app.W.ArchiveFileChooserDialog.GetCurrentFolder()
				if err != nil {
					log.Println("Error getting FileChooser CurrentFolder")
				}
			}
			if filename != "" {
				app.loadArchiveFromPath(filename)
			}
		}
	})

	_, err := app.W.ArchiveFileChooserDialog.AddButton("_Open", FILE_CHOOSER_RESPONSE_ACCEPT)
	checkDialogAddButtonErr(err)
	_, err = app.W.ArchiveFileChooserDialog.AddButton("_Cancel", gtk.RESPONSE_CANCEL)
	checkDialogAddButtonErr(err)
}

func (app *App) menuInitOpenURLDialog() {
	app.W.MenuItemOpenURL.Connect("activate", func() {
		res := gtk.ResponseType(app.W.OpenURLDialog.Run())
		app.W.OpenURLDialog.Hide()
		if res == gtk.RESPONSE_ACCEPT {
			url, err := app.W.OpenURLDialogURLEntry.GetText()
			if err != nil {
				log.Panicf("getting Open URL Dialog URL Entry text: %v", err)
			}
			referer, err := app.W.OpenURLDialogRefererEntry.GetText()
			if err != nil {
				log.Panicf("getting Open URL Dialog Referer Entry text: %v", err)
			}
			app.loadArchiveFromURL(url, referer)
		}
	})

	app.W.OpenURLDialogExplanationLabel.SetMarkup(
		"Specify the <i>direct</i> URL of the image of one of the comic pages, for example <tt>http://my-source/comic-1234/4.png</tt>." +
			" The URL must be in such format, that successive pages can be accessed by inserting their respective numbers into the URL." +
			" Otherwise, the particular source is not currently supported.\n" +
			"    If the program fails to" +
			" guess the pattern followed by the page URLs, try again, but this time manually specifying where the page number" +
			" is in the URL, by replacing it with the placeholder <tt>%d</tt>. For example, if the URL for page 7 is" +
			" <tt>http://my-source/comic-1234?page=7&amp;full=true</tt>, specify <tt>http://my-source/comic-1234?page=<b>%d</b>&amp;full=true</tt>" +
			" above. If the number needs to be padded with zeroes, you can specify its width, for example <tt>%03d</tt> for (<tt>001</tt>, <tt>002</tt>, …).\n" +
			"    Note that certain hosts might use various access restriction measures that could make this program unable to access" +
			" the images, even if they’re accessible when viewed directly on the host’s website. Below are extra options that might" +
			" be useful in this connection.",
	)

	_, err := app.W.OpenURLDialog.AddButton("_Cancel", gtk.RESPONSE_CANCEL)
	checkDialogAddButtonErr(err)
	okButton, err := app.W.OpenURLDialog.AddButton("_Open", gtk.RESPONSE_ACCEPT)
	checkDialogAddButtonErr(err)

	app.W.OpenURLDialog.SetDefault(okButton)
}

func (app *App) menuInitSaveImageDialog() {
	app.W.MenuItemSaveImage.Connect("activate", func() {
		baseName, err := filenamify.FilenamifyV2(app.archiveGetBaseName())
		baseName = strings.ReplaceAll(baseName, ".", "!")
		if err != nil {
			log.Panicf("filenamifying archive base name: %v", err)
		}
		filename := fmt.Sprintf("%s-%000d.png", baseName, app.S.ArchivePos+1)
		app.W.SaveImageFileChooserDialog.SetCurrentName(filename)

		res := gtk.ResponseType(app.W.SaveImageFileChooserDialog.Run())
		app.W.SaveImageFileChooserDialog.Hide()
		if res == gtk.RESPONSE_ACCEPT {
			filename := app.W.SaveImageFileChooserDialog.GetFilename()
			if filename != "" {
				app.saveImage(filename)
			}
		}
	})

	_, err := app.W.SaveImageFileChooserDialog.AddButton("_Save", gtk.RESPONSE_ACCEPT)
	checkDialogAddButtonErr(err)
	_, err = app.W.SaveImageFileChooserDialog.AddButton("_Cancel", gtk.RESPONSE_CANCEL)
	checkDialogAddButtonErr(err)
}

func (app *App) menuSetupAccels() {
	// NOTE: This can't be done in the glade file using the <accelerator> tag under the respective
	//       menu items, because then the bindings stop working when the menubar is hidden.
	//       Unfortunately, only primary keybindings can be set here. Auxilliary ones are bound under
	//       the MainWindow key-press-event signal handler.
	accels := []MenuWithAccels{
		{
			Menu: app.W.MenuFile,
			Path: menuMakeAccelPath("File"),
			Items: []MenuItemWithAccels{
				{app.W.MenuItemOpen, Accel{gdk.KEY_O, gdk.CONTROL_MASK}},
				{app.W.MenuItemOpenURL, Accel{gdk.KEY_O, gdk.CONTROL_MASK | gdk.SHIFT_MASK}},
				{app.W.MenuItemClose, Accel{gdk.KEY_W, gdk.CONTROL_MASK}},
				{app.W.MenuItemSaveImage, Accel{gdk.KEY_F9, 0}},
				{app.W.MenuItemQuit, Accel{gdk.KEY_Q, gdk.CONTROL_MASK}},
			},
		},
		{
			Menu: app.W.MenuEdit,
			Path: menuMakeAccelPath("Edit"),
			Items: []MenuItemWithAccels{
				{app.W.MenuItemCopyImageToClipboard, Accel{gdk.KEY_C, gdk.CONTROL_MASK}},
				{app.W.MenuItemPreferences, Accel{gdk.KEY_P, gdk.CONTROL_MASK}},
			},
		},
		{
			Menu: app.W.MenuView,
			Path: menuMakeAccelPath("View"),
			Items: []MenuItemWithAccels{
				{&app.W.MenuItemHideUI.MenuItem, Accel{gdk.KEY_M, gdk.MOD1_MASK}},
				{&app.W.MenuItemShrink.MenuItem, Accel{gdk.KEY_S, 0}},
				{&app.W.MenuItemEnlarge.MenuItem, Accel{gdk.KEY_E, 0}},
				{&app.W.MenuItemBestFit.MenuItem, Accel{gdk.KEY_B, 0}},
				{&app.W.MenuItemOriginal.MenuItem, Accel{gdk.KEY_O, 0}},
				{&app.W.MenuItemFitToWidth.MenuItem, Accel{gdk.KEY_W, 0}},
				{&app.W.MenuItemFitToHalfWidth.MenuItem, Accel{gdk.KEY_W, gdk.MOD1_MASK}},
				{&app.W.MenuItemFitToHeight.MenuItem, Accel{gdk.KEY_H, 0}},
				{&app.W.MenuItemFullscreen.MenuItem, Accel{gdk.KEY_F, 0}},
				{&app.W.MenuItemRandom.MenuItem, Accel{gdk.KEY_R, 0}},
				{&app.W.MenuItemDoublePage.MenuItem, Accel{gdk.KEY_D, 0}},
				{&app.W.MenuItemVFlip.MenuItem, Accel{gdk.KEY_V, 0}},
				{&app.W.MenuItemHFlip.MenuItem, Accel{gdk.KEY_V, gdk.SHIFT_MASK}},
				{&app.W.MenuItemMangaMode.MenuItem, Accel{gdk.KEY_M, gdk.CONTROL_MASK}},
			},
		},
		{
			Menu: app.W.MenuNavigation,
			Path: menuMakeAccelPath("Navigation"),
			Items: []MenuItemWithAccels{
				{app.W.MenuItemPreviousPage, Accel{gdk.KEY_Page_Up, 0}},
				{app.W.MenuItemNextPage, Accel{gdk.KEY_Page_Down, 0}},
				{app.W.MenuItemFirstPage, Accel{gdk.KEY_Home, 0}},
				{app.W.MenuItemLastPage, Accel{gdk.KEY_End, 0}},
				{app.W.MenuItemPreviousArchive, Accel{gdk.KEY_Page_Up, gdk.CONTROL_MASK}},
				{app.W.MenuItemNextArchive, Accel{gdk.KEY_Page_Down, gdk.CONTROL_MASK}},
				{app.W.MenuItemGoTo, Accel{gdk.KEY_G, 0}},
			},
		},
		{
			Menu: app.W.MenuBookmarks,
			Path: menuMakeAccelPath("Bookmarks"),
			Items: []MenuItemWithAccels{
				{app.W.MenuItemAddBookmark, Accel{gdk.KEY_B, gdk.CONTROL_MASK}},
			},
		},
		{
			Menu: app.W.MenuJumpmarks,
			Path: menuMakeAccelPath("Jumpmarks"),
			Items: []MenuItemWithAccels{
				{app.W.MenuItemToggleJumpmark, Accel{gdk.KEY_M, 0}},
				{app.W.MenuItemCycleJumpmarksBackward, Accel{gdk.KEY_bracketleft, 0}},
				{app.W.MenuItemCycleJumpmarksForward, Accel{gdk.KEY_bracketright, 0}},
				{app.W.MenuItemJumpmarksReturnFromCycling, Accel{gdk.KEY_BackSpace, 0}},
			},
		},
		{
			Menu: app.W.MenuAbout,
			Path: menuMakeAccelPath("About"),
			Items: []MenuItemWithAccels{
				{app.W.MenuItemAbout, Accel{gdk.KEY_F1, 0}},
			},
		},
	}

	err := setupMenuAccels(app.W.MainWindow, accels)
	if err != nil {
		log.Panicf("setting up menu accels: %v", err)
	}
}

func menuMakeAccelPath(menuCategory string) string {
	return fmt.Sprintf("<gomicsv>/%s", menuCategory)
}
