package main

import (
	"github.com/gotk3/gotk3/gtk"
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
		n := name
		eventBox.Connect("button-press-event", func() {
			settings.SetProperty("gtk-theme-name", n)
		})
		if n == currentTheme {
			rowToSelect = row
		}

		box.PackStart(lbl, false, false, 0)

		row.Add(eventBox)
		listBox.Add(row)
	}
	listBox.SelectRow(rowToSelect)

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
	box.PackStart(btn, false, false, 0)

	btn, _ = gtk.ButtonNewFromIconName("go-next", gtk.ICON_SIZE_BUTTON)
	box.PackStart(btn, false, false, 0)

	btn, _ = gtk.ButtonNewFromIconName("process-stop", gtk.ICON_SIZE_BUTTON)
	box.PackStart(btn, false, false, 0)

	entry, _ := gtk.EntryNew()
	box.PackStart(entry, true, true, 0)

	checkButton, _ := gtk.CheckButtonNew()
	checkButton.SetLabel("Check Button")
	grid.Attach(checkButton, 0, 1, 1, 1)

	radioButton, _ := gtk.RadioButtonNew(nil)
	radioButton.SetLabel("Radio Button")
	grid.Attach(radioButton, 0, 2, 1, 1)

	spinButton, _ := gtk.SpinButtonNewWithRange(0, 1000, 10)
	grid.Attach(spinButton, 0, 3, 1, 1)

	button, _ := gtk.ButtonNewFromIconName("search", gtk.ICON_SIZE_BUTTON)
	button.SetLabel("Button")
	grid.Attach(button, 1, 3, 1, 1)

	scale, _ := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	scale.SetDrawValue(true)
	scale.SetValue(60)
	grid.Attach(scale, 1, 1, 2, 1)

	separator, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	separator.SetProperty("valign", gtk.ALIGN_CENTER)
	grid.Attach(separator, 1, 2, 2, 1)

	combo, _ := gtk.ComboBoxTextNew()
	combo.Append("entry #1", "entry #1")
	combo.Append("entry #2", "entry #2")
	grid.Attach(combo, 2, 3, 1, 1)

	return frame
}

func setUpFontSelector(defaultFontName string) *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)

	fontButton, _ := gtk.FontButtonNew()
	fontButton.SetProperty("valign", gtk.ALIGN_CENTER)
	fontButton.SetFont(defaultFontName)
	box.PackEnd(fontButton, true, true, 6)

	label, _ := gtk.LabelNew("Default font:")
	box.PackEnd(label, false, false, 6)

	return box
}
