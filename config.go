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

package gomicsv

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	ConfigFilename = "config"
	ReadLaterDir   = "read-later"
)

type Config struct {
	ZoomMode             string
	Enlarge              bool
	Shrink               bool
	LastDirectory        string
	Fullscreen           bool
	HideUI               bool
	WindowWidth          int
	WindowHeight         int
	Random               bool
	Seamless             bool
	HFlip                bool
	VFlip                bool
	DoublePage           bool
	MangaMode            bool
	BackgroundColor      Color
	NSkip                int
	NPreload             int
	RememberRecent       bool
	RememberPosition     bool
	RememberPositionHTTP bool
	OneWide              bool
	EmbeddedOrientation  bool
	Interpolation        int
	ImageDiffThres       float32
	SceneScanSkip        int
	SmartScroll          bool
	HideIdleCursor       bool
	KamiteEnabled        bool
	KamitePort           int
	Bookmarks            []Bookmark
}

func (app *App) configFilePath() string {
	return filepath.Join(app.S.ConfigDirPath, ConfigFilename)
}

func (app *App) loadConfig() {
	app.Config.setDefaults()

	currentUser, err := user.Current()
	if err != nil {
		log.Panicf("getting the current user: %v", err)
	}
	app.Config.LastDirectory = currentUser.HomeDir

	if err := app.Config.load(app.configFilePath()); err != nil {
		if !os.IsNotExist(err) {
			log.Panicf("loading config: %v", err)
		}
	}
}

func (app *App) saveConfig() {
	if err := app.Config.save(app.configFilePath()); err != nil {
		log.Printf("Error saving config: %v", err)
	}
}

func (c *Config) setDefaults() {
	c.ZoomMode = "FitToWidth"
	c.Shrink = true
	c.Enlarge = true
	c.WindowWidth = 640
	c.WindowHeight = 480
	c.NSkip = 10
	c.NPreload = 2
	c.Seamless = true
	c.RememberRecent = true
	c.RememberPosition = false
	c.RememberPositionHTTP = false
	c.BackgroundColor = Color{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}
	c.Interpolation = 2
	c.EmbeddedOrientation = true
	c.ImageDiffThres = 0.4
	c.SceneScanSkip = 5
	c.SmartScroll = false
	c.HideIdleCursor = true
	c.KamiteEnabled = false
	c.KamitePort = 4110
}

func (c *Config) load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	if err = d.Decode(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) save(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (app *App) setZoomMode(mode string) {
	switch mode {
	case "FitToWidth":
		app.W.MenuItemFitToWidth.SetActive(true)
	case "FitToHeight":
		app.W.MenuItemFitToHeight.SetActive(true)
	case "BestFit":
		app.W.MenuItemBestFit.SetActive(true)
	default:
		app.W.MenuItemOriginal.SetActive(true)
		mode = "Original"
	}

	app.Config.ZoomMode = mode
	app.blit()
	app.updateStatus()
}

func (app *App) setHFlip(hflip bool) {
	app.Config.HFlip = hflip
	app.blit()
}

func (app *App) setVFlip(vflip bool) {
	app.Config.VFlip = vflip
	app.blit()
}

func (app *App) setFullscreen(fullscreen bool) {
	app.Config.Fullscreen = fullscreen

	if fullscreen {
		app.W.MainWindow.Fullscreen()
	} else {
		app.W.MainWindow.Unfullscreen()
	}

	app.W.MenuItemFullscreen.SetActive(fullscreen)
}

func (app *App) setHideUI(hideUI bool) {
	app.Config.HideUI = hideUI

	if hideUI {
		app.W.Menubar.Hide()
		app.W.Toolbar.Hide()
		app.W.Statusbar.Hide()
	} else {
		app.W.Menubar.Show()
		app.W.Toolbar.Show()
		app.W.Statusbar.Show()
	}

	app.W.MenuItemHideUI.SetActive(hideUI)
}

func (app *App) setShrink(shrink bool) {
	app.Config.Shrink = shrink
	app.W.MenuItemShrink.SetActive(shrink)
	app.blit()
	app.updateStatus()
}

func (app *App) setEnlarge(enlarge bool) {
	app.Config.Enlarge = enlarge
	app.W.MenuItemEnlarge.SetActive(enlarge)
	app.blit()
	app.updateStatus()
}

func (app *App) setRandom(random bool) {
	app.Config.Random = random
	app.W.MenuItemRandom.SetActive(random)
}

func (app *App) setSeamless(seamless bool) {
	app.Config.Seamless = seamless
	app.W.MenuItemSeamless.SetActive(seamless)
}

func (app *App) setDoublePage(doublePage bool) {
	app.W.ImageR.SetVisible(doublePage)
	app.Config.DoublePage = doublePage
	app.doSetPage(app.S.ArchivePos)
}

func (app *App) setMangaMode(mangaMode bool) {
	app.Config.MangaMode = mangaMode
	app.blit()
	app.updateStatus()
}

func (app *App) setBackgroundColor(color Color) {
	app.Config.BackgroundColor = color

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.Panicf("getting default screen: %v", err)
	}

	if app.S.BackgroundColorCssProvider != nil {
		gtk.RemoveProviderForScreen(screen, app.S.BackgroundColorCssProvider)
		app.S.BackgroundColorCssProvider = nil
	}

	provider, err := gtk.CssProviderNew()
	if err != nil {
		log.Panicf("creating CssProvider: %v", err)
	}
	err = provider.LoadFromData(fmt.Sprintf("#ScrolledWindow { background-color: %s; }", color.ToCSS()))
	if err != nil {
		log.Panicf("adding css to provider: %v", err)
	}
	gtk.AddProviderForScreen(screen, provider, 1)
	app.S.BackgroundColorCssProvider = provider
}

func (app *App) setInterpolation(interpolation int) {
	app.Config.Interpolation = interpolation
	app.blit()
}

func (app *App) setRememberRecent(rememberRecent bool) {
	app.Config.RememberRecent = rememberRecent
}

func (app *App) setRememberPosition(rememberPosition bool) {
	app.Config.RememberPosition = rememberPosition
	// TODO: Clear remembered when unset
}

func (app *App) setRememberPositionHTTP(rememberPositionHTTP bool) {
	app.Config.RememberPositionHTTP = rememberPositionHTTP
}

func (app *App) setOneWide(oneWide bool) {
	app.Config.OneWide = oneWide
	app.blit()
	app.updateStatus()
}

func (app *App) setSmartScroll(smartScroll bool) {
	app.Config.SmartScroll = smartScroll
}

func (app *App) setHideIdleCursor(hideIdleCursor bool) {
	app.Config.HideIdleCursor = hideIdleCursor
}

func (app *App) setEmbeddedOrientation(embeddedOrientation bool) {
	app.Config.EmbeddedOrientation = embeddedOrientation
	app.blit()
	app.updateStatus()
}

func (app *App) setKamiteEnabled(kamiteEnable bool) {
	app.Config.KamiteEnabled = kamiteEnable
}

func (app *App) setKamitePort(kamitePort int) {
	app.Config.KamitePort = kamitePort
}
