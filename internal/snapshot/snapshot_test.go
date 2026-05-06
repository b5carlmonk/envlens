package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envlens/internal/snapshot"
)

func TestNew_SetsFieldsCorrectly(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	before := time.Now().UTC()
	s := snapshot.New("staging", env)
	after := time.Now().UTC()

	if s.Label != "staging" {
		t.Errorf("expected label %q, got %q", "staging", s.Label)
	}
	if len(s.Env) != 2 {
		t.Errorf("expected 2 env entries, got %d", len(s.Env))
	}
	if s.Timestamp.Before(before) || s.Timestamp.After(after) {
		t.Errorf("timestamp %v is outside expected range", s.Timestamp)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	s := snapshot.New("production", env)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := snapshot.Save(s, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Label != s.Label {
		t.Errorf("expected label %q, got %q", s.Label, loaded.Label)
	}
	if loaded.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", loaded.Env["DB_HOST"])
	}
	if loaded.Env["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", loaded.Env["PORT"])
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSave_InvalidPath_ReturnsError(t *testing.T) {
	s := snapshot.New("test", map[string]string{})
	err := snapshot.Save(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(path, []byte("not valid json{"), 0644); err != nil {
		t.Fatalf("failed to write bad file: %v", err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
