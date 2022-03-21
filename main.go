package main

import (
	"fmt"
	// "log"

	"github.com/gotk3/gotk3/gtk"
)

var themeNames = [...]string{
	"Adapta",
	"Adapta-Eta",
	"Adwaita",
	"Adwaita-Dark-Green",
	"Adwaita-dark",
	"Aero",
	"Aero-dark",
	"Aquatix",
	"ArchLabs-Dark",
	"ArchLabs-Light",
}

var (
	listBox *gtk.ListBox
	margins = [...]string{"margin-start", "margin-end", "margin-top", "margin-bottom"}
)

func main() {
	gtk.Init(nil)

	builder, _ := gtk.BuilderNewFromFile("/home/piotr/Code/nwg-look/glade/main.glade")
	win, _ := getWindow(builder, "window")

	gtkSettings, _ := gtk.SettingsGetDefault()
	prop, _ := gtkSettings.GetProperty("gtk-theme-name")
	currentTheme, _ := prop.(string)

	fmt.Println("Current theme:", currentTheme)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	viewport, _ := getViewPort(builder, "viewport-list")

	listBox = setUpThemeListBox(gtkSettings, currentTheme)

	viewport.Add(listBox)

	grid, _ := getGrid(builder, "grid")

	preview := setUpWidgetsPreview()

	grid.Attach(preview, 1, 1, 1, 1)

	for _, prop := range margins {
		preview.SetProperty(prop, 6)
	}

	btnClose, _ := getButton(builder, "btn-close")
	btnClose.Connect("clicked", func() {
		gtk.MainQuit()
	})

	win.ShowAll()

	gtk.Main()
}
