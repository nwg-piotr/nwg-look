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
	dataDirs         []string
	cursorThemes     map[string]string
	cursorThemeNames map[string]string
	viewport         *gtk.Viewport
	listBox          *gtk.ListBox
	menuBar          *gtk.MenuBar
	fontSelector     *gtk.Box
	grid             *gtk.Grid
	preview          *gtk.Frame
	rowToFocus       *gtk.ListBoxRow
)

type gtkSettingsFields struct {
	themeName                 string
	iconThemeName             string
	fontName                  string
	cursorThemeName           string
	cursorThemeSize           int
	toolbarStyle              string
	toolbarIconSize           string
	buttonImages              bool
	menuImages                bool
	enableEventSounds         bool
	enableInputFeedbackSounds bool
	xftAntialias              int
	xftDpi                    int
	xftHinting                int
	xftHintstyle              string
	xftRgba                   string
}

func gtkSettingsFieldsDefault() gtkSettingsFields {
	s := gtkSettingsFields{}
	// 'ignored' and 'deprecated' values left for lxappearance compatibility
	s.themeName = "Adwaita"
	s.iconThemeName = "Adwaita"
	s.fontName = "Sans 10"
	s.cursorThemeName = ""
	s.cursorThemeSize = 0
	s.toolbarStyle = "GTK_TOOLBAR_ICONS"              // ignored
	s.toolbarIconSize = "GTK_ICON_SIZE_LARGE_TOOLBAR" // ignored
	s.buttonImages = false                            // deprecated
	s.menuImages = false                              // deprecated
	s.enableEventSounds = true
	s.enableInputFeedbackSounds = true
	s.xftAntialias = -1
	s.xftDpi = -1
	s.xftHinting = -1
	s.xftHintstyle = ""
	s.xftRgba = ""

	return s
}

var gtkSettings gtkSettingsFields
var originalGtkSettings []string

func displayThemes() {
	if listBox != nil {
		listBox.Destroy()
	}
	listBox = setUpThemeListBox(gtkSettings.themeName)
	viewport.Add(listBox)
	menuBar.Deactivate()
	rowToFocus.GrabFocus()

	if preview != nil {
		preview.Destroy()
	}
	preview = setUpWidgetsPreview()
	grid.Attach(preview, 1, 1, 1, 1)

	if fontSelector != nil {
		fontSelector.Destroy()
	}
	fontSelector = setUpFontSelector(gtkSettings.fontName)
	fontSelector.SetProperty("vexpand", true)
	fontSelector.SetProperty("valign", gtk.ALIGN_START)
	grid.Attach(fontSelector, 1, 2, 1, 1)

	viewport.ShowAll()
	grid.ShowAll()
}

func displayIconThemes() {
	if listBox != nil {
		listBox.Destroy()
	}
	listBox = setUpIconThemeListBox(gtkSettings.iconThemeName)
	viewport.Add(listBox)
	menuBar.Deactivate()
	rowToFocus.GrabFocus()

	if preview != nil {
		preview.Destroy()
	}
	preview = setUpIconsPreview()
	grid.Attach(preview, 1, 1, 1, 1)

	if fontSelector != nil {
		fontSelector.Destroy()
	}

	viewport.ShowAll()
	grid.ShowAll()
}

func displayCursorThemes() {
	if listBox != nil {
		listBox.Destroy()
	}
	listBox = setUpCursorThemeListBox(gtkSettings.cursorThemeName)
	viewport.Add(listBox)
	menuBar.Deactivate()
	rowToFocus.GrabFocus()

	if preview != nil {
		preview.Destroy()
	}

	preview = setUpCursorsPreview(cursorThemes[gtkSettings.cursorThemeName])
	grid.Attach(preview, 1, 1, 1, 1)

	if fontSelector != nil {
		fontSelector.Destroy()
	}

	viewport.ShowAll()
	grid.ShowAll()
}

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

	dataDirs = getDataDirs()
	cursorThemes, cursorThemeNames = getCursorThemes()

	gtk.Init(nil)

	gtkSettings = gtkSettingsFieldsDefault()

	loadGtkSettings()

	builder, _ := gtk.BuilderNewFromFile("/usr/share/nwg-look/main.glade")
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

	viewport, _ = getViewPort(builder, "viewport-list")
	grid, _ = getGrid(builder, "grid")

	menuBar, _ = getMenuBar(builder, "menubar")

	item1, _ := getMenuItem(builder, "item-widgets")
	item1.Connect("button-release-event", displayThemes)

	item2, _ := getMenuItem(builder, "item-icons")
	item2.Connect("button-release-event", displayIconThemes)

	item3, _ := getMenuItem(builder, "item-cursors")
	item3.Connect("button-release-event", displayCursorThemes)

	btnClose, _ := getButton(builder, "btn-close")
	btnClose.Connect("clicked", func() {
		gtk.MainQuit()
	})

	btnApply, _ := getButton(builder, "btn-apply")
	btnApply.SetSensitive(false)

	verLabel, _ := getLabel(builder, "version-label")
	verLabel.SetText(fmt.Sprintf("nwg-look v%s", version))

	displayThemes()

	win.ShowAll()

	gtk.Main()
}
