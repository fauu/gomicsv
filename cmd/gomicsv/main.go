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

package main

import (
	"fmt"
	"os"

	"github.com/fauu/gomicsv"
	flag "github.com/spf13/pflag"
)

var (
	buildDate     = ""
	versionString = ""

	referer = flag.String("referer", "", "HTTP Referer value to use for requests when the provided path is a URL")
	help    = flag.BoolP("help", "h", false, "Print usage message and exit")
	version = flag.BoolP("version", "v", false, "Print program version and exit")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [path] [options]:\n"+
				"      path               Path to the comic to load at startup. Could be filesystem path, a URL of one of the images or a URL template for all of the images\n"+
				"   Options:\n",
			gomicsv.AppName,
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Fprintf(os.Stderr, "%s version %s", gomicsv.AppNameDisplay, versionString)
		os.Exit(0)
	}

	nonFlagArgs := []string{os.Args[0]}
	nonFlagArgs = append(nonFlagArgs, flag.Args()...)
	initParams := gomicsv.AppStartupParams{
		Referer: *referer,
	}

	app := gomicsv.App{}
	gtkApplication := app.Init(nonFlagArgs, initParams, gomicsv.BuildInfo{Version: versionString, Date: buildDate})
	os.Exit(gtkApplication.Run(nonFlagArgs))
}
