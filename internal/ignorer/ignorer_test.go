package ignorer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlens/internal/ignorer"
)

func writeTempIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envignore")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp ignore file: %v", err)
	}
	return path
}

func TestApply_IgnoresExactKeys(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	opts := ignorer.Options{Keys: []string{"FOO", "BAZ"}}
	result := ignorer.Apply(env, opts)
	if _, ok := result["FOO"]; ok {
		t.Error("expected FOO to be ignored")
	}
	if _, ok := result["BAZ"]; ok {
		t.Error("expected BAZ to be ignored")
	}
	if result["BAR"] != "2" {
		t.Error("expected BAR to be retained")
	}
}

func TestApply_IgnoresByPrefix(t *testing.T) {
	env := map[string]string{"CI_TOKEN": "x", "CI_URL": "y", "APP_NAME": "z"}
	opts := ignorer.Options{Prefixes: []string{"CI_"}}
	result := ignorer.Apply(env, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 key, got %d", len(result))
	}
	if result["APP_NAME"] != "z" {
		t.Error("expected APP_NAME to be retained")
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	result := ignorer.Apply(env, ignorer.Options{})
	if len(result) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(result))
	}
}

func TestFromFile_ParsesKeysAndPrefixes(t *testing.T) {
	content := "# comment\n\nDEBUG_MODE\nTEST_*\nSECRET_KEY\n"
	path := writeTempIgnoreFile(t, content)

	opts, err := ignorer.FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(opts.Keys))
	}
	if len(opts.Prefixes) != 1 || opts.Prefixes[0] != "TEST_" {
		t.Errorf("expected prefix TEST_, got %v", opts.Prefixes)
	}
}

func TestFromFile_MissingFile_ReturnsError(t *testing.T) {
	_, err := ignorer.FromFile("/nonexistent/.envignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
