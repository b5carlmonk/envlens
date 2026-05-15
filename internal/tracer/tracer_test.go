package tracer

import (
	"testing"

	"github.com/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "API_KEY", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "SECRET", Type: differ.Removed, OldValue: "s3cr3t"},
	}
}

func TestNew_EmptyTrace(t *testing.T) {
	tr := New()
	if tr.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", tr.Len())
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	tr := New()
	tr.Add("deploy-1", "staging.env", "prod.env", sampleChanges())
	if tr.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", tr.Len())
	}
}

func TestAdd_StoresLabel(t *testing.T) {
	tr := New()
	tr.Add("release-v2", "a.env", "b.env", sampleChanges())
	if tr.Entries()[0].Label != "release-v2" {
		t.Errorf("expected label 'release-v2', got %s", tr.Entries()[0].Label)
	}
}

func TestAdd_StoresSourceAndTarget(t *testing.T) {
	tr := New()
	tr.Add("test", "src.env", "tgt.env", sampleChanges())
	e := tr.Entries()[0]
	if e.Source != "src.env" || e.Target != "tgt.env" {
		t.Errorf("unexpected source/target: %s / %s", e.Source, e.Target)
	}
}

func TestFilterByLabel_ReturnsMatching(t *testing.T) {
	tr := New()
	tr.Add("alpha", "a.env", "b.env", sampleChanges())
	tr.Add("beta", "c.env", "d.env", sampleChanges())
	tr.Add("alpha", "e.env", "f.env", sampleChanges())

	result := tr.FilterByLabel("alpha")
	if len(result) != 2 {
		t.Errorf("expected 2 entries for label 'alpha', got %d", len(result))
	}
}

func TestFilterByLabel_NoMatch_ReturnsEmpty(t *testing.T) {
	tr := New()
	tr.Add("alpha", "a.env", "b.env", sampleChanges())

	result := tr.FilterByLabel("nonexistent")
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestSummary_EmptyTrace(t *testing.T) {
	tr := New()
	s := tr.Summary()
	if s != "trace: no entries recorded" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithEntries(t *testing.T) {
	tr := New()
	tr.Add("deploy", "a.env", "b.env", sampleChanges())
	s := tr.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
