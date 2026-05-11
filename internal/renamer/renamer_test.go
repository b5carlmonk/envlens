package renamer_test

import (
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/renamer"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Removed, Key: "DB_HOST", OldValue: "localhost"},
		{Type: differ.Added, Key: "DATABASE_HOST", NewValue: "localhost"},
		{Type: differ.Modified, Key: "APP_ENV", OldValue: "staging", NewValue: "production"},
	}
}

func TestDetect_IdentifiesRename(t *testing.T) {
	result := renamer.Detect(sampleChanges())

	if len(result.Hints) != 1 {
		t.Fatalf("expected 1 rename hint, got %d", len(result.Hints))
	}
	h := result.Hints[0]
	if h.OldKey != "DB_HOST" {
		t.Errorf("expected OldKey DB_HOST, got %s", h.OldKey)
	}
	if h.NewKey != "DATABASE_HOST" {
		t.Errorf("expected NewKey DATABASE_HOST, got %s", h.NewKey)
	}
	if h.Value != "localhost" {
		t.Errorf("expected Value localhost, got %s", h.Value)
	}
}

func TestDetect_RemainingExcludesMatchedKeys(t *testing.T) {
	result := renamer.Detect(sampleChanges())

	for _, c := range result.Remaining {
		if c.Key == "DB_HOST" || c.Key == "DATABASE_HOST" {
			t.Errorf("matched key %s should not appear in Remaining", c.Key)
		}
	}
}

func TestDetect_RemainingContainsModified(t *testing.T) {
	result := renamer.Detect(sampleChanges())

	found := false
	for _, c := range result.Remaining {
		if c.Key == "APP_ENV" {
			found = true
		}
	}
	if !found {
		t.Error("expected APP_ENV in Remaining")
	}
}

func TestDetect_NoRenameWhenValuesDoNotMatch(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Removed, Key: "OLD_KEY", OldValue: "foo"},
		{Type: differ.Added, Key: "NEW_KEY", NewValue: "bar"},
	}
	result := renamer.Detect(changes)

	if len(result.Hints) != 0 {
		t.Errorf("expected 0 hints, got %d", len(result.Hints))
	}
	if len(result.Remaining) != 2 {
		t.Errorf("expected 2 remaining changes, got %d", len(result.Remaining))
	}
}

func TestDetect_EmptyChanges(t *testing.T) {
	result := renamer.Detect([]differ.Change{})

	if len(result.Hints) != 0 {
		t.Errorf("expected 0 hints, got %d", len(result.Hints))
	}
	if len(result.Remaining) != 0 {
		t.Errorf("expected 0 remaining, got %d", len(result.Remaining))
	}
}

func TestDetect_EmptyValueNotMatched(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Removed, Key: "EMPTY_OLD", OldValue: ""},
		{Type: differ.Added, Key: "EMPTY_NEW", NewValue: ""},
	}
	result := renamer.Detect(changes)

	if len(result.Hints) != 0 {
		t.Errorf("empty values should not be treated as renames, got %d hints", len(result.Hints))
	}
}
