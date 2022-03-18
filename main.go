package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	gtkSettings, _ := gtk.SettingsGetDefault()

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	win.Add(box)

	var names = [...]string{"Adwaita", "Adwaita-dark", "HighContrast", "Raleigh"}

	for _, name := range names {
		btn, _ := gtk.ButtonNew()
		n := name
		btn.SetLabel(name)
		btn.Connect("clicked", func() {
			gtkSettings.SetProperty("gtk-theme-name", n)
		})
		box.PackStart(btn, false, false, 10)
	}

	win.ShowAll()

	gtk.Main()
}
