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

package archive

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/fauu/gomicsv/pagecache"
	"github.com/fauu/gomicsv/util"
	"github.com/gotk3/gotk3/gdk"
)

var (
	ErrBounds = errors.New("Image index out of bounds.")
)

type Archive interface {
	Load(i int, autorotate bool, nPreload int) (*gdk.Pixbuf, error)
	Kind() Kind
	ArchiveName() string
	Name(i int) (string, error)
	Len() *int // nil represents unknown length
	Close() error
}

const (
	MaxArchiveEntries = 4096 * 64
)

func NewArchive(path string, pageCache *pagecache.PageCache, httpReferer string) (Archive, error) {
	if util.IsLikelyHTTPURL(path) {
		return NewHTTP(path, pageCache, httpReferer)
	}

	f, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if f.IsDir() {
		return NewDir(path)
	}

	ext := strings.ToLower(filepath.Ext(path))[1:]
	switch ext {
	case "zip", "cbz":
		return NewZip(path)
	case "7z", "rar", "tar", "tgz", "gz", "tbz2", "cb7", "cbr", "cbt", "lha":
		return nil, errors.New("Archive type not supported, please unpack it first")
	}

	return nil, errors.New("Unknown archive type")
}
