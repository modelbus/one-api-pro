package common

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

// ValidateEmbeddedTheme verifies that the given theme has its build output
// (<themesRoot>/<theme>/index.html) embedded in buildFS. If the theme is
// missing, the returned error lists every theme that IS embedded, so the
// operator can either switch to one of them or rebuild the binary with the
// expected theme.
//
// This guards against the failure mode where the database options table or
// environment variable references a theme name that is not embedded in the
// running binary, which previously caused the admin UI to render as a blank
// page (HTTP 200 with empty HTML) and made the service look unreachable.
//
// themesRoot is the embed-relative directory that contains one subdirectory
// per theme (each holding its own index.html).
func ValidateEmbeddedTheme(buildFS embed.FS, themesRoot, theme string) error {
	if strings.TrimSpace(theme) == "" {
		return fmt.Errorf("theme is empty")
	}
	themesRoot = strings.Trim(themesRoot, "/")
	if themesRoot == "" {
		return fmt.Errorf("themesRoot is empty")
	}

	indexPath := fmt.Sprintf("%s/%s/index.html", themesRoot, theme)
	if _, err := buildFS.ReadFile(indexPath); err == nil {
		return nil
	}

	available := ListEmbeddedThemes(buildFS, themesRoot)
	hint := "no themes are embedded"
	if len(available) > 0 {
		hint = "embedded themes: [" + strings.Join(available, ", ") + "]"
	}
	return fmt.Errorf(
		"theme %q is not embedded in this binary (missing %s); %s. "+
			"Fix: update the database option `theme` (or the THEME env var) to one of the embedded themes, "+
			"or rebuild the binary with the expected theme added to web/THEMES",
		theme, indexPath, hint,
	)
}

// ListEmbeddedThemes returns the sorted list of theme names whose build output
// is present under themesRoot in buildFS. A directory counts as a theme only
// if it contains an index.html at its root.
func ListEmbeddedThemes(buildFS embed.FS, themesRoot string) []string {
	root := strings.Trim(themesRoot, "/")
	if root == "" {
		return nil
	}
	entries, err := fs.ReadDir(buildFS, root)
	if err != nil {
		return nil
	}
	var themes []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		indexPath := root + "/" + name + "/index.html"
		if _, err := buildFS.ReadFile(indexPath); err == nil {
			themes = append(themes, name)
		}
	}
	sort.Strings(themes)
	return themes
}
