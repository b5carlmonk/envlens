package patcher_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envlens/internal/differ"
	"github.com/envlens/internal/patcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestApply_AddsNewKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	changes := []differ.Change{{Type: differ.Added, Key: "BAZ", NewValue: "qux"}}
	result, err := patcher.Apply(path, changes, patcher.Options{Strategy: patcher.StrategyOverwrite})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 1 || result.Applied[0] != "BAZ" {
		t.Errorf("expected BAZ in applied, got %v", result.Applied)
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "BAZ=qux") {
		t.Errorf("expected BAZ=qux in file, got: %s", raw)
	}
}

func TestApply_RemovesKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nDEL=gone\n")
	changes := []differ.Change{{Type: differ.Removed, Key: "DEL", OldValue: "gone"}}
	_, err := patcher.Apply(path, changes, patcher.Options{Strategy: patcher.StrategyOverwrite})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw, _ := os.ReadFile(path)
	if strings.Contains(string(raw), "DEL=") {
		t.Errorf("expected DEL to be removed, got: %s", raw)
	}
}

func TestApply_SkipStrategy_KeepsExisting(t *testing.T) {
	path := writeTempEnv(t, "FOO=original\n")
	changes := []differ.Change{{Type: differ.Modified, Key: "FOO", OldValue: "original", NewValue: "updated"}}
	result, err := patcher.Apply(path, changes, patcher.Options{Strategy: patcher.StrategySkip})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "FOO=original") {
		t.Errorf("expected original value preserved, got: %s", raw)
	}
}

func TestApply_ErrorStrategy_ReturnsError(t *testing.T) {
	path := writeTempEnv(t, "FOO=original\n")
	changes := []differ.Change{{Type: differ.Modified, Key: "FOO", OldValue: "original", NewValue: "updated"}}
	_, err := patcher.Apply(path, changes, patcher.Options{Strategy: patcher.StrategyError})
	if err == nil {
		t.Error("expected error for conflict, got nil")
	}
}

func TestApply_DryRun_DoesNotWriteFile(t *testing.T) {
	original := "FOO=bar\n"
	path := writeTempEnv(t, original)
	changes := []differ.Change{{Type: differ.Added, Key: "NEW", NewValue: "val"}}
	result, err := patcher.Apply(path, changes, patcher.Options{Strategy: patcher.StrategyOverwrite, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected DryRun flag set in result")
	}
	raw, _ := os.ReadFile(path)
	if string(raw) != original {
		t.Errorf("expected file unchanged in dry-run, got: %s", raw)
	}
}
