package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

func setUpThemeListBox(currentTheme string) *gtk.ListBox {
	settings, _ := gtk.SettingsGetDefault()
	listBox, _ := gtk.ListBoxNew()
	var rowToSelect *gtk.ListBoxRow

	for _, name := range getThemeNames() {
		row, _ := gtk.ListBoxRowNew()

		eventBox, _ := gtk.EventBoxNew()
		box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		eventBox.Add(box)

		lbl, _ := gtk.LabelNew(name)
		lbl.SetProperty("margin-start", 6)
		lbl.SetProperty("margin-end", 6)
		n := name
		eventBox.Connect("button-press-event", func() {
			settings.SetProperty("gtk-theme-name", n)
			gtkSettings.themeName = n
		})
		row.Connect("focus-in-event", func() {
			settings.SetProperty("gtk-theme-name", n)
			gtkSettings.themeName = n
		})
		if n == currentTheme {
			rowToSelect = row
		}

		box.PackStart(lbl, false, false, 0)

		row.Add(eventBox)
		listBox.Add(row)
	}
	if rowToSelect != nil {
		listBox.SelectRow(rowToSelect)
		rowToFocus = rowToSelect
	}

	return listBox
}

func setUpIconThemeListBox(currentIconTheme string) *gtk.ListBox {
	settings, _ := gtk.SettingsGetDefault()
	listBox, _ := gtk.ListBoxNew()
	var rowToSelect *gtk.ListBoxRow

	// map[displayName]folderName
	namesMap := getIconThemeNames()
	var displayNames []string
	for name, _ := range namesMap {
		displayNames = append(displayNames, name)
	}
	sort.Slice(displayNames, func(i, j int) bool {
		return strings.ToUpper(displayNames[i]) < strings.ToUpper(displayNames[j])
	})

	for _, name := range displayNames {
		row, _ := gtk.ListBoxRowNew()

		eventBox, _ := gtk.EventBoxNew()
		box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		eventBox.Add(box)

		lbl, _ := gtk.LabelNew(name)
		lbl.SetProperty("margin-start", 6)
		lbl.SetProperty("margin-end", 6)
		n := name
		eventBox.Connect("button-press-event", func() {
			settings.SetProperty("gtk-icon-theme-name", namesMap[n])
			gtkSettings.iconThemeName = n
		})
		row.Connect("focus-in-event", func() {
			settings.SetProperty("gtk-icon-theme-name", namesMap[n])
			gtkSettings.iconThemeName = n
		})

		if namesMap[n] == currentIconTheme || n == currentIconTheme {
			rowToSelect = row
		}

		box.PackStart(lbl, false, false, 0)

		row.Add(eventBox)
		listBox.Add(row)
	}
	if rowToSelect != nil {
		listBox.SelectRow(rowToSelect)
		rowToFocus = rowToSelect
	}

	return listBox
}

func setUpCursorThemeListBox(currentCursorTheme string) *gtk.ListBox {
	settings, _ := gtk.SettingsGetDefault()
	listBox, _ := gtk.ListBoxNew()
	var rowToSelect *gtk.ListBoxRow

	var names []string
	for name, _ := range cursorThemeNames {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.ToUpper(names[i]) < strings.ToUpper(names[j])
	})

	for _, name := range names {
		row, _ := gtk.ListBoxRowNew()

		eventBox, _ := gtk.EventBoxNew()
		box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		eventBox.Add(box)

		lbl, _ := gtk.LabelNew(name)
		lbl.SetProperty("margin-start", 6)
		lbl.SetProperty("margin-end", 6)
		n := name
		eventBox.Connect("button-press-event", func() {
			settings.SetProperty("gtk-cursor-theme-name", cursorThemeNames[n])
			gtkSettings.cursorThemeName = cursorThemeNames[n]
			displayCursorThemes()
		})
		row.Connect("focus-in-event", func() {
			settings.SetProperty("gtk-cursor-theme-name", cursorThemeNames[n])
			gtkSettings.cursorThemeName = cursorThemeNames[n]
		})
		if cursorThemeNames[n] == currentCursorTheme {
			rowToSelect = row
		}

		box.PackStart(lbl, false, false, 0)

		row.Add(eventBox)
		listBox.Add(row)
	}
	if rowToSelect != nil {
		listBox.SelectRow(rowToSelect)
		rowToFocus = rowToSelect
	}

	return listBox
}

