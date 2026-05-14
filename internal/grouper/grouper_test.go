package grouper_test

import (
	"testing"

	"github.com/envlens/internal/differ"
	"github.com/envlens/internal/grouper"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "DB_PORT", Type: differ.Added, NewValue: "5432"},
		{Key: "AUTH_SECRET", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "AUTH_TOKEN", Type: differ.Removed, OldValue: "token"},
		{Key: "STANDALONE", Type: differ.Added, NewValue: "yes"},
		{Key: "APP_ENV", Type: differ.Added, NewValue: "prod"},
	}
}

func TestByPrefix_GroupsCorrectly(t *testing.T) {
	result := grouper.ByPrefix(sampleChanges(), "_")

	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}

	names := make(map[string]int)
	for _, g := range result.Groups {
		names[g.Name] = len(g.Changes)
	}

	if names["DB"] != 2 {
		t.Errorf("expected DB group to have 2 changes, got %d", names["DB"])
	}
	if names["AUTH"] != 2 {
		t.Errorf("expected AUTH group to have 2 changes, got %d", names["AUTH"])
	}
	if names["APP"] != 1 {
		t.Errorf("expected APP group to have 1 change, got %d", names["APP"])
	}
}

func TestByPrefix_UngroupedWhenNoSeparator(t *testing.T) {
	result := grouper.ByPrefix(sampleChanges(), "_")

	if len(result.Ungrouped) != 1 {
		t.Fatalf("expected 1 ungrouped change, got %d", len(result.Ungrouped))
	}
	if result.Ungrouped[0].Key != "STANDALONE" {
		t.Errorf("expected STANDALONE in ungrouped, got %s", result.Ungrouped[0].Key)
	}
}

func TestByPrefix_EmptyChanges(t *testing.T) {
	result := grouper.ByPrefix([]differ.Change{}, "_")

	if len(result.Groups) != 0 {
		t.Errorf("expected no groups, got %d", len(result.Groups))
	}
	if len(result.Ungrouped) != 0 {
		t.Errorf("expected no ungrouped, got %d", len(result.Ungrouped))
	}
}

func TestByCustom_GroupsByLabel(t *testing.T) {
	labels := map[string][]string{
		"database": {"DB_"},
		"auth":     {"AUTH_"},
	}
	result := grouper.ByCustom(sampleChanges(), labels)

	if len(result.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result.Groups))
	}

	names := make(map[string]int)
	for _, g := range result.Groups {
		names[g.Name] = len(g.Changes)
	}

	if names["database"] != 2 {
		t.Errorf("expected database group to have 2 changes, got %d", names["database"])
	}
	if names["auth"] != 2 {
		t.Errorf("expected auth group to have 2 changes, got %d", names["auth"])
	}
}

func TestByCustom_UnmatchedGoToUngrouped(t *testing.T) {
	labels := map[string][]string{
		"database": {"DB_"},
	}
	result := grouper.ByCustom(sampleChanges(), labels)

	if len(result.Ungrouped) != 4 {
		t.Errorf("expected 4 ungrouped changes, got %d", len(result.Ungrouped))
	}
}

func TestByCustom_EmptyLabels_AllUngrouped(t *testing.T) {
	result := grouper.ByCustom(sampleChanges(), map[string][]string{})

	if len(result.Groups) != 0 {
		t.Errorf("expected no groups, got %d", len(result.Groups))
	}
	if len(result.Ungrouped) != len(sampleChanges()) {
		t.Errorf("expected all changes ungrouped")
	}
}
