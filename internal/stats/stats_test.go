package stats_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/differ"
	"github.com/yourusername/envlens/internal/stats"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "DB_PORT", Type: differ.Modified, OldValue: "5432", NewValue: "5433"},
		{Key: "AUTH_SECRET", Type: differ.Removed, OldValue: "old"},
		{Key: "AUTH_TOKEN", Type: differ.Added, NewValue: "tok"},
		{Key: "APP_ENV", Type: differ.Modified, OldValue: "dev", NewValue: "prod"},
	}
}

func TestCompute_Counts(t *testing.T) {
	r := stats.Compute(sampleChanges())
	if r.Total != 5 {
		t.Errorf("expected Total=5, got %d", r.Total)
	}
	if r.Added != 2 {
		t.Errorf("expected Added=2, got %d", r.Added)
	}
	if r.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", r.Removed)
	}
	if r.Modified != 2 {
		t.Errorf("expected Modified=2, got %d", r.Modified)
	}
}

func TestCompute_EmptyChanges(t *testing.T) {
	r := stats.Compute(nil)
	if r.Total != 0 || r.Added != 0 || r.Removed != 0 || r.Modified != 0 {
		t.Error("expected all zero counts for empty input")
	}
}

func TestCompute_TopPrefixes(t *testing.T) {
	r := stats.Compute(sampleChanges())
	if len(r.TopPrefixes) == 0 {
		t.Fatal("expected at least one top prefix")
	}
	// DB and AUTH each have 2 changes, APP has 1
	top := r.TopPrefixes[0]
	if top.Count < 2 {
		t.Errorf("expected top prefix count >= 2, got %d", top.Count)
	}
}

func TestCompute_TopPrefixes_LimitedToFive(t *testing.T) {
	var changes []differ.Change
	prefixes := []string{"AA", "BB", "CC", "DD", "EE", "FF", "GG"}
	for _, p := range prefixes {
		changes = append(changes, differ.Change{Key: p + "_KEY", Type: differ.Added})
	}
	r := stats.Compute(changes)
	if len(r.TopPrefixes) > 5 {
		t.Errorf("expected at most 5 top prefixes, got %d", len(r.TopPrefixes))
	}
}

func TestSummary_ContainsAllFields(t *testing.T) {
	r := stats.Compute(sampleChanges())
	s := stats.Summary(r)
	for _, want := range []string{"total=", "added=", "removed=", "modified="} {
		if !strings.Contains(s, want) {
			t.Errorf("summary missing field %q: %s", want, s)
		}
	}
}

func TestSummary_EmptyChanges(t *testing.T) {
	r := stats.Compute(nil)
	s := stats.Summary(r)
	if !strings.Contains(s, "total=0") {
		t.Errorf("expected total=0 in summary, got: %s", s)
	}
}
