package trimmer_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/differ"
	"github.com/yourusername/envlens/internal/trimmer"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"CLEAN_KEY":       "no-whitespace",
		"LEADING_SPACE":   "  hello",
		"TRAILING_SPACE":  "world   ",
		"BOTH_SPACES":     "  both  ",
		"EMPTY_VALUE":     "",
	}
}

func TestApply_ReturnsCleanMap(t *testing.T) {
	r := trimmer.Apply(sampleEnv())

	if got := r.Clean["CLEAN_KEY"]; got != "no-whitespace" {
		t.Errorf("expected 'no-whitespace', got %q", got)
	}
	if got := r.Clean["LEADING_SPACE"]; got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
	if got := r.Clean["TRAILING_SPACE"]; got != "world" {
		t.Errorf("expected 'world', got %q", got)
	}
	if got := r.Clean["BOTH_SPACES"]; got != "both" {
		t.Errorf("expected 'both', got %q", got)
	}
}

func TestApply_TrimmedContainsOnlyDirtyKeys(t *testing.T) {
	r := trimmer.Apply(sampleEnv())

	if _, ok := r.Trimmed["CLEAN_KEY"]; ok {
		t.Error("CLEAN_KEY should not appear in Trimmed")
	}
	if _, ok := r.Trimmed["EMPTY_VALUE"]; ok {
		t.Error("EMPTY_VALUE should not appear in Trimmed")
	}

	expected := []string{"LEADING_SPACE", "TRAILING_SPACE", "BOTH_SPACES"}
	for _, key := range expected {
		if _, ok := r.Trimmed[key]; !ok {
			t.Errorf("expected %q in Trimmed", key)
		}
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	env := sampleEnv()
	trimmer.Apply(env)

	if env["LEADING_SPACE"] != "  hello" {
		t.Error("Apply must not mutate the original map")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	r := trimmer.Apply(map[string]string{})

	if len(r.Trimmed) != 0 {
		t.Errorf("expected 0 trimmed keys, got %d", len(r.Trimmed))
	}
	if len(r.Clean) != 0 {
		t.Errorf("expected 0 clean keys, got %d", len(r.Clean))
	}
}

func TestToChanges_ReturnsModifiedType(t *testing.T) {
	r := trimmer.Apply(sampleEnv())
	changes := trimmer.ToChanges(r)

	if len(changes) != len(r.Trimmed) {
		t.Errorf("expected %d changes, got %d", len(r.Trimmed), len(changes))
	}
	for _, c := range changes {
		if c.Type != differ.Modified {
			t.Errorf("expected Modified type for key %q, got %v", c.Key, c.Type)
		}
		if c.OldVal == c.NewVal {
			t.Errorf("OldVal and NewVal should differ for key %q", c.Key)
		}
	}
}
