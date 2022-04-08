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

func loadGtkConfig() {
	// parse gtk settings file
	originalGtkConfig = []string{}
	configFile := filepath.Join(configHome(), "gtk-3.0/settings.ini")
	if pathExists(configFile) {
		lines, err := loadTextFile(configFile)
		if err == nil {
			log.Infof(">>> Parsing original %s", configFile)
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
	log.Debugf("gtk-theme-name: %s", gtkConfig.themeName)
	log.Debugf("gtk-icon-theme-name: %s", gtkConfig.iconThemeName)
	log.Debugf("gtk-font-name: %s", gtkConfig.fontName)
	log.Debugf("gtk-cursor-theme-name: %s", gtkConfig.cursorThemeName)
	log.Debugf("gtk-cursor-theme-size: %v", gtkConfig.cursorThemeSize)
	log.Debugf("gtk-toolbar-style: %s", gtkConfig.toolbarStyle)
	log.Debugf("gtk-toolbar-icon-size: %s", gtkConfig.toolbarIconSize)
	log.Debugf("gtk-button-images: %v", gtkConfig.buttonImages)
	log.Debugf("gtk-menu-images: %v", gtkConfig.menuImages)
	log.Debugf("gtk-enable-event-sounds: %v", gtkConfig.enableEventSounds)
	log.Debugf("gtk-enable-input-feedback-sounds: %v", gtkConfig.enableInputFeedbackSounds)
	log.Debugf("gtk-xft-antialias: %v", gtkConfig.xftAntialias)
	log.Debugf("gtk-xft-hinting: %v", gtkConfig.xftHinting)
	log.Debugf("gtk-xft-hintstyle: %v", gtkConfig.xftHintstyle)
	log.Debugf("gtk-xft-rgba: %v", gtkConfig.xftRgba)
}

func intValue(s string) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	// -1 is default
	return -1
}

func readGsettings() {
	log.Info(">>> Reading gsettings")

	val, err := getGsettingsValue("org.gnome.desktop.interface", "gtk-theme")
	if err == nil {
		gsettings.gtkTheme = val
		log.Infof("gtk-theme: %s", gsettings.gtkTheme)
	} else {
		log.Warnf("Couldn't read gtk-theme, leaving default %s",
			gsettings.gtkTheme)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "icon-theme")
	if err == nil {
		gsettings.iconTheme = val
		log.Infof("icon-theme: %s", gsettings.iconTheme)
	} else {
		log.Warnf("Couldn't read icon-theme, leaving default %s",
			gsettings.iconTheme)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "font-name")
	if err == nil {
		gsettings.fontName = val
		log.Infof("font-name: %s", gsettings.fontName)
	} else {
		log.Warnf("Couldn't read font-name, leaving default %s",
			gsettings.fontName)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "cursor-theme")
	if err == nil {
		gsettings.cursorTheme = val
		log.Infof("cursor-theme: %s", gsettings.cursorTheme)
	} else {
		gsettings.cursorTheme = ""
		log.Warnf("Couldn't read cursor-theme, leaving default %s",
			gsettings.cursorTheme)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "cursor-size")
	if err == nil {
		v, e := strconv.Atoi(val)
		if e == nil {
			gsettings.cursorSize = v
			log.Infof("cursor-size: %v", gsettings.cursorSize)
		}
	} else {
		log.Warnf("Couldn't read cursorSize, leaving default %s",
			gsettings.cursorSize)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "toolbar-style")
	if err == nil {
		gsettings.toolbarStyle = val
		log.Infof("toolbar-style: %s", gsettings.toolbarStyle)
	} else {
		log.Warnf("Couldn't read toolbar-style, leaving default %s",
			gsettings.toolbarStyle)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "toolbar-icons-size")
	if err == nil {
		gsettings.toolbarIconsSize = val
		log.Infof("toolbar-icons-size: %s", gsettings.toolbarIconsSize)
	} else {
		log.Warnf("Couldn't read toolbar-icons-size, leaving default %s",
			gsettings.toolbarIconsSize)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "font-hinting")
	if err == nil {
		gsettings.fontHinting = val
		log.Infof("font-hinting: %s", gsettings.fontHinting)
	} else {
		log.Warnf("Couldn't read font-hinting, leaving default %s",
			gsettings.fontHinting)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "font-antialiasing")
	if err == nil {
		gsettings.fontAntialiasing = val
		log.Infof("font-antialiasing: %s", gsettings.fontAntialiasing)
	} else {
		log.Warnf("Couldn't read font-antialiasing, leaving default %s",
			gsettings.fontAntialiasing)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "font-rgba-order")
	if err == nil {
		gsettings.fontRgbaOrder = val
		log.Infof("font-rgba-order: %s", gsettings.fontRgbaOrder)
	} else {
		log.Warnf("Couldn't read font-rgba-order, leaving default %s",
			gsettings.fontRgbaOrder)
	}

	val, err = getGsettingsValue("org.gnome.desktop.interface", "text-scaling-factor")
	if err == nil {
		v, e := strconv.ParseFloat(val, 32)
		if e == nil {
			gsettings.textScalingFactor = v
			log.Infof("text-scaling-factor: %v", gsettings.textScalingFactor)
		}
	} else {
		log.Warnf("Couldn't read textScalingFactor, leaving default %s",
			gsettings.textScalingFactor)
	}

	val, err = getGsettingsValue("org.gnome.desktop.sound", "event-sounds")
	if err == nil {
		if val == "true" {
			gsettings.eventSounds = true
		} else {
			gsettings.eventSounds = false
		}
		log.Infof("event-sounds: %v", gsettings.eventSounds)
	} else {
		log.Warnf("Couldn't read event-sounds, leaving default %v",
			gsettings.eventSounds)
	}

	val, err = getGsettingsValue("org.gnome.desktop.sound", "input-feedback-sounds")
	if err == nil {
		if val == "true" {
			gsettings.inputFeedbackSounds = true
		} else {
			gsettings.inputFeedbackSounds = false
		}
		log.Infof("input-feedback-sounds: %v", gsettings.inputFeedbackSounds)
	} else {
		log.Warnf("Couldn't read input-feedback-sounds, leaving default %v",
			gsettings.inputFeedbackSounds)
	}
}

func saveGsettings() {
	gsettingsFile := filepath.Join(dataHome(), "nwg-look/")
	makeDir(gsettingsFile)
	log.Infof(">>> Backing up gsettings to %s", gsettingsFile)

	lines := []string{"# Generated by nwg-look, do not edit this file."}

	for _, key := range []string{
		"gtk-theme",
		"icon-theme",
		"font-name",
		"cursor-theme",
		"cursor-size",
		"toolbar-style",
		"toolbar-icons-size",
		"font-hinting",
		"font-antialiasing",
		"font-rgba-order",
		"text-scaling-factor"} {
		val, err := getGsettingsValue("org.gnome.desktop.interface", key)
		if err == nil {
			line := fmt.Sprintf("%s=%s", key, val)
			lines = append(lines, line)
		} else {
			log.Warnf("Couldn't get gsettings key: $s", key)
		}
	}
	for _, key := range []string{"event-sounds", "input-feedback-sounds"} {
		val, err := getGsettingsValue("org.gnome.desktop.sound", key)
		if err == nil {
			line := fmt.Sprintf("%s=%s", key, val)
			lines = append(lines, line)
		} else {
			log.Warnf("Couldn't get gsettings key: $s", key)
		}
	}

	saveTextFile(lines, filepath.Join(dataHome(), "nwg-look/gsettings"))
}

func getGsettingsValue(schema, key string) (string, error) {
	cmd := exec.Command("gsettings", "get", schema, key)
	out, err := cmd.CombinedOutput()
	if err == nil {
		s := fmt.Sprintf("%s", strings.TrimSpace(string(out)))
		if strings.HasPrefix(s, "'") {
			return s[1 : len(s)-1], nil
		}
		return s, nil
	}
	return "", err
}

func applyGsettings() {
	gnomeSchema := "org.gnome.desktop.interface"
	log.Info(">>> Applying gsettings")
	log.Infof(">> %s", gnomeSchema)

	cmd := exec.Command("gsettings", "set", gnomeSchema, "gtk-theme", gsettings.gtkTheme)
	err := cmd.Run()
	if err != nil {
		log.Warnf("gtk-theme: %s", err)
	} else {
		log.Infof("gtk-theme: %s OK", gsettings.gtkTheme)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "icon-theme", gsettings.iconTheme)
	err = cmd.Run()
	if err != nil {
		log.Warnf("icon-theme: %s", err)
	} else {
		log.Infof("icon-theme: %s OK", gsettings.iconTheme)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "cursor-theme", gsettings.cursorTheme)
	err = cmd.Run()
	if err != nil {
		log.Warnf("cursor-theme: %s", err)
	} else {
		log.Infof("cursor-theme: %s OK", gsettings.cursorTheme)
	}

	var val string

	val = strconv.Itoa(gsettings.cursorSize)
	cmd = exec.Command("gsettings", "set", gnomeSchema, "cursor-size", val)
	err = cmd.Run()
	if err != nil {
		log.Warnf("cursor-size: %s", err)
	} else {
		log.Infof("cursor-size: %s OK", val)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "font-name", gsettings.fontName)
	err = cmd.Run()
	if err != nil {
		log.Warnf("font-name: %s %s", gsettings.fontName, err)
	} else {
		log.Infof("font-name: %s OK", gsettings.fontName)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "font-hinting", gsettings.fontHinting)
	err = cmd.Run()
	if err != nil {
		log.Warnf("font-hinting: %s %s", gsettings.fontHinting, err)
	} else {
		log.Infof("font-hinting: %s OK", gsettings.fontHinting)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "font-antialiasing", gsettings.fontAntialiasing)
	err = cmd.Run()
	if err != nil {
		log.Warnf("font-antialiasing: %s %s", gsettings.fontAntialiasing, err)
	} else {
		log.Infof("font-antialiasing: %s OK", gsettings.fontAntialiasing)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "font-rgba-order", gsettings.fontRgbaOrder)
	err = cmd.Run()
	if err != nil {
		log.Warnf("font-rgba-order: %s %s", gsettings.fontRgbaOrder, err)
	} else {
		log.Infof("font-rgba-order: %s OK", gsettings.fontRgbaOrder)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "text-scaling-factor", fmt.Sprintf("%f", gsettings.textScalingFactor))
	err = cmd.Run()
	if err != nil {
		log.Warnf("text-scaling-factor: %s %s", gsettings.textScalingFactor, err)
	} else {
		log.Infof("text-scaling-factor: %v OK", gsettings.textScalingFactor)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "toolbar-style", gsettings.toolbarStyle)
	err = cmd.Run()
	if err != nil {
		log.Warnf("toolbar-style: %s %s", gsettings.toolbarStyle, err)
	} else {
		log.Infof("toolbar-style: %s OK", gsettings.toolbarStyle)
	}

	cmd = exec.Command("gsettings", "set", gnomeSchema, "toolbar-icons-size", gsettings.toolbarIconsSize)
	err = cmd.Run()
	if err != nil {
		log.Warnf("toolbar-icons-size: %s %s", gsettings.toolbarIconsSize, err)
	} else {
		log.Infof("toolbar-icons-size: %s OK", gsettings.toolbarIconsSize)
	}

	gnomeSchema = "org.gnome.desktop.sound"
	log.Infof(">> %s", gnomeSchema)

	if gsettings.eventSounds {
		val = "true"
	} else {
		val = "false"
	}
	cmd = exec.Command("gsettings", "set", gnomeSchema, "event-sounds", val)
	err = cmd.Run()
	if err != nil {
		log.Warnf("event-sounds: %s %s", val, err)
	} else {
		log.Infof("event-sounds: %s OK", val)
	}

	if gsettings.inputFeedbackSounds {
		val = "true"
	} else {
		val = "false"
	}
	cmd = exec.Command("gsettings", "set", gnomeSchema, "input-feedback-sounds", val)
	err = cmd.Run()
	if err != nil {
		log.Warnf("input-feedback-sounds: %s %s", val, err)
	} else {
		log.Infof("input-feedback-sounds: %s OK", val)
	}
}

func applyGsettingsFromFile() {
	gsettingsFile := filepath.Join(dataHome(), "nwg-look/gsettings")
	if pathExists(gsettingsFile) {
		log.Infof("Loading gsettings from %s", gsettingsFile)
		lines, err := loadTextFile(gsettingsFile)
		if err != nil {
			log.Fatalf("Failed loading file: %s", err)
		}
		var key, value string
		for _, line := range lines {
			if !strings.HasPrefix(line, "#") {
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					key = parts[0]
					value = parts[1]

					switch key {
					case "gtk-theme":
						gsettings.gtkTheme = value
					case "icon-theme":
						gsettings.iconTheme = value
					case "font-name":
						gsettings.fontName = value
					case "cursor-theme":
						gsettings.cursorTheme = value
					case "cursor-size":
						v, err := strconv.Atoi(value)
						if err == nil {
							gsettings.cursorSize = v
						}
					case "toolbar-style":
						gsettings.toolbarStyle = value
					case "toolbar-icons-size":
						gsettings.toolbarIconsSize = value
					case "font-hinting":
						gsettings.fontHinting = value
					case "font-antialiasing":
						gsettings.fontAntialiasing = value
					case "font-rgba-order":
						gsettings.fontRgbaOrder = value
					case "event-sounds":
						gsettings.eventSounds = value == "true"
					case "input-feedback-sounds":
						gsettings.inputFeedbackSounds = value == "true"
					}
				}
			}
		}
		applyGsettings()
	} else {
		log.Warnf("Couldn't find file: %s", gsettingsFile)
		os.Exit(1)
	}
}

