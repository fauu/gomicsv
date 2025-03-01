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

package archive

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"sort"

	"github.com/fauu/gomicsv/pixbuf"
	"github.com/gotk3/gotk3/gdk"
	"github.com/nwaples/rardecode/v2"
)

type Rar struct {
	files  RarMembers // Sorted by name
	reader *rardecode.ReadCloser
	name   string
}

type RarMember struct {
	Header *rardecode.FileHeader
	Offset int // in terms of files
}

type RarMembers []RarMember

func (p RarMembers) Len() int           { return len(p) }
func (p RarMembers) Less(i, j int) bool { return strcmp(p[i].Header.Name, p[j].Header.Name, true) }
func (p RarMembers) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// NewRar reads supported image filenames from a given rar archive and sorts them
func NewRar(name string) (*Rar, error) {
	var err error

	ar := new(Rar)

	ar.name = filepath.Base(name)
	ar.files = make([]RarMember, 0, MaxArchiveEntries)
	ar.reader, err = rardecode.OpenReader(name)
	if err != nil {
		return nil, err
	}

	offsetAcc := -1
	for {
		offsetAcc += 1

		header, err := ar.reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error reading a file inside the rar archive: %v", err)
		}
		if header.IsDir {
			continue
		}
		if !extensionMatches(header.Name, imageExtensions) {
			continue
		}
		ar.files = append(ar.files, RarMember{
			Header: header,
			Offset: offsetAcc,
		})
	}

	if len(ar.files) == 0 {
		return nil, errors.New(ar.name + ": no supported images in the rar file")
	}

	sort.Sort(RarMembers(ar.files))

	ar.Close()

	return ar, nil
}

func (ar *Rar) checkbounds(i int) error {
	if i < 0 || i >= len(ar.files) {
		return ErrBounds
	}
	return nil
}

func (ar *Rar) Load(i int, autorotate bool, _nPreload int) (*gdk.Pixbuf, error) {
	var err error
	if err = ar.checkbounds(i); err != nil {
		return nil, err
	}

	ar.reader, err = rardecode.OpenReader(ar.name)
	if err != nil {
		return nil, err
	}

	targetOffset := ar.files[i].Offset
	offsetAcc := -1
	for {
		offsetAcc += 1

		header, err := ar.reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error reading a file inside the rar archive: %v", err)
		}
		if header.IsDir {
			continue
		}
		if offsetAcc == targetOffset {
			defer ar.Close()
			return pixbuf.Load(ar.reader, autorotate)
		}
	}

	return nil, fmt.Errorf(ar.name + ": could not find a file inside the rar archive")
}

func (ar *Rar) Kind() Kind {
	return Packed
}

func (ar *Rar) ArchiveName() string {
	return ar.name
}

func (ar *Rar) Name(i int) (string, error) {
	if err := ar.checkbounds(i); err != nil {
		return "", err
	}
	return ar.files[i].Header.Name, nil
}

func (ar *Rar) Len() *int {
	l := len(ar.files)
	return &l
}

func (ar *Rar) Close() error {
	return ar.reader.Close()
}