func setUpWidgetsPreview() *gtk.Frame {
	frame, _ := gtk.FrameNew("Widget style preview")
	frame.SetProperty("margin", 6)
	frame.SetProperty("valign", gtk.ALIGN_START)

	grid, _ := gtk.GridNew()
	grid.SetRowSpacing(6)
	grid.SetColumnSpacing(12)
	grid.SetProperty("margin", 6)
	frame.Add(grid)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	box.SetProperty("hexpand", true)
	grid.Attach(box, 0, 0, 3, 1)

	btn, _ := gtk.ButtonNewFromIconName("go-previous", gtk.ICON_SIZE_BUTTON)
	btn.SetProperty("can-focus", false)
	box.PackStart(btn, false, false, 0)

	btn, _ = gtk.ButtonNewFromIconName("go-next", gtk.ICON_SIZE_BUTTON)
	btn.SetProperty("can-focus", false)
	box.PackStart(btn, false, false, 0)

	btn, _ = gtk.ButtonNewFromIconName("process-stop", gtk.ICON_SIZE_BUTTON)
	btn.SetProperty("can-focus", false)
	box.PackStart(btn, false, false, 0)

	entry, _ := gtk.EntryNew()
	entry.SetProperty("can-focus", false)
	box.PackStart(entry, true, true, 0)

	checkButton, _ := gtk.CheckButtonNew()
	checkButton.SetProperty("can-focus", false)
	checkButton.SetLabel("Check Button")
	grid.Attach(checkButton, 0, 1, 1, 1)

	radioButton, _ := gtk.RadioButtonNew(nil)
	radioButton.SetProperty("can-focus", false)
	radioButton.SetLabel("Radio Button")
	grid.Attach(radioButton, 0, 2, 1, 1)

	spinButton, _ := gtk.SpinButtonNewWithRange(0, 1000, 10)
	spinButton.SetProperty("can-focus", false)
	grid.Attach(spinButton, 0, 3, 1, 1)

	button, _ := gtk.ButtonNewFromIconName("search", gtk.ICON_SIZE_BUTTON)
	button.SetProperty("can-focus", false)
	button.SetLabel("Button")
	grid.Attach(button, 1, 3, 1, 1)

	scale, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	scale.SetProperty("can-focus", false)
	scale.SetDrawValue(true)
	scale.SetValue(50)
	grid.Attach(scale, 1, 1, 2, 1)

	separator, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	separator.SetProperty("can-focus", false)
	separator.SetProperty("valign", gtk.ALIGN_CENTER)
	grid.Attach(separator, 1, 2, 2, 1)

	combo, _ := gtk.ComboBoxTextNew()
	combo.Append("entry #1", "entry #1")
	combo.Append("entry #2", "entry #2")
	combo.SetProperty("can-focus", false)
	grid.Attach(combo, 2, 3, 1, 1)

	return frame
}

