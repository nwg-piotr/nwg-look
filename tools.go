// tools
package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

func configHome() string {
	cHome := os.Getenv("XDG_CONFIG_HOME")
	if cHome != "" {
		return cHome
	}
	return filepath.Join(os.Getenv("HOME"), ".config/")
}

func loadGtkSettings() {
	// settings, _ := gtk.SettingsGetDefault()
	// prop, _ := settings.GetProperty("gtk-theme-name")
	// gtkSettings.themeName, _ = prop.(string)
	// log.Infof("Current theme: %s", gtkSettings.themeName)

	// parse gtk settings file
	originalGtkConfig = []string{}
	configFile := filepath.Join(configHome(), "gtk-3.0/settings.ini")
	if pathExists(configFile) {
		lines, err := loadTextFile(configFile)
		if err == nil {
			log.Infof("Loaded %s", configFile)
		} else {
			log.Warnf("Couldn't load %s", configFile)
		}

		for _, line := range lines {
			// In case users settings.ini had some lines we didn't expect,
			// we'll append them back from here.
			if !strings.HasPrefix(line, "[") {
				originalGtkConfig = append(originalGtkConfig, line)
			}

			if !strings.HasPrefix(line, "[") && !strings.HasPrefix(line, "#") &&
				strings.Contains(line, "=") {
				parts := strings.Split(line, "=")
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				switch key {
				case "gtk-theme-name":
					gtkConfig.themeName = value
				case "gtk-icon-theme-name":
					gtkConfig.iconThemeName = value
				case "gtk-font-name":
					gtkConfig.fontName = value
				case "gtk-cursor-theme-name":
					gtkConfig.cursorThemeName = value
				case "gtk-cursor-theme-size":
					i := intValue(value)
					if i != -1 {
						gtkConfig.cursorThemeSize = i
					} else {
						gtkConfig.cursorThemeSize = 0
					}
				case "gtk-toolbar-style":
					gtkConfig.toolbarStyle = value
				case "gtk-toolbar-icon-size":
					gtkConfig.toolbarIconSize = value
				case "gtk-button-images":
					gtkConfig.buttonImages = value == "1"
				case "gtk-menu-images":
					gtkConfig.menuImages = value == "1"
				case "gtk-enable-event-sounds":
					gtkConfig.enableEventSounds = value == "1"
				case "gtk-enable-input-feedback-sounds":
					gtkConfig.enableInputFeedbackSounds = value == "1"
				case "gtk-xft-antialias":
					gtkConfig.xftAntialias = intValue(value)
				case "gtk-xft-dpi":
					gtkConfig.xftDpi = intValue(value)
				case "gtk-xft-hinting":
					gtkConfig.xftHinting = intValue(value)
				case "gtk-xft-hintstyle":
					gtkConfig.xftHintstyle = value
				case "gtk-xft-rgba":
					gtkConfig.xftRgba = value
				default:
					log.Warnf("Unsupported config key: %s", key)
				}
			}
		}
	} else {
		log.Warnf("Could'n find %s", configFile)
	}
	log.Infof("gtk-theme-name:                   %s [default: Adwaita]", gtkConfig.themeName)
	log.Infof("gtk-icon-theme-name:              %s [default: Adwaita]", gtkConfig.iconThemeName)
	log.Infof("gtk-font-name:                    %s [default: Sans 10]", gtkConfig.fontName)
	log.Infof("gtk-cursor-theme-name:            %s [default: none]", gtkConfig.cursorThemeName)
	log.Infof("gtk-cursor-theme-size:            %v [default: 0]", gtkConfig.cursorThemeSize)
	log.Infof("gtk-toolbar-style:                %s [ignored]", gtkConfig.toolbarStyle)
	log.Infof("gtk-toolbar-icon-size:            %s [ignored]", gtkConfig.toolbarIconSize)
	log.Infof("gtk-button-images:                %v [default: false]", gtkConfig.buttonImages)
	log.Infof("gtk-menu-images:                  %v [default: false]", gtkConfig.menuImages)
	log.Infof("gtk-enable-event-sounds:          %v [default: true]", gtkConfig.enableEventSounds)
	log.Infof("gtk-enable-input-feedback-sounds: %v [default: true]", gtkConfig.enableInputFeedbackSounds)
	log.Infof("gtk-xft-antialias:                %v [0=no, 1=yes, -1=default]", gtkConfig.xftAntialias)
	log.Infof("gtk-xft-dpi:                      %v [1024*dots/inch. -1 for default]", gtkConfig.xftDpi)
	log.Infof("gtk-xft-hinting:                  %v [0=no, 1=yes, -1=default]", gtkConfig.xftHinting)
	log.Infof("gtk-xft-hintstyle:                %v [hintnone|hintslight|hintmedium|hintfull]", gtkConfig.xftHintstyle)
	log.Infof("gtk-xft-rgba:                     %v [none|rgb|bgr|vrgb|vbgr]", gtkConfig.xftRgba)

	// Apply setting to the window
	gtkSettings.SetProperty("gtk-theme-name", gtkConfig.themeName)
	gtkSettings.SetProperty("gtk-icon-theme-name", gtkConfig.iconThemeName)
	gtkSettings.SetProperty("gtk-font-name", gtkConfig.fontName)
	gtkSettings.SetProperty("gtk-cursor-theme-name", gtkConfig.cursorThemeName)
	// In docs 0 is default, but setting 0 prevents the cursor theme from loading!
	if gtkConfig.cursorThemeSize > 0 {
		gtkSettings.SetProperty("gtk-cursor-theme-size", gtkConfig.cursorThemeSize)
	}
	gtkSettings.SetProperty("gtk-button-images", gtkConfig.buttonImages)
	gtkSettings.SetProperty("gtk-menu-images", gtkConfig.menuImages)
	gtkSettings.SetProperty("gtk-enable-event-sounds", gtkConfig.enableEventSounds)
	gtkSettings.SetProperty("gtk-enable-input-feedback-sounds", gtkConfig.enableInputFeedbackSounds)
	gtkSettings.SetProperty("gtk-xft-antialias", gtkConfig.xftAntialias)
	gtkSettings.SetProperty("gtk-xft-dpi", gtkConfig.xftDpi)
	gtkSettings.SetProperty("gtk-xft-hinting", gtkConfig.xftHinting)
	gtkSettings.SetProperty("gtk-xft-hintstyle", gtkConfig.xftHintstyle)
	gtkSettings.SetProperty("gtk-xft-rgba", gtkConfig.xftRgba)
}

