package annotator

import (
	"testing"

	"github.com/yourusername/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_PASSWORD", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "APP_PORT", Type: differ.Added, NewValue: "8080"},
		{Key: "LEGACY_HOST", Type: differ.Removed, OldValue: "localhost"},
		{Key: "_INTERNAL_FLAG", Type: differ.Added, NewValue: "true"},
		{Key: "DATABASE_URL", Type: differ.Modified, OldValue: "a", NewValue: "b"},
	}
}

func TestApply_ReturnsSameCount(t *testing.T) {
	changes := sampleChanges()
	result := Apply(changes, Options{})
	if len(result.Annotations) != len(changes) {
		t.Errorf("expected %d annotations, got %d", len(changes), len(result.Annotations))
	}
}

func TestApply_SensitiveKey_TaggedSensitive(t *testing.T) {
	changes := []differ.Change{
		{Key: "DB_PASSWORD", Type: differ.Modified, OldValue: "x", NewValue: "y"},
	}
	result := Apply(changes, Options{})
	if result.Annotations[0].Tag != TagSensitive {
		t.Errorf("expected sensitive tag, got %s", result.Annotations[0].Tag)
	}
}

func TestApply_DeprecatedKey_TaggedDeprecated(t *testing.T) {
	changes := []differ.Change{
		{Key: "LEGACY_HOST", Type: differ.Removed, OldValue: "localhost"},
	}
	opts := Options{DeprecatedKeys: []string{"LEGACY"}}
	result := Apply(changes, opts)
	if result.Annotations[0].Tag != TagDeprecated {
		t.Errorf("expected deprecated tag, got %s", result.Annotations[0].Tag)
	}
}

func TestApply_InternalPrefix_TaggedInternal(t *testing.T) {
	changes := []differ.Change{
		{Key: "_INTERNAL_FLAG", Type: differ.Added, NewValue: "true"},
	}
	opts := Options{InternalPrefixes: []string{"_INTERNAL"}}
	result := Apply(changes, opts)
	if result.Annotations[0].Tag != TagInternal {
		t.Errorf("expected internal tag, got %s", result.Annotations[0].Tag)
	}
}

func TestApply_RequiredKey_TaggedRequired(t *testing.T) {
	changes := []differ.Change{
		{Key: "DATABASE_URL", Type: differ.Modified, OldValue: "a", NewValue: "b"},
	}
	opts := Options{RequiredKeys: []string{"DATABASE_URL"}}
	result := Apply(changes, opts)
	if result.Annotations[0].Tag != TagRequired {
		t.Errorf("expected required tag, got %s", result.Annotations[0].Tag)
	}
}

func TestApply_UnknownKey_TaggedUnknown(t *testing.T) {
	changes := []differ.Change{
		{Key: "APP_PORT", Type: differ.Added, NewValue: "8080"},
	}
	result := Apply(changes, Options{})
	if result.Annotations[0].Tag != TagUnknown {
		t.Errorf("expected unknown tag, got %s", result.Annotations[0].Tag)
	}
}

func TestApply_EmptyChanges_ReturnsEmpty(t *testing.T) {
	result := Apply([]differ.Change{}, Options{})
	if len(result.Annotations) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(result.Annotations))
	}
}
