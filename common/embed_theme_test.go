package common

import (
	"embed"
	"testing"
)

//go:embed embed_theme_testdata/themes
var themeTestFS embed.FS

const themeTestRoot = "embed_theme_testdata/themes"

func TestListEmbeddedThemes_FindsBoth(t *testing.T) {
	got := ListEmbeddedThemes(themeTestFS, themeTestRoot)
	want := []string{"alpha", "beta"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestValidateEmbeddedTheme_OK(t *testing.T) {
	if err := ValidateEmbeddedTheme(themeTestFS, themeTestRoot, "alpha"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateEmbeddedTheme_Missing(t *testing.T) {
	err := ValidateEmbeddedTheme(themeTestFS, themeTestRoot, "gamma")
	if err == nil {
		t.Fatal("expected error for missing theme, got nil")
	}
	msg := err.Error()
	for _, want := range []string{`"gamma"`, "alpha", "beta", "embed_theme_testdata/themes/gamma/index.html"} {
		if !contains(msg, want) {
			t.Errorf("error %q should contain %q", msg, want)
		}
	}
}

func TestValidateEmbeddedTheme_Empty(t *testing.T) {
	if err := ValidateEmbeddedTheme(themeTestFS, themeTestRoot, ""); err == nil {
		t.Fatal("expected error for empty theme, got nil")
	}
}

func TestValidateEmbeddedTheme_RootWithTrailingSlash(t *testing.T) {
	if err := ValidateEmbeddedTheme(themeTestFS, themeTestRoot+"/", "alpha"); err != nil {
		t.Fatalf("expected no error with trailing slash on root, got: %v", err)
	}
}

func contains(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