func intValue(s string) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	return -1
}

func applyGtkSettings() {
	gnomeSchema := "org.gnome.desktop.interface"

	cmd := exec.Command("gsettings", "set", gnomeSchema, "gtk-theme", gtkConfig.themeName)
	err := cmd.Run()
	if err != nil {
		log.Warn(err)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "icon-theme", gtkConfig.iconThemeName)
	err = cmd.Run()
	if err != nil {
		log.Warn(err)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "cursor-theme", gtkConfig.cursorThemeName)
	err = cmd.Run()
	if err != nil {
		log.Warn(err)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "font-name", gtkConfig.fontName)
	err = cmd.Run()
	if err != nil {
		log.Warn(err)
	}
}

func saveGtkSettings() {
	configFile := filepath.Join(configHome(), "gtk-3.0/settings.ini")
	lines := []string{"[Settings]"}

	lines = append(lines, fmt.Sprintf("gtk-theme-name=%s", gtkConfig.themeName))
	lines = append(lines, fmt.Sprintf("gtk-icon-theme-name=%s", gtkConfig.iconThemeName))
	lines = append(lines, fmt.Sprintf("gtk-font-name=%s", gtkConfig.fontName))
	lines = append(lines, fmt.Sprintf("gtk-cursor-theme-name=%s", gtkConfig.cursorThemeName))
	lines = append(lines, fmt.Sprintf("gtk-cursor-theme-size=%v", gtkConfig.cursorThemeSize))
	lines = append(lines, fmt.Sprintf("gtk-toolbar-style=%s", gtkConfig.toolbarStyle))
	lines = append(lines, fmt.Sprintf("gtk-toolbar-icon-size=%s", gtkConfig.toolbarIconSize))
	v := 0
	if gtkConfig.buttonImages {
		v = 1
	}
	lines = append(lines, fmt.Sprintf("gtk-button-images=%v", v))
	if gtkConfig.menuImages {
		v = 1
	} else {
		v = 0
	}
	lines = append(lines, fmt.Sprintf("gtk-menu-images=%v", v))
	if gtkConfig.enableEventSounds {
		v = 1
	} else {
		v = 0
	}
	lines = append(lines, fmt.Sprintf("gtk-enable-event-sounds=%v", v))
	if gtkConfig.enableInputFeedbackSounds {
		v = 1
	} else {
		v = 0
	}
	lines = append(lines, fmt.Sprintf("gtk-enable-input-feedback-sounds=%v", v))
	lines = append(lines, fmt.Sprintf("gtk-xft-antialias=%v", gtkConfig.xftAntialias))
	lines = append(lines, fmt.Sprintf("gtk-xft-dpi=%v", gtkConfig.xftDpi))
	lines = append(lines, fmt.Sprintf("gtk-xft-hinting=%v", gtkConfig.xftHinting))
	lines = append(lines, fmt.Sprintf("gtk-xft-hintstyle=%s", gtkConfig.xftHintstyle))
	lines = append(lines, fmt.Sprintf("gtk-xft-rgba=%s", gtkConfig.xftRgba))

	// append unsupported lines / comments from the original settings.ini file
	for _, l := range originalGtkConfig {
		if l != "" && !isDefined(l) {
			lines = append(lines, l)
		}
	}

	saveTextFile(lines, configFile)
}