func saveGtkIni() {
	configFile := filepath.Join(configHome(), "gtk-3.0/settings.ini")
	log.Info(">>> Exporting settings.ini")

	lines := []string{"[Settings]"}

	lines = append(lines, fmt.Sprintf("gtk-theme-name=%s", gsettings.gtkTheme))
	lines = append(lines, fmt.Sprintf("gtk-icon-theme-name=%s", gsettings.iconTheme))
	lines = append(lines, fmt.Sprintf("gtk-font-name=%s", gsettings.fontName))
	lines = append(lines, fmt.Sprintf("gtk-cursor-theme-name=%s", gsettings.cursorTheme))
	lines = append(lines, fmt.Sprintf("gtk-cursor-theme-size=%v", gsettings.cursorSize))

	// Ignored
	lines = append(lines, fmt.Sprintf("gtk-toolbar-style=%s", gtkConfig.toolbarStyle))
	lines = append(lines, fmt.Sprintf("gtk-toolbar-icon-size=%s", gtkConfig.toolbarIconSize))

	// Deprecated
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

	if gsettings.eventSounds {
		v = 1
	} else {
		v = 0
	}
	lines = append(lines, fmt.Sprintf("gtk-enable-event-sounds=%v", v))

	if gsettings.inputFeedbackSounds {
		v = 1
	} else {
		v = 0
	}
	lines = append(lines, fmt.Sprintf("gtk-enable-input-feedback-sounds=%v", v))

	if gsettings.fontAntialiasing != "none" {
		v = 1
	} else {
		v = 0
	}
	lines = append(lines, fmt.Sprintf("gtk-xft-antialias=%v", v))

	if gsettings.fontHinting != "none" {
		v = 1
	} else {
		v = 0
	}
	lines = append(lines, fmt.Sprintf("gtk-xft-hinting=%v", v))

	var fh string
	switch gsettings.fontHinting {
	case "slight":
		fh = "hintslight"
	case "medium":
		fh = "hintmedium"
	case "full":
		fh = "hintfull"
	default:
		fh = "hintnone"
	}
	lines = append(lines, fmt.Sprintf("gtk-xft-hintstyle=%s", fh))

	lines = append(lines, fmt.Sprintf("gtk-xft-rgba=%s", gsettings.fontRgbaOrder))

	// append unsupported lines / comments from the original settings.ini file
	for _, l := range originalGtkConfig {
		if l != "" && !isSupported(l) {
			lines = append(lines, l)
		}
	}

	for _, l := range lines {
		log.Debug(l)
	}
	log.Infof(">> Saving to %s", configFile)
	saveTextFile(lines, configFile)
}

func isSupported(line string) bool {
	supported := []string{
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
		"gtk-xft-hinting",
		"gtk-xft-hintstyle",
		"gtk-xft-rgba",
	}
	for _, d := range supported {
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

func dataHome() string {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return xdgDataHome
	}
	return filepath.Join(os.Getenv("HOME"), ".local/share")
}

func getDataDirs() []string {
	var dirs []string
	xdgDataDirs := ""

	dirs = append(dirs, dataHome())

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
