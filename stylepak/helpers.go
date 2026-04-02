package stylepak

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func validateTheme(theme string) error {
	re := regexp.MustCompile(`^[A-Za-z0-9._\-]+$`)
	if !re.MatchString(theme) {
		return fmt.Errorf("invalid theme name %q (allowed: letters, digits, '.', '_', '-')", theme)
	}
	return nil
}

func clean(paths ...string) error {
	for _, p := range paths {
		if err := os.RemoveAll(p); err != nil {
			return fmt.Errorf("remove path %q: %w", p, err)
		}
	}
	return nil
}

func getCurrentGTKTheme(r Runner) (string, error) {
	out, err := r.Output("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	if err != nil {
		return "", fmt.Errorf("failed to get current GTK theme: %w", err)
	}
	return strings.Trim(string(out), "'\n "), nil
}

func findThemePath(theme, dataHome, home string) (string, error) {
	paths := []string{
		filepath.Join(dataHome, "themes"),
		filepath.Join(home, ".themes"),
		"/usr/share/themes",
	}

	for _, base := range paths {
		full := filepath.Join(base, theme)
		if fi, err := os.Stat(full); err == nil && fi.IsDir() {
			return full, nil
		}
	}
	return "", errors.New("theme directory not found in known locations")
}

func detectGTKVersion(themePath string) (string, error) {
	entries, err := os.ReadDir(themePath)
	if err != nil {
		return "", fmt.Errorf("read theme dir %q: %w", themePath, err)
	}

	re := regexp.MustCompile(`gtk-3\.(\d+)$`)
	var versions []int

	for _, e := range entries {
		if m := re.FindStringSubmatch(e.Name()); len(m) == 2 {
			var v int
			fmt.Sscanf(m[1], "%d", &v)
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return "", errors.New("no GTK 3.x directories found")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	return fmt.Sprintf("%d", versions[0]), nil
}

func getFlatpakArchitectures(r Runner) ([]string, error) {
	out, err := r.Output("flatpak", "list", "--runtime", "--columns=arch")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	set := make(map[string]struct{})

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			set[l] = struct{}{}
		}
	}

	if len(set) == 0 {
		return nil, errors.New("no flatpak architectures detected")
	}

	var arches []string
	for k := range set {
		arches = append(arches, k)
	}
	sort.Strings(arches)

	return arches, nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk path %q: %w", path, err)
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("compute relative path: %w", err)
		}

		target := filepath.Join(dst, rel)

		switch {
		case info.Mode()&os.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("read symlink %q: %w", path, err)
			}
			if err := os.Symlink(link, target); err != nil {
				return fmt.Errorf("create symlink %q: %w", target, err)
			}

		case info.IsDir():
			if err := os.MkdirAll(target, info.Mode()); err != nil {
				return fmt.Errorf("create dir %q: %w", target, err)
			}

		default:
			if err := copyFile(path, target, info.Mode()); err != nil {
				return err
			}
		}

		return nil
	})
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file %q: %w", src, err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("open destination file %q: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy file %q -> %q: %w", src, dst, err)
	}

	return nil
}