func isDefined(line string) bool {
	defined := []string{
		"gtk-theme-name",
		"gtk-icon-theme-name",
		"gtk-font-name",
		"gtk-cursor-theme-name",
		"gtk-cursor-theme-size",
		"gtk-toolbar-style",
		"gtk-toolbar-icon-size",
		"gtk-button-images",
		"gtk-menu-images",
		"gtk-enable-event-sounds",
		"gtk-enable-input-feedback-sounds",
		"gtk-xft-antialias",
		"gtk-xft-dpi",
		"gtk-xft-hinting",
		"gtk-xft-hintstyle",
		"gtk-xft-rgba",
	}
	for _, d := range defined {
		if strings.HasPrefix(line, d) {
			return true
		}
	}
	return false
}

func getThemeNames() []string {
	var dirs []string

	// get theme dirs
	for _, dir := range dataDirs {
		if pathExists(filepath.Join(dir, "themes")) {
			dirs = append(dirs, filepath.Join(dir, "themes"))
		}
	}

	home := os.Getenv("HOME")
	if home != "" {
		if pathExists(filepath.Join(home, ".themes")) {
			dirs = append(dirs, filepath.Join(home, ".themes"))
		}
	}

	exclusions := []string{"Default", "Emacs"}
	var names []string
	for _, d := range dirs {
		files, err := listFiles(d)
		if err == nil {
			for _, f := range files {
				if f.IsDir() {
					subdirs, err := listFiles(filepath.Join(d, f.Name()))
					if err == nil {
						for _, sd := range subdirs {
							if sd.IsDir() && strings.HasPrefix(sd.Name(), "gtk-") {
								if !isIn(names, f.Name()) {
									if !isIn(exclusions, f.Name()) {
										names = append(names, f.Name())
										log.Debugf("Theme found: %s", f.Name())
									} else {
										log.Debugf("Excluded theme: %s", f.Name())
									}
									break
								}
							}
						}
					}
				}
			}
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	return names
}

// returns map[displayName]folderName
func getIconThemeNames() map[string]string {
	var dirs []string
	name2folderName := make(map[string]string)

	// get icon theme dirs
	for _, dir := range dataDirs {
		if pathExists(filepath.Join(dir, "icons")) {
			dirs = append(dirs, filepath.Join(dir, "icons"))
		}
	}

	home := os.Getenv("HOME")
	if home != "" {
		if pathExists(filepath.Join(home, ".icons")) {
			dirs = append(dirs, filepath.Join(home, ".icons"))
		}
	}

	exclusions := []string{"default", "hicolor", "locolor"}
	var names []string
	for _, d := range dirs {
		files, err := listFiles(d)
		if err == nil {
			for _, f := range files {
				if f.IsDir() {
					if !isIn(exclusions, f.Name()) {
						name, hasDirs, err := iconThemeName(filepath.Join(d, f.Name()))
						if err == nil && hasDirs {
							names = append(names, name)
							name2folderName[name] = f.Name()
							log.Debugf("Icon theme found: %s", name)
						}
					} else {
						log.Debugf("Excluded icon theme: %s", f.Name())
					}
				}
			}
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.ToUpper(names[i]) < strings.ToUpper(names[j])
	})

	return name2folderName
}

func getCursorThemes() (map[string]string, map[string]string) {
	var dirs []string
	name2path := make(map[string]string)
	name2FolderName := make(map[string]string)

	// get icon theme dirs
	for _, dir := range dataDirs {
		if pathExists(filepath.Join(dir, "icons")) {
			dirs = append(dirs, filepath.Join(dir, "icons"))
		}
	}

	home := os.Getenv("HOME")
	if home != "" {
		if pathExists(filepath.Join(home, ".icons")) {
			dirs = append(dirs, filepath.Join(home, ".icons"))
		}
	}

	exclusions := []string{"default", "hicolor", "locolor"}
	for _, d := range dirs {
		files, err := listFiles(d)
		if err == nil {
			for _, f := range files {
				if f.IsDir() {
					if !isIn(exclusions, f.Name()) {
						content, _ := listFiles(filepath.Join(d, f.Name()))
						if err == nil {
							for _, item := range content {
								if item.Name() == "cursors" {
									name, _, err := iconThemeName(filepath.Join(d, f.Name()))
									if err == nil {
										name2FolderName[name] = f.Name()
									}
									log.Debugf("Cursor theme found: %s", f.Name())
									name2path[f.Name()] = filepath.Join(d, f.Name(), "cursors")
								}
							}
						}
					}
				}
			}
		}
	}

	return name2path, name2FolderName
}

func getDataDirs() []string {
	var dirs []string
	xdgDataDirs := ""

	home := os.Getenv("HOME")
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		dirs = append(dirs, xdgDataHome)
	} else if home != "" {
		dirs = append(dirs, filepath.Join(home, ".local/share"))
	}

	if os.Getenv("XDG_DATA_DIRS") != "" {
		xdgDataDirs = os.Getenv("XDG_DATA_DIRS")
	} else {
		xdgDataDirs = "/usr/local/share/:/usr/share/"
	}

	for _, d := range strings.Split(xdgDataDirs, ":") {
		dirs = append(dirs, d)
	}

	var confirmedDirs []string
	for _, d := range dirs {
		if pathExists(d) {
			confirmedDirs = append(confirmedDirs, d)
		}
	}
	return confirmedDirs
}

func iconThemeName(path string) (string, bool, error) {
	name := ""
	hasDirs := false

	lines, err := loadTextFile(filepath.Join(path, "index.theme"))
	if err != nil {
		return name, hasDirs, err
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "Name=") || strings.HasPrefix(line, "Name =") {
			name = strings.Split(line, "=")[1]
			name = strings.TrimSpace(name)
			break
		}
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "Directories=") || strings.HasPrefix(line, "Directories =") {
			hasDirs = true
			break
		}
	}
	return name, hasDirs, err
}

