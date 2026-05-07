package summarizer_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/summarizer"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "API_KEY", Type: differ.Removed, OldValue: "secret"},
		{Key: "PORT", Type: differ.Modified, OldValue: "8080", NewValue: "9090"},
		{Key: "APP_NAME", Type: differ.Unchanged, OldValue: "myapp", NewValue: "myapp"},
	}
}

func TestSummarize_Counts(t *testing.T) {
	s := summarizer.Summarize(sampleChanges())

	if s.Added != 1 {
		t.Errorf("expected 1 added, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", s.Modified)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", s.Unchanged)
	}
	if s.Total != 4 {
		t.Errorf("expected total 4, got %d", s.Total)
	}
}

func TestSummarize_Narrative_ContainsCounts(t *testing.T) {
	s := summarizer.Summarize(sampleChanges())

	for _, want := range []string{"1 added", "1 removed", "1 modified", "1 unchanged"} {
		if !strings.Contains(s.Narrative, want) {
			t.Errorf("narrative missing %q: %s", want, s.Narrative)
		}
	}
}

func TestSummarize_EmptyChanges(t *testing.T) {
	s := summarizer.Summarize([]differ.Change{})

	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
	if !strings.Contains(s.Narrative, "No environment variables") {
		t.Errorf("unexpected narrative for empty input: %s", s.Narrative)
	}
}

func TestSummarize_OnlyAdded(t *testing.T) {
	changes := []differ.Change{
		{Key: "NEW_KEY", Type: differ.Added, NewValue: "value"},
	}
	s := summarizer.Summarize(changes)

	if s.Added != 1 || s.Removed != 0 || s.Modified != 0 {
		t.Errorf("unexpected counts: %+v", s)
	}
	if !strings.Contains(s.Narrative, "1 added") {
		t.Errorf("narrative missing '1 added': %s", s.Narrative)
	}
}

func TestSummarize_NarrativeContainsTotal(t *testing.T) {
	s := summarizer.Summarize(sampleChanges())

	if !strings.Contains(s.Narrative, "4") {
		t.Errorf("narrative should mention total count: %s", s.Narrative)
	}
}
