// tools
package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

func loadGtkSettings() {
	settings, _ := gtk.SettingsGetDefault()
	prop, _ := settings.GetProperty("gtk-theme-name")
	gtkSettings.themeName, _ = prop.(string)
	log.Infof("Current theme: %s", gtkSettings.themeName)

	prop, _ = settings.GetProperty("gtk-icon-theme-name")
	gtkSettings.iconThemeName, _ = prop.(string)
	log.Infof("Icon theme: %s", gtkSettings.iconThemeName)

	prop, _ = settings.GetProperty("gtk-font-name")
	gtkSettings.fontName, _ = prop.(string)
	log.Infof("Default font: %s", gtkSettings.fontName)

	prop, _ = settings.GetProperty("gtk-cursor-theme-name")
	gtkSettings.cursorThemeName, _ = prop.(string)
	log.Infof("Cursor theme: %s", gtkSettings.cursorThemeName)
}

func getAppDirs() []string {
	var dirs []string
	xdgDataDirs := ""

	home := os.Getenv("HOME")
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if os.Getenv("XDG_DATA_DIRS") != "" {
		xdgDataDirs = os.Getenv("XDG_DATA_DIRS")
	} else {
		xdgDataDirs = "/usr/local/share/:/usr/share/"
	}
	if xdgDataHome != "" {
		dirs = append(dirs, filepath.Join(xdgDataHome, "applications"))
	} else if home != "" {
		dirs = append(dirs, filepath.Join(home, ".local/share/applications"))
	}
	for _, d := range strings.Split(xdgDataDirs, ":") {
		dirs = append(dirs, filepath.Join(d, "applications"))
	}
	flatpakDirs := []string{filepath.Join(home, ".local/share/flatpak/exports/share/applications"),
		"/var/lib/flatpak/exports/share/applications"}

	for _, d := range flatpakDirs {
		if pathExists(d) && !isIn(dirs, d) {
			dirs = append(dirs, d)
		}
	}
	var confirmedDirs []string
	for _, d := range dirs {
		if pathExists(d) {
			confirmedDirs = append(confirmedDirs, d)
		}
	}
	return confirmedDirs
}

func isIn(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func pathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// getWindow returns *gtk.Window object from the glade resource
func getWindow(b *gtk.Builder, id string) (*gtk.Window, error) {

	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}

	window, ok := obj.(*gtk.Window)
	if !ok {
		return nil, err
	}

	return window, nil
}

func getScrolledWindow(b *gtk.Builder, id string) (*gtk.ScrolledWindow, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}

	window, ok := obj.(*gtk.ScrolledWindow)
	if !ok {
		return nil, err
	}

	return window, nil
}

func getViewPort(b *gtk.Builder, id string) (*gtk.Viewport, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}

	viewport, ok := obj.(*gtk.Viewport)
	if !ok {
		return nil, err
	}

	return viewport, nil
}

func getButton(b *gtk.Builder, id string) (*gtk.Button, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}

	btn, ok := obj.(*gtk.Button)
	if !ok {
		return nil, err
	}

	return btn, nil
}

func getGrid(b *gtk.Builder, id string) (*gtk.Grid, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}

	grid, ok := obj.(*gtk.Grid)
	if !ok {
		return nil, err
	}

	return grid, nil
}