func loadTextFile(path string) ([]string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	var output []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		output = append(output, line)
	}
	return output, nil
}

func saveTextFile(text []string, path string) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	datawriter := bufio.NewWriter(file)

	for _, data := range text {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}

func listFiles(dir string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err == nil {
		return files, nil
	}
	return nil, err
}

func isIn(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func pathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func tempDir() string {
	if os.Getenv("TMPDIR") != "" {
		return os.Getenv("TMPDIR")
	} else if os.Getenv("TEMP") != "" {
		return os.Getenv("TEMP")
	} else if os.Getenv("TMP") != "" {
		return os.Getenv("TMP")
	}
	return "/tmp"
}

func makeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err == nil {
			log.Debugf("Creating dir: %s", dir)
		}
	}
}

// Assert types to gtk.Builder objects
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

func getLabel(b *gtk.Builder, id string) (*gtk.Label, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}
	label, ok := obj.(*gtk.Label)
	if !ok {
		return nil, err
	}
	return label, nil
}

func getMenuBar(b *gtk.Builder, id string) (*gtk.MenuBar, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}
	menuBar, ok := obj.(*gtk.MenuBar)
	if !ok {
		return nil, err
	}
	return menuBar, nil
}

func getMenuItem(b *gtk.Builder, id string) (*gtk.MenuItem, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}
	item, ok := obj.(*gtk.MenuItem)
	if !ok {
		return nil, err
	}
	return item, nil
}
