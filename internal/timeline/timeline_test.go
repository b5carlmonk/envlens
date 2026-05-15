package timeline_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/timeline"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Added, Key: "DB_HOST", NewValue: "localhost"},
		{Type: differ.Removed, Key: "OLD_KEY", OldValue: "old"},
		{Type: differ.Modified, Key: "APP_ENV", OldValue: "staging", NewValue: "production"},
	}
}

func TestNew_EmptyTimeline(t *testing.T) {
	tl := timeline.New()
	if tl.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", tl.Len())
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	tl := timeline.New()
	tl.Add("deploy-1", "staging.env", "prod.env", sampleChanges())
	if tl.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", tl.Len())
	}
}

func TestAdd_StoresLabel(t *testing.T) {
	tl := timeline.New()
	tl.Add("release-v2", "a.env", "b.env", sampleChanges())
	if tl.Entries[0].Label != "release-v2" {
		t.Errorf("expected label 'release-v2', got %q", tl.Entries[0].Label)
	}
}

func TestAdd_StoresSourceAndTarget(t *testing.T) {
	tl := timeline.New()
	tl.Add("test", "source.env", "target.env", sampleChanges())
	e := tl.Entries[0]
	if e.Source != "source.env" || e.Target != "target.env" {
		t.Errorf("unexpected source/target: %q / %q", e.Source, e.Target)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	tl := timeline.New()
	tl.Add("deploy-1", "a.env", "b.env", sampleChanges())
	tl.Add("deploy-2", "b.env", "c.env", sampleChanges())

	dir := t.TempDir()
	path := filepath.Join(dir, "timeline.json")

	if err := tl.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := timeline.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", loaded.Len())
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0644)

	_, err := timeline.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := timeline.Load("/nonexistent/path/timeline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestKeyHistory_ReturnsMatchingEntries(t *testing.T) {
	tl := timeline.New()
	tl.Add("deploy-1", "a.env", "b.env", []differ.Change{
		{Type: differ.Added, Key: "DB_HOST", NewValue: "localhost"},
	})
	tl.Add("deploy-2", "b.env", "c.env", []differ.Change{
		{Type: differ.Modified, Key: "APP_ENV", OldValue: "staging", NewValue: "prod"},
	})
	tl.Add("deploy-3", "c.env", "d.env", []differ.Change{
		{Type: differ.Modified, Key: "DB_HOST", OldValue: "localhost", NewValue: "db.prod"},
	})

	history := tl.KeyHistory("DB_HOST")
	if len(history) != 2 {
		t.Errorf("expected 2 history entries for DB_HOST, got %d", len(history))
	}
}

func TestSave_ProducesValidJSON(t *testing.T) {
	tl := timeline.New()
	tl.Add("v1", "a.env", "b.env", sampleChanges())

	dir := t.TempDir()
	path := filepath.Join(dir, "tl.json")
	_ = tl.Save(path)

	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("expected valid JSON output: %v", err)
	}
}
