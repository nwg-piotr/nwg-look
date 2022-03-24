package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.0.1"

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
	var debug = flag.Bool("d", false, "turn on Debug messages")
	var displayVersion = flag.Bool("v", false, "display Version information")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("nwg-look version %s\n", version)
		os.Exit(0)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	gtk.Init(nil)

	loadGtkSettings()

	builder, _ := gtk.BuilderNewFromFile("/home/piotr/Code/nwg-look/glade/main.glade")
	win, _ := getWindow(builder, "window")

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.Connect("key-release-event", func(window *gtk.Window, event *gdk.Event) bool {
		key := &gdk.EventKey{Event: event}
		if key.KeyVal() == gdk.KEY_Escape {
			gtk.MainQuit()
			return true
		}
		return false
	})

	viewport, _ := getViewPort(builder, "viewport-list")

	listBox = setUpThemeListBox(gtkSettings.themeName)
	viewport.Add(listBox)

	grid, _ := getGrid(builder, "grid")

	preview := setUpWidgetsPreview()
	grid.Attach(preview, 1, 1, 1, 1)

	fontSelector := setUpFontSelector(gtkSettings.fontName)
	fontSelector.SetProperty("vexpand", true)
	fontSelector.SetProperty("valign", gtk.ALIGN_START)
	grid.Attach(fontSelector, 1, 2, 1, 1)

	btnClose, _ := getButton(builder, "btn-close")
	btnClose.Connect("clicked", func() {
		gtk.MainQuit()
	})

	verLabel, _ := getLabel(builder, "version-label")
	verLabel.SetText(fmt.Sprintf("nwg-look v%s", version))

	win.ShowAll()

	gtk.Main()
}
