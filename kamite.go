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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gotk3/gotk3/gdk"
)

const (
	kamiteRecognizeImageSnipWidthPx  = 550
	kamiteRecognizeImageSnipHeightPx = 900

	kamiteCMDEndpointBaseTpl     = "http://localhost:%d/cmd/"
	kamiteOCRImageEndpoint       = "ocr/image"
	kamiteOCRManualBlockEndpoint = "ocr/manual-block"

	bytesPerPixel = 3
)

func (app *App) kamiteRecognizeManualBlock() {
	endpoint := kamiteMakeEndpointURL(app.Config.KamitePort, kamiteOCRManualBlockEndpoint)
	_, err := http.PostForm(endpoint, nil)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
	}
}

func (app *App) kamiteRecognizeImageUnderCursorBlock() {
	// 1. Determine global pointer coordinates
	pointerDevice, err := getDefaultPointerDevice()
	if err != nil {
		log.Panicf("getting the default pointer device: %v", err)
	}
	swx0, swy0, err := app.W.ScrolledWindow.Widget.TranslateCoordinates(app.W.MainWindow, 0, 0)
	if err != nil {
		log.Panicf("translating widget coordinates: %v", err)
	}
	mainWindowWindow, err := app.W.MainWindow.GetWindow()
	if err != nil {
		log.Panicf("getting GdkWindow of MainWindow widget: %v", err)
	}
	_, x, y, _ := mainWindowWindow.GetDevicePosition(pointerDevice)

	// 2. Check if cursor is over the image container
	if (x < swx0 || x > swx0+app.W.ScrolledWindow.GetAllocatedWidth()) ||
		(y < swy0 || y > swy0+app.W.ScrolledWindow.GetAllocatedHeight()) {
		return
	}

	// 3. Determine over which of the images the cursor is and set the source Pixbuf accordingly
	lx0, ly0, err := app.W.ImageL.Widget.TranslateCoordinates(app.W.ScrolledWindow, 0, 0)
	if err != nil {
		log.Panicf("translating widget coordinates: %v", err)
	}
	rx0, ry0, err := app.W.ImageR.Widget.TranslateCoordinates(app.W.ScrolledWindow, 0, 0)
	if err != nil {
		log.Panicf("translating widget coordinates: %v", err)
	}
	scrolledWindowWindow, err := app.W.ScrolledWindow.GetWindow()
	if err != nil {
		log.Panicf("getting GdkWindow of ScrolledWindow widget: %v", err)
	}
	_, x, y, _ = scrolledWindowWindow.GetDevicePosition(pointerDevice)
	xOffset, yOffset := 0, 0
	var srcPixbuf *gdk.Pixbuf
	if (x > lx0 && x < lx0+app.W.ImageL.GetAllocatedWidth()) &&
		(y > ly0 && y < ly0+app.W.ImageL.GetAllocatedHeight()) {
		srcPixbuf = app.W.ImageL.GetPixbuf()
		xOffset = lx0
		yOffset = ly0
	} else if (x > rx0 && x < rx0+app.W.ImageR.GetAllocatedWidth()) &&
		(y > ry0 && y < ry0+app.W.ImageR.GetAllocatedHeight()) {
		srcPixbuf = app.W.ImageR.GetPixbuf()
		xOffset = rx0
		yOffset = ry0
	} else {
		return
	}
	targetX, targetY := x-xOffset, y-yOffset // Relative to the source Pixbuf

	srcW, srcH := srcPixbuf.GetWidth(), srcPixbuf.GetHeight()
	srcNChannels := srcPixbuf.GetNChannels()
	srcPixels := srcPixbuf.GetPixels()
	srcRowstride := srcPixbuf.GetRowstride()

	// 4. Grab area around the cursor
	snipW := kamiteRecognizeImageSnipWidthPx
	snipH := kamiteRecognizeImageSnipHeightPx
	snipSourceX0, snipSourceY0 := targetX-(snipW/2), targetY-(snipH/2)
	snipBytes := make([]byte, snipW*snipH*bytesPerPixel)
	for y := 0; y < snipH; y++ {
		for x := 0; x < snipW; x++ {
			srcX, srcY := snipSourceX0+x, snipSourceY0+y
			var r, g, b byte
			if srcX < 0 || srcY < 0 || srcX >= srcW || srcY >= srcH {
				// Beyond source Pixbuf bounds
				r, g, b = 255, 255, 255
			} else {
				idx := srcY*srcRowstride + srcX*srcNChannels
				r, g, b = srcPixels[idx], srcPixels[idx+1], srcPixels[idx+2]
			}
			idx := ((y * snipW) + x) * bytesPerPixel
			snipBytes[idx] = r
			snipBytes[idx+1] = g
			snipBytes[idx+2] = b
		}
	}

	// 5. Send
	go kamiteSendOCRImageCommand(
		app.Config.KamitePort,
		snipBytes,
		snipW,
		snipH,
	)
}

type KamiteOCRImageCommandParams struct {
	BytesB64 string `json:"bytesB64"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func kamiteSendOCRImageCommand(port int, imgBytes []byte, w, h int) {
	paramsJSON, err := json.Marshal(KamiteOCRImageCommandParams{
		Width:    w,
		Height:   h,
		BytesB64: base64.StdEncoding.EncodeToString(imgBytes),
	})
	if err != nil {
		log.Printf("Error encoding Kamite command params: %v", err)
		return
	}

	_, err = http.Post(
		kamiteMakeEndpointURL(port, kamiteOCRImageEndpoint),
		"application/json",
		bytes.NewReader(paramsJSON),
	)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
	}
}

func kamiteMakeEndpointURL(port int, suffix string) string {
	path, _ := url.JoinPath(fmt.Sprintf(kamiteCMDEndpointBaseTpl, port), suffix)
	return path
}
