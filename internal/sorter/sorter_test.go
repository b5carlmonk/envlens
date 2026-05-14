package sorter_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/differ"
	"github.com/yourusername/envlens/internal/sorter"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "ZEBRA_URL", Type: differ.Modified, OldValue: "a", NewValue: "b"},
		{Key: "APP_NAME", Type: differ.Unchanged, OldValue: "app", NewValue: "app"},
		{Key: "DB_PASSWORD", Type: differ.Removed, OldValue: "secret", NewValue: ""},
		{Key: "NEW_FEATURE", Type: differ.Added, OldValue: "", NewValue: "true"},
		{Key: "ALPHA_KEY", Type: differ.Added, OldValue: "", NewValue: "1"},
	}
}

func TestApply_SortByKey_Ascending(t *testing.T) {
	result := sorter.Apply(sampleChanges(), sorter.Options{By: sorter.SortByKey})

	expected := []string{"ALPHA_KEY", "APP_NAME", "DB_PASSWORD", "NEW_FEATURE", "ZEBRA_URL"}
	for i, key := range expected {
		if result[i].Key != key {
			t.Errorf("index %d: expected key %q, got %q", i, key, result[i].Key)
		}
	}
}

func TestApply_SortByKey_Descending(t *testing.T) {
	result := sorter.Apply(sampleChanges(), sorter.Options{By: sorter.SortByKey, Descending: true})

	expected := []string{"ZEBRA_URL", "NEW_FEATURE", "DB_PASSWORD", "APP_NAME", "ALPHA_KEY"}
	for i, key := range expected {
		if result[i].Key != key {
			t.Errorf("index %d: expected key %q, got %q", i, key, result[i].Key)
		}
	}
}

func TestApply_SortByType_GroupsCorrectly(t *testing.T) {
	result := sorter.Apply(sampleChanges(), sorter.Options{By: sorter.SortByType})

	if result[0].Type != differ.Added || result[1].Type != differ.Added {
		t.Errorf("expected first two entries to be Added, got %v and %v", result[0].Type, result[1].Type)
	}
	if result[2].Type != differ.Removed {
		t.Errorf("expected third entry to be Removed, got %v", result[2].Type)
	}
	if result[3].Type != differ.Modified {
		t.Errorf("expected fourth entry to be Modified, got %v", result[3].Type)
	}
	if result[4].Type != differ.Unchanged {
		t.Errorf("expected last entry to be Unchanged, got %v", result[4].Type)
	}
}

func TestApply_SortByType_WithinTypeAlphabetical(t *testing.T) {
	result := sorter.Apply(sampleChanges(), sorter.Options{By: sorter.SortByType})

	// Both Added entries: ALPHA_KEY and NEW_FEATURE — ALPHA_KEY should come first.
	if result[0].Key != "ALPHA_KEY" {
		t.Errorf("expected ALPHA_KEY first among Added, got %q", result[0].Key)
	}
	if result[1].Key != "NEW_FEATURE" {
		t.Errorf("expected NEW_FEATURE second among Added, got %q", result[1].Key)
	}
}

func TestApply_EmptyChanges_ReturnsEmpty(t *testing.T) {
	result := sorter.Apply([]differ.Change{}, sorter.Options{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := sampleChanges()
	firstKey := original[0].Key

	sorter.Apply(original, sorter.Options{By: sorter.SortByKey})

	if original[0].Key != firstKey {
		t.Errorf("original slice was mutated: expected %q at index 0, got %q", firstKey, original[0].Key)
	}
}

func TestApply_DefaultSortBy_FallsBackToKey(t *testing.T) {
	result := sorter.Apply(sampleChanges(), sorter.Options{})

	if result[0].Key != "ALPHA_KEY" {
		t.Errorf("expected ALPHA_KEY first with default sort, got %q", result[0].Key)
	}
}
