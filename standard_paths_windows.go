//go:build windows

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
	"errors"
	"os"
	"path/filepath"
)

func getConfigLocation(appName string) (string, error) {
	return getUserDataLocation(appName)
}

func getUserDataLocation(appName string) (string, error) {
	if appName == "" {
		return "", errors.New("'appName' is empty")
	}
	if envAppdata := os.Getenv("APPDATA"); envAppdata != "" {
		return filepath.Join(envAppdata, appName), nil
	}
	return "", errors.New("getting user appdata directory")
}