func setUpIconsPreview() *gtk.Frame {
	frame, _ := gtk.FrameNew("Icon theme preview")
	frame.SetProperty("margin", 6)
	frame.SetProperty("valign", gtk.ALIGN_START)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 12)
	box.SetProperty("hexpand", true)
	frame.Add(box)

	flowBox, _ := gtk.FlowBoxNew()
	flowBox.SetMaxChildrenPerLine(7)
	flowBox.SetMinChildrenPerLine(7)
	box.PackStart(flowBox, false, false, 0)
	icons := []string{
		"user-home",
		"user-desktop",
		"folder",
		"folder-remote",
		"user-trash",
		"x-office-document",
		"application-x-executable",
		"image-x-generic",
		"package-x-generic",
		"emblem-mail",
		"utilities-terminal",
		"chromium",
		"firefox",
		"gimp"}
	for _, name := range icons {
		img, err := gtk.ImageNewFromIconName(name, gtk.ICON_SIZE_DIALOG)
		if err == nil {
			flowBox.Add(img)
			log.Debugf("Added icon: '%s'", name)
		} else {
			log.Warnf("Couldn't create image: '%s'", name)
		}
	}

	flowBox, _ = gtk.FlowBoxNew()
	box.PackStart(flowBox, false, false, 12)
	icons = []string{
		"network-wired-symbolic",
		"network-wireless-symbolic",
		"bluetooth-active-symbolic",
		"computer-symbolic",
		"audio-volume-high-symbolic",
		"battery-low-charging-symbolic",
		"display-brightness-medium-symbolic",
	}
	for _, name := range icons {
		img, err := gtk.ImageNewFromIconName(name, gtk.ICON_SIZE_MENU)
		if err == nil {
			flowBox.Add(img)
			log.Debugf("Added icon: '%s'", name)
		} else {
			log.Warnf("Couldn't create image: '%s'", name)
		}
	}

	return frame
}

func setUpCursorsPreview(path string) *gtk.Frame {
	frame, _ := gtk.FrameNew("Cursor theme preview")
	frame.SetProperty("margin", 6)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 12)
	box.SetProperty("margin", 12)
	box.SetProperty("hexpand", true)
	frame.Add(box)

	flowBox, _ := gtk.FlowBoxNew()
	flowBox.SetMaxChildrenPerLine(8)
	box.Add(flowBox)

	images := []string{
		"left_ptr",
		"hand2",
		"sb_v_double_arrow",
		"fleur",
		"xterm",
		"left_side",
		"top_left_corner",
		"h_double_arrow",
	}

	if path != "" {
		// As I have no better idea, we'll use the external `xcur2png` tool
		// to extract images from xcursor files, and save them to tmp dir.
		cursorsDir := filepath.Join(tempDir(), "nwg-look-cursors")

		dir, err := ioutil.ReadDir(cursorsDir)
		if err == nil {
			for _, d := range dir {
				os.RemoveAll(filepath.Join([]string{cursorsDir, d.Name()}...))
			}
		}
		// just in case it didn't yet exist
		makeDir(cursorsDir)

		for _, name := range images {
			imgPath := filepath.Join(path, name)

			args := []string{imgPath, "-d", cursorsDir, "-c", cursorsDir, "-q"}
			cmd := exec.Command("xcur2png", args...)

			cmd.Run()

			fName := fmt.Sprintf("%s_000.png", name)
			pngPath := filepath.Join(cursorsDir, fName)
			pixbuf, err := gdk.PixbufNewFromFileAtSize(pngPath, 24, 24)
			if err == nil {
				img, err := gtk.ImageNewFromPixbuf(pixbuf)
				if err == nil {
					flowBox.Add(img)
					p, _ := img.GetParent()
					parent, _ := p.(*gtk.FlowBoxChild)
					parent.SetProperty("can-focus", false)

					log.Debugf("Added icon: '%s'", pngPath)
				} else {
					log.Warnf("Couldn't create pixbuf from '%s'", pngPath)
				}
			} else {
				log.Warnf("Couldn't create image from '%s'", pngPath)
			}
		}
	}

	return frame
}

func setUpFontSelector(defaultFontName string) *gtk.Box {
	settings, _ := gtk.SettingsGetDefault()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)

	fontButton, _ := gtk.FontButtonNew()
	fontButton.SetProperty("valign", gtk.ALIGN_CENTER)
	fontButton.SetFont(defaultFontName)
	fontButton.Connect("font-set", func() {
		fontName := fontButton.GetFont()
		settings.SetProperty("gtk-font-name", fontName)
		gtkSettings.fontName = fontName
	})
	box.PackEnd(fontButton, true, true, 6)

	label, _ := gtk.LabelNew("Default font:")
	box.PackEnd(label, false, false, 6)

	return box
}
