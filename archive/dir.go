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

package archive

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/fauu/gomicsv/pixbuf"
	"github.com/gotk3/gotk3/gdk"
)

type Dir struct {
	filenames filenames
	name      string
	path      string
}

/* Reads filenames from a directory, and sorts them */
func NewDir(path string) (*Dir, error) {
	var err error

	ar := new(Dir)

	ar.name = filepath.Base(path)
	ar.path = path

	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	ar.filenames = make([]string, 0, len(filenames))

	for _, name := range filenames {
		if !extensionMatches(name, imageExtensions) {
			continue
		}
		ar.filenames = append(ar.filenames, name)
	}

	if len(ar.filenames) == 0 {
		return nil, errors.New(ar.name + ": no images in the directory")
	}

	sort.Sort(ar.filenames)

	return ar, nil
}

func (ar *Dir) checkbounds(i int) error {
	if i < 0 || i >= len(ar.filenames) {
		return ErrBounds
	}
	return nil
}

func (ar *Dir) Load(i int, autorotate bool, _nPreload int) (*gdk.Pixbuf, error) {
	if err := ar.checkbounds(i); err != nil {
		return nil, err
	}

	path := filepath.Join(ar.path, ar.filenames[i])
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return pixbuf.Load(f, autorotate)
}

func (ar *Dir) Name(i int) (string, error) {
	if err := ar.checkbounds(i); err != nil {
		return "", err
	}

	return ar.filenames[i], nil
}

func (ar *Dir) Kind() Kind {
	return Unpacked
}

func (ar *Dir) Len() *int {
	l := len(ar.filenames)
	return &l
}

func (ar *Dir) Close() error {
	return nil
}
