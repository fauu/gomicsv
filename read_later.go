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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fauu/gomicsv/archive"
)

func (app *App) maybeSaveReadingPosition() {
	var http = app.S.Archive.Kind() == archive.HTTPKind
	if (!http && !app.Config.RememberPosition) ||
		(http && !app.Config.RememberPositionHTTP) {
		return
	}

	readLaterFilePath := app.readLaterFilePath(app.S.ArchivePath)
	f, err := os.Create(readLaterFilePath)
	if err != nil {
		log.Printf("Error creating read later file '%s': %v", readLaterFilePath, err)
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "archive-pos=%d", app.S.ArchivePos)
}

func (app *App) loadReadingPosition(archivePath string) (int, error) {
	readLaterFilePath := app.readLaterFilePath(archivePath)
	f, err := os.Open(readLaterFilePath)
	if err != nil {
		return 0, fmt.Errorf("opening read later file '%s': %v", readLaterFilePath, err)
	}
	defer f.Close()

	var archivePos int
	_, err = fmt.Fscanf(f, "archive-pos=%d", &archivePos)
	if err != nil {
		return 0, fmt.Errorf("reading 'archive-pos' from read later file '%s': %v", readLaterFilePath, err)
	}
	return archivePos, nil
}

func (app *App) readLaterFilePath(archivePath string) string {
	filename := strings.ToUpper(md5String(archivePath))
	return filepath.Join(app.S.ReadLaterDirPath, filename)
}
