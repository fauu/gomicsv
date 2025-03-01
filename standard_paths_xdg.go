//go:build freebsd || linux || netbsd || openbsd || solaris

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
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func getConfigLocation(appName string) (string, error) {
	if appName == "" {
		return "", errors.New("'appName' cannot be empty")
	}
	if envHome := os.Getenv("XDG_CONFIG_HOME"); envHome != "" {
		return filepath.Join(envHome, appName), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting user home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config", appName), nil
}

func getUserDataLocation(appName string) (string, error) {
	if appName == "" {
		return "", errors.New("'appName' cannot be empty")
	}
	if envHome := os.Getenv("XDG_DATA_HOME"); envHome != "" {
		return filepath.Join(envHome, appName), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting user home directory: %v", err)
	}
	dataRoot := filepath.Join(homeDir, ".local", "share")
	if unix.Access(dataRoot, unix.W_OK) != nil {
		return "", fmt.Errorf("accessing user data directory '%s'", dataRoot)
	}
	return filepath.Join(dataRoot, appName), nil
}
