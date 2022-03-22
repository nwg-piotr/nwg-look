// tools
package main

import (
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
