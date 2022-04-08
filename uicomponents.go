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
			gtkSettings.SetProperty("gtk-theme-name", n)
			gsettings.gtkTheme = n
		})
		row.Connect("focus-in-event", func() {
			gtkSettings.SetProperty("gtk-theme-name", n)
			gsettings.gtkTheme = n
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
	listBox, _ := gtk.ListBoxNew()
	var rowToSelect *gtk.ListBoxRow

	// map[displayName]folderName
	namesMap := getIconThemeNames()
	var displayNames []string
	for name := range namesMap {
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
			gtkSettings.SetProperty("gtk-icon-theme-name", namesMap[n])
			gsettings.iconTheme = n
		})
		row.Connect("focus-in-event", func() {
			gtkSettings.SetProperty("gtk-icon-theme-name", namesMap[n])
			gsettings.iconTheme = n
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
	listBox, _ := gtk.ListBoxNew()
	var rowToSelect *gtk.ListBoxRow

	var names []string
	for name := range cursorThemeNames {
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
			gtkSettings.SetProperty("gtk-cursor-theme-name", cursorThemeNames[n])
			gsettings.cursorTheme = cursorThemeNames[n]
			displayCursorThemes()
		})
		row.Connect("focus-in-event", func() {
			gtkSettings.SetProperty("gtk-cursor-theme-name", cursorThemeNames[n])
			gsettings.cursorTheme = cursorThemeNames[n]
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

	progressBar, _ := gtk.ProgressBarNew()
	progressBar.SetFraction(0.3)
	progressBar.SetText("30%")
	progressBar.SetShowText(true)
	progressBar.SetProperty("margin-bottom", 6)
	grid.Attach(progressBar, 0, 4, 3, 1)

	return frame
}

func setUpThemeSettingsForm(defaultFontName string) *gtk.Grid {
	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(12)
	grid.SetRowSpacing(6)
	grid.SetProperty("margin", 12)
	label, _ := gtk.LabelNew("Default font:")
	label.SetProperty("halign", gtk.ALIGN_END)
	grid.Attach(label, 0, 0, 1, 1)

	fontButton, _ := gtk.FontButtonNew()
	fontButton.SetProperty("valign", gtk.ALIGN_CENTER)
	fontButton.SetFont(defaultFontName)
	fontButton.Connect("font-set", func() {
		fontName := fontButton.GetFont()
		gtkSettings.SetProperty("gtk-font-name", fontName)
		gsettings.fontName = fontName
	})
	grid.Attach(fontButton, 1, 0, 1, 1)

	label, _ = gtk.LabelNew("Default font:")
	label.SetProperty("halign", gtk.ALIGN_END)
	grid.Attach(label, 1, 0, 1, 1)

	label, _ = gtk.LabelNew("Color scheme:")
	grid.Attach(label, 0, 1, 1, 1)

	combo, _ := gtk.ComboBoxTextNew()
	combo.Append("default", "default")
	combo.Append("prefer-dark", "prefer dark")
	combo.Append("prefer-light", "prefer light")
	combo.SetActiveID(gsettings.colorScheme)
	combo.SetProperty("can-focus", false)
	combo.Connect("changed", func() {
		id := combo.GetActiveID()
		gsettings.colorScheme = id
		if id == "prefer-dark" {
			gtkConfig.applicationPreferDarkTheme = true
			gtkSettings.SetProperty("gtk-application-prefer-dark-theme", true)
		} else {
			gtkConfig.applicationPreferDarkTheme = false
			gtkSettings.SetProperty("gtk-application-prefer-dark-theme", false)
		}
	})
	grid.Attach(combo, 1, 1, 1, 1)

	return grid
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

func setUpCursorSizeSelector() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	box.SetProperty("margin", 12)
	box.SetProperty("hexpand", true)
	box.SetProperty("vexpand", true)
	box.SetProperty("valign", gtk.ALIGN_START)

	lbl, _ := gtk.LabelNew("Cursor size:")
	box.PackStart(lbl, false, false, 0)

	sb, _ := gtk.SpinButtonNewWithRange(6, 1024, 1)
	sb.SetValue(float64(gsettings.cursorSize))
	sb.Connect("value-changed", func() {
		v := int(sb.GetValue())
		gtkSettings.SetProperty("gtk-cursor-theme-size", v)
		gsettings.cursorSize = v
	})
	box.PackStart(sb, false, false, 6)
	lbl, _ = gtk.LabelNew("(default: 24)")
	box.PackStart(lbl, false, false, 0)

	return box
}

func setUpFontSettingsForm() *gtk.Frame {
	// We wont be applying these properties to gtk.Settings for preview,
	// as they remain unchanged in once open window.

	frame, _ := gtk.FrameNew("Font settings")
	frame.SetProperty("margin", 6)
	g, _ := gtk.GridNew()
	g.SetRowSpacing(12)
	g.SetColumnSpacing(12)
	g.SetProperty("margin", 6)
	g.SetProperty("hexpand", true)
	g.SetProperty("vexpand", true)
	frame.Add(g)

	lbl, _ := gtk.LabelNew("Font hinting:")
	lbl.SetProperty("halign", gtk.ALIGN_END)
	g.Attach(lbl, 0, 0, 1, 1)

	comboHinting, _ := gtk.ComboBoxTextNew()
	comboHinting.Append("none", "none")
	comboHinting.Append("slight", "slight")
	comboHinting.Append("medium", "medium")
	comboHinting.Append("full", "full")
	comboHinting.SetActiveID(gsettings.fontHinting)
	g.Attach(comboHinting, 1, 0, 1, 1)

	comboHinting.Connect("changed", func() {
		id := comboHinting.GetActiveID()
		gsettings.fontHinting = id
	})

	lbl, _ = gtk.LabelNew("Font antialiasing:")
	lbl.SetProperty("halign", gtk.ALIGN_END)
	g.Attach(lbl, 0, 1, 1, 1)

	comboRgba, _ := gtk.ComboBoxTextNew()

	comboAntialiasing, _ := gtk.ComboBoxTextNew()
	comboAntialiasing.Append("none", "none")
	comboAntialiasing.Append("grayscale", "grayscale")
	comboAntialiasing.Append("rgba", "rgba")
	comboAntialiasing.SetActiveID(gsettings.fontAntialiasing)
	g.Attach(comboAntialiasing, 1, 1, 1, 1)

	comboAntialiasing.Connect("changed", func() {
		id := comboAntialiasing.GetActiveID()
		gsettings.fontAntialiasing = id
		comboRgba.SetSensitive(id == "rgba")
	})

	lbl, _ = gtk.LabelNew("Font RGBA order:")
	lbl.SetProperty("halign", gtk.ALIGN_END)
	g.Attach(lbl, 0, 2, 1, 1)

	comboRgba.Append("rgb", "RGB")
	comboRgba.Append("bgr", "BGR")
	comboRgba.Append("vrgb", "VRGB")
	comboRgba.Append("vbgr", "VBGR")
	comboRgba.SetActiveID(gsettings.fontRgbaOrder)
	comboRgba.SetSensitive(comboAntialiasing.GetActiveID() == "rgba")
	g.Attach(comboRgba, 1, 2, 1, 1)

	comboRgba.Connect("changed", func() {
		gsettings.fontRgbaOrder = comboRgba.GetActiveID()
	})

	lbl, _ = gtk.LabelNew("Text scaling factor:")
	lbl.SetProperty("halign", gtk.ALIGN_END)
	g.Attach(lbl, 0, 3, 1, 1)

	sb, _ := gtk.SpinButtonNewWithRange(0.5, 3, 0.01)
	sb.SetValue(gsettings.textScalingFactor)
	sb.Connect("value-changed", func() {
		v := sb.GetValue()
		gsettings.textScalingFactor = v
	})
	g.Attach(sb, 1, 3, 1, 1)

	return frame
}

func setUpOtherSettingsForm() *gtk.Frame {
	// We won't be applying these properties to gtk.Settings for preview,
	// as they remain unchanged in once open window.

	frame, _ := gtk.FrameNew("Other settings")
	frame.SetProperty("margin", 6)
	g, _ := gtk.GridNew()
	g.SetRowSpacing(12)
	g.SetColumnSpacing(12)
	g.SetProperty("margin", 6)
	g.SetProperty("hexpand", true)
	g.SetProperty("vexpand", true)
	frame.Add(g)

	lbl, _ := gtk.LabelNew("")
	lbl.SetMarkup("<b>UI settings</b> (deprecated)")
	lbl.SetProperty("halign", gtk.ALIGN_START)
	g.Attach(lbl, 0, 0, 1, 1)

	lbl, _ = gtk.LabelNew("Toolbar style:")
	lbl.SetProperty("halign", gtk.ALIGN_END)
	g.Attach(lbl, 0, 1, 1, 1)

	comboToolbarStyle, _ := gtk.ComboBoxTextNew()
	comboToolbarStyle.SetTooltipText("deprecated since GTK 3.10, ignored")
	comboToolbarStyle.Append("both", "Text below icons")
	comboToolbarStyle.Append("both-horiz", "Text next to icons")
	comboToolbarStyle.Append("icons", "Icons")
	comboToolbarStyle.Append("text", "Text")
	comboToolbarStyle.SetActiveID(gsettings.toolbarStyle)
	g.Attach(comboToolbarStyle, 1, 1, 1, 1)

	comboToolbarStyle.Connect("changed", func() {
		gsettings.toolbarStyle = comboToolbarStyle.GetActiveID()

		switch gsettings.toolbarStyle {
		case "both":
			gtkConfig.toolbarStyle = "GTK_TOOLBAR_BOTH"
		case "icons":
			gtkConfig.toolbarStyle = "GTK_TOOLBAR_ICONS"
		case "text":
			gtkConfig.toolbarStyle = "GTK_TOOLBAR_TEXT"
		default:
			gtkConfig.toolbarStyle = "GTK_TOOLBAR_BOTH_HORIZ"
		}
	})

	lbl, _ = gtk.LabelNew("Toolbar icon size:")
	lbl.SetProperty("halign", gtk.ALIGN_END)
	g.Attach(lbl, 0, 2, 1, 1)

	comboToolbarIconSize, _ := gtk.ComboBoxTextNew()
	comboToolbarIconSize.SetTooltipText("deprecated since GTK 3.10, ignored")
	comboToolbarIconSize.Append("small", "Small")
	comboToolbarIconSize.Append("large", "Large")
	comboToolbarIconSize.SetActiveID(gsettings.toolbarIconsSize)
	g.Attach(comboToolbarIconSize, 1, 2, 1, 1)

	comboToolbarIconSize.Connect("changed", func() {
		gsettings.toolbarIconsSize = comboToolbarIconSize.GetActiveID()

		if gsettings.toolbarIconsSize == "small" {
			gtkConfig.toolbarIconSize = "GTK_ICON_SIZE_SMALL_TOOLBAR"
		} else {
			gtkConfig.toolbarIconSize = "GTK_ICON_SIZE_LARGE_TOOLBAR"
		}
	})

	cbBtn, _ := gtk.CheckButtonNewWithLabel("Show button images")
	cbBtn.SetTooltipText("deprecated since GTK 3.10")
	cbBtn.SetActive(gtkConfig.buttonImages)
	cbBtn.Connect("toggled", func() {
		gtkConfig.buttonImages = cbBtn.GetActive()
	})
	g.Attach(cbBtn, 0, 3, 1, 1)

	cbMnu, _ := gtk.CheckButtonNewWithLabel("Show menu images")
	cbMnu.SetTooltipText("deprecated since GTK 3.10")
	cbMnu.SetActive(gtkConfig.menuImages)
	cbMnu.Connect("toggled", func() {
		gtkConfig.menuImages = cbMnu.GetActive()
	})
	g.Attach(cbMnu, 0, 4, 1, 1)

	lbl, _ = gtk.LabelNew("")
	lbl.SetMarkup("<b>Sound effects</b>")
	lbl.SetProperty("halign", gtk.ALIGN_START)
	g.Attach(lbl, 0, 5, 1, 1)

	cbEventSounds, _ := gtk.CheckButtonNewWithLabel("Enable event sounds")
	cbEventSounds.SetActive(gsettings.eventSounds)
	cbEventSounds.Connect("toggled", func() {
		gsettings.eventSounds = cbEventSounds.GetActive()
		gtkConfig.enableEventSounds = cbEventSounds.GetActive()
	})
	g.Attach(cbEventSounds, 0, 6, 1, 1)

	cbInputSounds, _ := gtk.CheckButtonNewWithLabel("Enable input feedback sounds")
	cbInputSounds.SetActive(gsettings.inputFeedbackSounds)
	cbInputSounds.Connect("toggled", func() {
		gsettings.inputFeedbackSounds = cbInputSounds.GetActive()
		gtkConfig.enableInputFeedbackSounds = cbInputSounds.GetActive()
	})
	g.Attach(cbInputSounds, 0, 7, 2, 1)

	return frame
}
