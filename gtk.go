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
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func getDefaultPointerDevice() (*gdk.Device, error) {
	display, err := gdk.DisplayGetDefault()
	if err != nil {
		return nil, err
	}
	seat, err := display.GetDefaultSeat()
	if err != nil {
		return nil, err
	}
	pointerDevice, err := seat.GetPointer()
	if err != nil {
		return nil, err
	}
	return pointerDevice, nil
}

func pixbufCopyPixels(src *gdk.Pixbuf, dst *gdk.Pixbuf, dstStartX int) {
	nChannels := src.GetNChannels() // ASSUMPTION: Number of channels matches
	srcRowstride, dstRowstride := src.GetRowstride(), dst.GetRowstride()
	srcPixels, dstPixels := src.GetPixels(), dst.GetPixels()
	for y := 0; y < src.GetHeight(); y++ {
		for x := 0; x < src.GetWidth(); x++ {
			xOffset := x * nChannels
			srcIdxBase := y*srcRowstride + xOffset
			dstIdxBase := y*dstRowstride + dstStartX*nChannels + xOffset
			for i := 0; i < nChannels; i++ {
				dstPixels[dstIdxBase+i] = srcPixels[srcIdxBase+i]
			}
		}
	}
}

type MenuWithAccels struct {
	Menu  *gtk.Menu
	Path  string
	Items []MenuItemWithAccels
}

type MenuItemWithAccels struct {
	Item  *gtk.MenuItem
	Accel Accel
}

type Accel struct {
	Key  uint
	Mods gdk.ModifierType
}

func setupMenuAccels(mainWindow *gtk.ApplicationWindow, menuSetups []MenuWithAccels) error {
	for _, menuSetup := range menuSetups {
		accelGroup, err := gtk.AccelGroupNew()
		if err != nil {
			return err
		}
		mainWindow.AddAccelGroup(accelGroup)

		menuSetup.Menu.SetAccelPath(menuSetup.Path)
		menuSetup.Menu.SetAccelGroup(accelGroup)

		for _, item := range menuSetup.Items {
			itemPath := fmt.Sprintf("%s/%s", menuSetup.Path, item.Item.GetLabel())
			gtk.AccelMapAddEntry(itemPath, item.Accel.Key, item.Accel.Mods)
			accelGroup.ConnectByPath(itemPath, item.Item.Activate)
		}
	}

	return nil
}

func checkDialogAddButtonErr(err error) {
	if err != nil {
		log.Panicf("adding a button to a dialog: %v", err)
	}
}
