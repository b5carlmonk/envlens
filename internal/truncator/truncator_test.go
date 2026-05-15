package truncator_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/truncator"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Added, Key: "SHORT_VAL", NewValue: "hello"},
		{Type: differ.Modified, Key: "LONG_VAL", OldValue: strings.Repeat("a", 80), NewValue: strings.Repeat("b", 100)},
		{Type: differ.Removed, Key: "GONE_KEY", OldValue: strings.Repeat("x", 70)},
	}
}

func TestApply_ShortValuesUnchanged(t *testing.T) {
	changes := sampleChanges()
	opts := truncator.DefaultOptions()
	result := truncator.Apply(changes, opts)

	if result[0].NewValue != "hello" {
		t.Errorf("expected 'hello', got %q", result[0].NewValue)
	}
}

func TestApply_LongNewValueTruncated(t *testing.T) {
	changes := sampleChanges()
	opts := truncator.DefaultOptions()
	result := truncator.Apply(changes, opts)

	v := result[1].NewValue
	if !strings.HasSuffix(v, "...") {
		t.Errorf("expected ellipsis suffix, got %q", v)
	}
	if len([]rune(v)) != opts.MaxLength+len([]rune(opts.Ellipsis)) {
		t.Errorf("unexpected length: %d", len([]rune(v)))
	}
}

func TestApply_LongOldValueTruncated(t *testing.T) {
	changes := sampleChanges()
	opts := truncator.DefaultOptions()
	result := truncator.Apply(changes, opts)

	v := result[1].OldValue
	if !strings.HasSuffix(v, "...") {
		t.Errorf("expected ellipsis suffix on OldValue, got %q", v)
	}
}

func TestApply_RemovedKeyOldValueTruncated(t *testing.T) {
	changes := sampleChanges()
	opts := truncator.DefaultOptions()
	result := truncator.Apply(changes, opts)

	v := result[2].OldValue
	if !strings.HasSuffix(v, "...") {
		t.Errorf("expected truncation for removed key, got %q", v)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	changes := sampleChanges()
	origLen := len([]rune(changes[1].NewValue))
	opts := truncator.DefaultOptions()
	truncator.Apply(changes, opts)

	if len([]rune(changes[1].NewValue)) != origLen {
		t.Error("Apply mutated the original slice")
	}
}

func TestApply_CustomMaxLength(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Added, Key: "K", NewValue: "abcdefghij"},
	}
	opts := truncator.Options{MaxLength: 5, Ellipsis: "~"}
	result := truncator.Apply(changes, opts)

	if result[0].NewValue != "abcde~" {
		t.Errorf("unexpected truncation result: %q", result[0].NewValue)
	}
}

func TestApply_ZeroMaxLengthUsesDefault(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Added, Key: "K", NewValue: strings.Repeat("z", 100)},
	}
	opts := truncator.Options{MaxLength: 0}
	result := truncator.Apply(changes, opts)

	if len([]rune(result[0].NewValue)) != 64+3 {
		t.Errorf("expected default max 64 + ellipsis, got length %d", len([]rune(result[0].NewValue)))
	}
}
