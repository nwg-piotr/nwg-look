// This file is derived from stylepak (https://github.com/refi64/stylepak)
// Original work licensed under the Mozilla Public License 2.0.
// Modifications and Go translation by Eslam Allam eslamallam73@gmail.com.

package stylepak

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const gtkThemeVer = "3.22"

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

// InstallUserTheme converts and installs a GTK theme for the current user.
func InstallUserTheme(theme string, runner Runner) error {
	if runner == nil {
		runner = ExecRunner{}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}

	cacheHome := getenvDefault("XDG_CACHE_HOME", filepath.Join(home, ".cache"))
	dataHome := getenvDefault("XDG_DATA_HOME", filepath.Join(home, ".local", "share"))
	stylepakCache := filepath.Join(cacheHome, "stylepak")

	// Resolve theme if not provided
	if theme == "" {
		theme, err = getCurrentGTKTheme(runner)
		if err != nil {
			return fmt.Errorf("resolve current GTK theme: %w", err)
		}
	}

	// Validate theme name
	if err := validateTheme(theme); err != nil {
		return err
	}

	appID := "org.gtk.Gtk3theme." + theme
	log.Info("Converting theme:", theme)

	// Locate theme path
	themePath, err := findThemePath(theme, dataHome, home)
	if err != nil {
		return fmt.Errorf("locate theme %q: %w", theme, err)
	}
	log.Info("Found theme located at:", themePath)

	rootDir := filepath.Join(stylepakCache, theme)
	repoDir := filepath.Join(rootDir, "repo")
	buildDir := filepath.Join(rootDir, "build")

	// Clean previous state
	if err := clean(rootDir, repoDir); err != nil {
		return fmt.Errorf("cleanup previous build dirs: %w", err)
	}

	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		return fmt.Errorf("create repo dir %q: %w", repoDir, err)
	}

	// Initialize OSTree repo
	if err := runner.Run("ostree", "--repo="+repoDir, "init", "--mode=archive"); err != nil {
		return fmt.Errorf("initialize ostree repo: %w", err)
	}
	if err := runner.Run("ostree", "--repo="+repoDir, "config", "set", "core.min-free-space-percent", "0"); err != nil {
		return fmt.Errorf("configure ostree repo: %w", err)
	}

	// Prepare build dir
	if err := clean(buildDir); err != nil {
		return fmt.Errorf("cleanup build dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(buildDir, "files"), 0o755); err != nil {
		return fmt.Errorf("create build files dir: %w", err)
	}

	// Detect GTK version
	gtkVer, err := detectGTKVersion(themePath)
	if err != nil {
		return fmt.Errorf("detect GTK version: %w", err)
	}

	src := filepath.Join(themePath, "gtk-3."+gtkVer)
	dst := filepath.Join(buildDir, "files")

	if err := copyDir(src, dst); err != nil {
		return fmt.Errorf("copy theme files from %q to %q: %w", src, dst, err)
	}

	// Write appdata
	appDataDir := filepath.Join(dst, "share", "appdata")
	if err := os.MkdirAll(appDataDir, 0o755); err != nil {
		return fmt.Errorf("create appdata dir: %w", err)
	}

	appData := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<component type="runtime">
  <id>%s</id>
  <metadata_license>CC0-1.0</metadata_license>
  <name>%s GTK Theme</name>
  <summary>%s (generated via stylepak)</summary>
</component>`, appID, theme, theme)

	appDataPath := filepath.Join(appDataDir, appID+".appdata.xml")
	if err := os.WriteFile(appDataPath, []byte(appData), 0o644); err != nil {
		return fmt.Errorf("write appdata file %q: %w", appDataPath, err)
	}

	// appstream-compose
	if err := runner.Run("appstream-compose",
		"--prefix="+dst,
		"--basename="+appID,
		"--origin=flatpak",
		appID,
	); err != nil {
		return fmt.Errorf("run appstream-compose: %w", err)
	}

	// Initial commit
	if err := runner.Run("ostree", "--repo="+repoDir, "commit", "-b", "base", "--tree=dir="+buildDir); err != nil {
		return fmt.Errorf("ostree base commit: %w", err)
	}

	arches, err := getFlatpakArchitectures(runner)
	if err != nil {
		return fmt.Errorf("get flatpak architectures: %w", err)
	}

	var bundles []string

	for _, arch := range arches {
		bundle := filepath.Join(rootDir, fmt.Sprintf("%s-%s.flatpak", appID, arch))

		if err := clean(buildDir); err != nil {
			return fmt.Errorf("cleanup build dir for arch %s: %w", arch, err)
		}

		if err := runner.Run("ostree", "--repo="+repoDir, "checkout", "-U", "base", buildDir); err != nil {
			return fmt.Errorf("ostree checkout for arch %s: %w", arch, err)
		}

		metadata := fmt.Sprintf(`[Runtime]
name=%s
runtime=%s/%s/%s
sdk=%s/%s/%s`, appID, appID, arch, gtkThemeVer, appID, arch, gtkThemeVer)

		metaPath := filepath.Join(buildDir, "metadata")
		if err := os.WriteFile(metaPath, []byte(metadata), 0o644); err != nil {
			return fmt.Errorf("write metadata for arch %s: %w", arch, err)
		}

		if err := runner.Run("ostree", "--repo="+repoDir, "commit",
			"-b", fmt.Sprintf("runtime/%s/%s/%s", appID, arch, gtkThemeVer),
			"--add-metadata-string", "xa.metadata="+metadata,
			"--link-checkout-speedup",
			buildDir,
		); err != nil {
			return fmt.Errorf("ostree commit for arch %s: %w", arch, err)
		}

		if err := runner.Run("flatpak", "build-bundle",
			"--runtime",
			"--arch="+arch,
			repoDir,
			bundle,
			appID,
			gtkThemeVer,
		); err != nil {
			return fmt.Errorf("build flatpak bundle for arch %s: %w", arch, err)
		}

		bundles = append(bundles, bundle)
	}

	for _, bundle := range bundles {
		if err := runner.Run("flatpak", "install", "-y", "--user", bundle); err != nil {
			return fmt.Errorf("install bundle %q: %w", bundle, err)
		}
	}

	return nil
}
