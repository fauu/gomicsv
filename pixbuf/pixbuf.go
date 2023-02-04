/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2023 Piotr Grabowski
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

package pixbuf

import (
	"bytes"
	"io"
	"log"

	"github.com/gotk3/gotk3/gdk"
)

func Load(r io.Reader, autorotate bool) (*gdk.Pixbuf, error) {
	w, _ := gdk.PixbufLoaderNew()
	defer w.Close()
	_, err := io.Copy(w, r)
	if err != nil {
		return nil, err
	}

	pixbuf, err := w.GetPixbuf()
	if err != nil {
		return nil, err
	}

	if !autorotate {
		return pixbuf, nil
	}

	return pixbuf.ApplyEmbeddedOrientation()
}

func MustLoad(data []byte) *gdk.Pixbuf {
	pixbuf, err := Load(bytes.NewBuffer(data), true)
	if err != nil {
		log.Panicf("loading pixbuf: %v", err)
	}
	return pixbuf
}
