/*
GTK3 settings editor adapted to work in the sway / wlroots environment
Project: https://github.com/nwg-piotr/nwg-look
Author's email: nwg.piotr@gmail.com
Copyright (c) 2022 Piotr Miller
License: MIT
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.1.1"

var (
	originalGtkConfig     []string // we will append not parsed settings.ini lines from here
	gtkConfig             gtkConfigProperties
	gtkSettings           *gtk.Settings
	gsettings             gsettingsValues
	dataDirs              []string
	cursorThemes          map[string]string // theme name to path
	cursorThemeNames      map[string]string // theme name to theme folder name
	viewport              *gtk.Viewport
	scrolledWindow        *gtk.ScrolledWindow
	listBox               *gtk.ListBox
	menuBar               *gtk.MenuBar
	themeSettingsSelector *gtk.Grid
	grid                  *gtk.Grid
	preview               *gtk.Frame
	cursorSizeSelector    *gtk.Box
	fontSettingsForm      *gtk.Frame
	rowToFocus            *gtk.ListBoxRow
)

type gtkConfigProperties struct {
	themeName                  string
	iconThemeName              string
	fontName                   string
	cursorThemeName            string
	cursorThemeSize            int
	toolbarStyle               string
	toolbarIconSize            string
	buttonImages               bool
	menuImages                 bool
	enableEventSounds          bool
	enableInputFeedbackSounds  bool
	xftAntialias               int
	fontAntialiasing           string
	xftDpi                     int
	xftHinting                 int
	xftHintstyle               string
	xftRgba                    string
	applicationPreferDarkTheme bool
}

type gsettingsValues struct {
	// org.gnome.desktop.interface
	gtkTheme          string
	iconTheme         string
	fontName          string
	cursorTheme       string
	cursorSize        int
	toolbarStyle      string
	toolbarIconsSize  string
	fontHinting       string
	fontAntialiasing  string
	fontRgbaOrder     string
	textScalingFactor float64
	colorScheme       string
	// org.gnome.desktop.sound
	eventSounds         bool
	inputFeedbackSounds bool
}

func gsettingsNewWithDefaults() gsettingsValues {
	g := gsettingsValues{}
	g.gtkTheme = "Adwaita"
	g.iconTheme = "Adwaita"
	g.fontName = "Sans 10"
	g.cursorTheme = "Adwaita"
	g.cursorSize = 24
	g.toolbarStyle = "both-horiz"
	g.toolbarIconsSize = "large"
	g.fontHinting = "medium"
	g.fontAntialiasing = "grayscale"
	g.fontRgbaOrder = "rgb"
	g.textScalingFactor = 1.0
	g.eventSounds = true
	g.inputFeedbackSounds = false
	g.colorScheme = "default"

	return g
}

func gtkConfigPropertiesNewWithDefaults() gtkConfigProperties {
	s := gtkConfigProperties{}
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
	s.applicationPreferDarkTheme = false

	val, err := getGsettingsValue("org.gnome.desktop.interface", "font-antialiasing")
	if err == nil {
		s.fontAntialiasing = val
	} else {
		log.Warn(err)
	}

	s.xftHinting = -1
	s.xftHintstyle = "hintmedium"
	s.xftRgba = "none"

	return s
}

func displayThemes() {
	destroyContent()
	rowToFocus = nil

	listBox = setUpThemeListBox(gsettings.gtkTheme)
	viewport.Add(listBox)
	menuBar.Deactivate()
	if rowToFocus != nil {
		rowToFocus.GrabFocus()
	}

	preview = setUpWidgetsPreview()
	grid.Attach(preview, 1, 1, 1, 1)

	themeSettingsSelector = setUpThemeSettingsForm(gsettings.fontName)
	themeSettingsSelector.SetProperty("vexpand", true)
	themeSettingsSelector.SetProperty("valign", gtk.ALIGN_START)
	grid.Attach(themeSettingsSelector, 1, 2, 1, 1)

	viewport.ShowAll()
	grid.ShowAll()
}

func displayIconThemes() {
	destroyContent()
	rowToFocus = nil

	listBox = setUpIconThemeListBox(gsettings.iconTheme)
	viewport.Add(listBox)
	menuBar.Deactivate()
	if rowToFocus != nil {
		rowToFocus.GrabFocus()
	}

	preview = setUpIconsPreview()
	grid.Attach(preview, 1, 1, 1, 1)

	viewport.ShowAll()
	grid.ShowAll()
}

func displayCursorThemes() {
	destroyContent()
	rowToFocus = nil

	listBox = setUpCursorThemeListBox(gsettings.cursorTheme)
	viewport.Add(listBox)
	menuBar.Deactivate()
	if rowToFocus != nil {
		rowToFocus.GrabFocus()
	}

	preview = setUpCursorsPreview(cursorThemes[gsettings.cursorTheme])
	grid.Attach(preview, 1, 1, 1, 1)

	cursorSizeSelector = setUpCursorSizeSelector()
	grid.Attach(cursorSizeSelector, 1, 2, 1, 1)

	viewport.ShowAll()
	grid.ShowAll()
}

func displayFontSettingsForm() {
	destroyContent()

	preview = setUpFontSettingsForm()
	grid.Attach(preview, 0, 1, 1, 1)
	menuBar.Deactivate()
	grid.ShowAll()
	scrolledWindow.Hide()
}

func displayOtherSettingsForm() {
	destroyContent()

	preview = setUpOtherSettingsForm()
	grid.Attach(preview, 0, 1, 1, 1)
	menuBar.Deactivate()
	grid.ShowAll()
	scrolledWindow.Hide()
}

func destroyContent() {
	if listBox != nil {
		listBox.Destroy()
	}
	if preview != nil {
		preview.Destroy()
	}
	if themeSettingsSelector != nil {
		themeSettingsSelector.Destroy()
	}
	if cursorSizeSelector != nil {
		cursorSizeSelector.Destroy()
	}
}

func main() {
	var debug = flag.Bool("d", false, "turn on Debug messages")
	var displayVersion = flag.Bool("v", false, "display Version information")
	var applyGs = flag.Bool("a", false, "Apply stored gsetting and quit")
	var doNotSave = flag.Bool("n", false, "do Not save gtk settings.ini")
	var restoreDefaults = flag.Bool("r", false, "Restore default values and quit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("nwg-look version %s\n", version)
		os.Exit(0)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// initialize gsettings type with default gtk values
	gsettings = gsettingsNewWithDefaults()

	// initialize gtkConfigProperties type with default gtk.Settings values
	gtkConfig = gtkConfigPropertiesNewWithDefaults()

	if *restoreDefaults {
		fmt.Print("Restore default gtk settings? y/N ")
		var input string
		fmt.Scanln(&input)
		fmt.Println(input)
		if strings.ToUpper(input) == "Y" {
			applyGsettings()
			saveGsettingsBackup()
			if !*doNotSave {
				saveGtkIni()
			}
		}
		os.Exit(0)
	}

	if *applyGs {
		applyGsettingsFromFile()
		os.Exit(0)
	}

	dataDirs = getDataDirs()
	cursorThemes, cursorThemeNames = getCursorThemes()

	gtk.Init(nil)

	// update gtkConfig from gtk-3.0/settings.ini
	if !*doNotSave {
		loadGtkConfig()
	}

	readGsettings()

	gtkSettings, _ = gtk.SettingsGetDefault()

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
	scrolledWindow, _ = getScrolledWindow(builder, "scrolled-window")
	grid, _ = getGrid(builder, "grid")

	menuBar, _ = getMenuBar(builder, "menubar")

	item1, _ := getMenuItem(builder, "item-widgets")
	item1.Connect("button-release-event", displayThemes)

	item2, _ := getMenuItem(builder, "item-icons")
	item2.Connect("button-release-event", displayIconThemes)

	item3, _ := getMenuItem(builder, "item-cursors")
	item3.Connect("button-release-event", displayCursorThemes)

	item4, _ := getMenuItem(builder, "item-font")
	item4.Connect("button-release-event", displayFontSettingsForm)

	item5, _ := getMenuItem(builder, "item-other")
	item5.Connect("button-release-event", displayOtherSettingsForm)

	btnClose, _ := getButton(builder, "btn-close")
	btnClose.Connect("clicked", func() {
		gtk.MainQuit()
	})

	btnApply, _ := getButton(builder, "btn-apply")
	btnApply.Connect("clicked", func() {
		applyGsettings()
		saveGsettingsBackup()
		if !*doNotSave {
			saveGtkIni()
		}
		saveIndexTheme()
	})

	verLabel, _ := getLabel(builder, "version-label")
	verLabel.SetMarkup(fmt.Sprintf("<b>nwg-look</b> v%s <a href='https://github.com/nwg-piotr/nwg-look'>GitHub</a>", version))

	displayThemes()

	win.ShowAll()

	gtk.Main()
}
