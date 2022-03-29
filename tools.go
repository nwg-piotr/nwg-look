// tools
package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

func loadGtkSettings() {
	settings, _ := gtk.SettingsGetDefault()
	prop, _ := settings.GetProperty("gtk-theme-name")
	gtkSettings.themeName, _ = prop.(string)
	log.Infof("Current theme: %s", gtkSettings.themeName)

	prop, _ = settings.GetProperty("gtk-icon-theme-name")
	gtkSettings.iconThemeName, _ = prop.(string)
	log.Infof("Icon theme: %s", gtkSettings.iconThemeName)

	prop, _ = settings.GetProperty("gtk-font-name")
	gtkSettings.fontName, _ = prop.(string)
	log.Infof("Default font: %s", gtkSettings.fontName)

	prop, _ = settings.GetProperty("gtk-cursor-theme-name")
	gtkSettings.cursorThemeName, _ = prop.(string)
	log.Infof("Cursor theme: %s", gtkSettings.cursorThemeName)
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
