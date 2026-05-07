// This file is derived from stylepak (https://github.com/refi64/stylepak)
// Original work licensed under the Mozilla Public License 2.0.
// Modifications and Go translation by Eslam Allam eslamallam73@gmail.com.

package stylepak

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)


const (
	nwgLookMarker = ".nwg-look"
)

type themeInstallAction int

const (
	themeActionFreshInstall themeInstallAction = iota
	themeActionUpdate
	themeActionSkip // user-managed, don't touch
)

// Runner abstracts command execution (useful for testing/mocking)
type Runner interface {
	Run(name string, args ...string) error
	Output(name string, args ...string) ([]byte, error)
}

type ExecRunner struct{}

func (r ExecRunner) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %s %v: %w", name, args, err)
	}
	return nil
}

func (r ExecRunner) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command output failed: %s %v: %w", name, args, err)
	}
	return out, nil
}

func determineInstallAction(themeDir string) (themeInstallAction, error) {
	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		return themeActionFreshInstall, nil
	} else if err != nil {
		return 0, fmt.Errorf("stat theme dir %q: %w", themeDir, err)
	}

	markerPath := filepath.Join(themeDir, nwgLookMarker)
	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		return themeActionSkip, nil
	} else if err != nil {
		return 0, fmt.Errorf("stat marker %q: %w", markerPath, err)
	}

	return themeActionUpdate, nil
}

func removeStaleThemes(themesDir, currentTheme string) error {
	entries, err := os.ReadDir(themesDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read themes dir %q: %w", themesDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == currentTheme {
			continue
		}
		markerPath := filepath.Join(themesDir, entry.Name(), nwgLookMarker)
		if _, err := os.Stat(markerPath); err != nil {
			continue
		}
		log.Infof("Removing previously managed theme: %s", entry.Name())
		if err := os.RemoveAll(filepath.Join(themesDir, entry.Name())); err != nil {
			return fmt.Errorf("remove previous theme %q: %w", entry.Name(), err)
		}
	}
	return nil
}

func copyThemeFiles(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read theme dir %q: %w", src, err)
	}

	copied := 0
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "gtk-") {
			continue
		}
		entryDst := filepath.Join(dst, entry.Name())
		if err := os.MkdirAll(entryDst, 0o755); err != nil {
			return fmt.Errorf("create dir %q: %w", entryDst, err)
		}
		if err := copyDir(filepath.Join(src, entry.Name()), entryDst); err != nil {
			return fmt.Errorf("copy theme dir %q: %w", entry.Name(), err)
		}
		copied++
		log.Info("Copied theme dir:", entry.Name())
	}

	if copied == 0 {
		return fmt.Errorf("no gtk-* directories found in theme %q", src)
	}

	if data, err := os.ReadFile(filepath.Join(src, "index.theme")); err == nil {
		if err := os.WriteFile(filepath.Join(dst, "index.theme"), data, 0o644); err != nil {
			return fmt.Errorf("write index.theme: %w", err)
		}
		log.Info("Copied index.theme")
	}

	return nil
}

func InstallUserTheme(theme string, runner Runner) error {
	if runner == nil {
		runner = ExecRunner{}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}

	dataHome := getenvDefault("XDG_DATA_HOME", filepath.Join(home, ".local", "share"))

	if theme == "" {
		theme, err = getCurrentGTKTheme(runner)
		if err != nil {
			return fmt.Errorf("resolve current GTK theme: %w", err)
		}
	}

	if err := validateTheme(theme); err != nil {
		return err
	}

	log.Info("Converting theme:", theme)

	themesDir := filepath.Join(home, ".themes")
	userThemeDir := filepath.Join(themesDir, theme)
	markerPath := filepath.Join(userThemeDir, nwgLookMarker)

	// Find system theme path, explicitly excluding ~/.themes to avoid
	// copying a theme into itself
	themePath, err := findThemePath(theme, dataHome, home)
	if err != nil {
		return fmt.Errorf("locate theme %q: %w", theme, err)
	}
	if themePath == userThemeDir {
		return fmt.Errorf("theme %q is already in ~/.themes and was not installed by nwg-look", theme)
	}
	log.Info("Found theme located at:", themePath)

	if err := removeStaleThemes(themesDir, theme); err != nil {
		return fmt.Errorf("remove stale themes: %w", err)
	}

	action, err := determineInstallAction(userThemeDir)
	if err != nil {
		return err
	}

	switch action {
	case themeActionFreshInstall:
		log.Infof("Installing theme to ~/.themes/%s", theme)
		if err := os.MkdirAll(userThemeDir, 0o755); err != nil {
			return fmt.Errorf("create ~/.themes/%s: %w", theme, err)
		}
		if err := copyThemeFiles(themePath, userThemeDir); err != nil {
			return err
		}
		if err := os.WriteFile(markerPath, []byte("managed by nwg-look\n"), 0o644); err != nil {
			return fmt.Errorf("write nwg-look marker: %w", err)
		}
	case themeActionUpdate:
		log.Infof("Updating existing nwg-look-managed theme at ~/.themes/%s", theme)
		if err := copyThemeFiles(themePath, userThemeDir); err != nil {
			return err
		}
	case themeActionSkip:
		log.Infof("~/.themes/%s already exists and was not created by nwg-look, skipping copy", theme)
	}

	if err := runner.Run("flatpak", "override", "--user",
		"--filesystem="+themesDir+":ro",
	); err != nil {
		return fmt.Errorf("flatpak override ~/.themes: %w", err)
	}
	log.Info("Configured flatpak to access ~/.themes")

	log.Infof("Successfully installed theme %s", theme)
	return nil
}
