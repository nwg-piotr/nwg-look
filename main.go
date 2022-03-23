package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

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
)

type gtkSettingsFields struct {
	themeName       string
	iconThemeName   string
	fontName        string
	cursorThemeName string
	cursorThemeSize int
}

var gtkSettings gtkSettingsFields

func main() {
	log.SetLevel(log.DebugLevel)

	gtk.Init(nil)

	loadGtkSettings()

	builder, _ := gtk.BuilderNewFromFile("/home/piotr/Code/nwg-look/glade/main.glade")
	win, _ := getWindow(builder, "window")

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	viewport, _ := getViewPort(builder, "viewport-list")

	listBox = setUpThemeListBox(gtkSettings.themeName)

	viewport.Add(listBox)

	grid, _ := getGrid(builder, "grid")

	preview := setUpWidgetsPreview()

	grid.Attach(preview, 1, 1, 1, 1)

	// fontSelector := setUpFontSelector(gtkSettings.fontName)
	// fontSelector.SetProperty("vexpand", true)
	// fontSelector.SetProperty("valign", gtk.ALIGN_START)
	// grid.Attach(fontSelector, 1, 2, 1, 1)

	fontChooser, err := getFontChooser(builder, "font-chooser")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(">>>", fontChooser)
	fontChooser.SetFont(gtkSettings.fontName)

	btnClose, _ := getButton(builder, "btn-close")
	btnClose.Connect("clicked", func() {
		gtk.MainQuit()
	})

	win.ShowAll()

	gtk.Main()
}
