// tools
package main

import (
	"github.com/gotk3/gotk3/gtk"
)

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
