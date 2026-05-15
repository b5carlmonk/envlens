package changelog_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/changelog"
	"github.com/user/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "API_KEY", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "DEPRECATED_VAR", Type: differ.Removed, OldValue: "gone"},
	}
}

func TestNew_EmptyChangelog(t *testing.T) {
	c := changelog.New()
	if c.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", c.Len())
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	c := changelog.New()
	c.Add("dev.env", "prod.env", sampleChanges())
	if c.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", c.Len())
	}
}

func TestAdd_MultipleEntries(t *testing.T) {
	c := changelog.New()
	c.Add("dev.env", "staging.env", sampleChanges())
	c.Add("staging.env", "prod.env", sampleChanges())
	if c.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", c.Len())
	}
}

func TestAdd_StoresSourceAndTarget(t *testing.T) {
	c := changelog.New()
	c.Add("dev.env", "prod.env", sampleChanges())
	e := c.Entries[0]
	if e.Source != "dev.env" {
		t.Errorf("expected source dev.env, got %s", e.Source)
	}
	if e.Target != "prod.env" {
		t.Errorf("expected target prod.env, got %s", e.Target)
	}
}

func TestRenderText_NoEntries(t *testing.T) {
	c := changelog.New()
	out := changelog.RenderText(c)
	if !strings.Contains(out, "No changelog entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestRenderText_ContainsSourceAndTarget(t *testing.T) {
	c := changelog.New()
	c.Add("dev.env", "prod.env", sampleChanges())
	out := changelog.RenderText(c)
	if !strings.Contains(out, "dev.env") {
		t.Errorf("expected source in output")
	}
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected target in output")
	}
}

func TestRenderText_ContainsChangeKeys(t *testing.T) {
	c := changelog.New()
	c.Add("a.env", "b.env", sampleChanges())
	out := changelog.RenderText(c)
	for _, key := range []string{"DB_HOST", "API_KEY", "DEPRECATED_VAR"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %s in output", key)
		}
	}
}

func TestRenderText_ContainsChangeTypes(t *testing.T) {
	c := changelog.New()
	c.Add("a.env", "b.env", sampleChanges())
	out := changelog.RenderText(c)
	for _, ct := range []string{"ADDED", "MODIFIED", "REMOVED"} {
		if !strings.Contains(out, ct) {
			t.Errorf("expected change type %s in output", ct)
		}
	}
}
