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
	"fmt"
	"log"
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
	subPath, err := findFirstDirContainingSupportedImage(path)
	if err != nil {
		return nil, fmt.Errorf("searching for supported images in provided path: %v", err)
	}
	if subPath == nil {
		return nil, errors.New("Could not find a directory containing supported images")
	}

	dir, err := os.Open(*subPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	ar := Dir{
		filenames: nil,
		name:      filepath.Base(*subPath),
		path:      *subPath,
	}

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

	return &ar, nil
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

func (ar *Dir) ArchiveName() string {
	return ar.name
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

func findFirstDirContainingSupportedImage(path string) (*string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			if subdirs == nil {
				subdirs = make([]string, 0, len(entries))
			}
			subdirs = append(subdirs, entry.Name())
		} else if extensionMatches(entry.Name(), imageExtensions) {
			// Found supported image directly in path
			return &path, nil
		}
	}

	if len(subdirs) == 0 {
		// Did not find supported image directly in path and have no subdirectories to check
		return nil, nil
	} else {
		// Check subdirectories
		for _, subDir := range subdirs {
			var subDirPath = filepath.Join(path, subDir)
			res, err := findFirstDirContainingSupportedImage(subDirPath)
			if err != nil {
				log.Printf("Error while searching for supported images in subdirectory: %v", err)
			} else if res != nil {
				// Found supported image in subdirectory
				return res, nil
			}
		}
	}

	// Did not find supported image directly in path or in subdirectories
	return nil, nil
}
