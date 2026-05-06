package filter_test

import (
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/filter"
)

var sampleChanges = []differ.Change{
	{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
	{Key: "DB_PORT", Type: differ.Modified, OldValue: "5432", NewValue: "5433"},
	{Key: "APP_SECRET", Type: differ.Removed, OldValue: "abc123"},
	{Key: "APP_NAME", Type: differ.Added, NewValue: "envlens"},
	{Key: "REDIS_URL", Type: differ.Modified, OldValue: "redis://old", NewValue: "redis://new"},
}

func TestFilter_ByPrefix(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{Prefix: "DB_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 changes with prefix DB_, got %d", len(result))
	}
	for _, c := range result {
		if c.Key != "DB_HOST" && c.Key != "DB_PORT" {
			t.Errorf("unexpected key %q in prefix-filtered results", c.Key)
		}
	}
}

func TestFilter_ByKeyContains(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{KeyContains: "APP"})
	if len(result) != 2 {
		t.Fatalf("expected 2 changes containing APP, got %d", len(result))
	}
}

func TestFilter_ByType_Added(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{Types: []string{"added"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 added changes, got %d", len(result))
	}
	for _, c := range result {
		if c.Type != differ.Added {
			t.Errorf("expected Added type, got %q", c.Type)
		}
	}
}

func TestFilter_ByType_MultipleTypes(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{Types: []string{"removed", "modified"}})
	if len(result) != 3 {
		t.Fatalf("expected 3 changes for removed+modified, got %d", len(result))
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{})
	if len(result) != len(sampleChanges) {
		t.Fatalf("expected all %d changes, got %d", len(sampleChanges), len(result))
	}
}

func TestFilter_PrefixAndType_Combined(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{
		Prefix: "APP_",
		Types:  []string{"added"},
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %q", result[0].Key)
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	result := filter.Apply(nil, filter.Options{Prefix: "DB_"})
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %v", result)
	}
}
