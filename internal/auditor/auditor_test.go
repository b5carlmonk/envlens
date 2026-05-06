package auditor

import (
	"testing"

	"github.com/yourusername/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "APP_NAME", Type: differ.Added, NewValue: "envlens"},
		{Key: "DB_PASSWORD", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "LEGACY_FLAG", Type: differ.Removed, OldValue: "true"},
	}
}

func TestRun_ReturnsResult(t *testing.T) {
	result := Run("old.env", "new.env", sampleChanges())
	if result.Source != "old.env" {
		t.Errorf("expected source old.env, got %s", result.Source)
	}
	if result.Target != "new.env" {
		t.Errorf("expected target new.env, got %s", result.Target)
	}
	if len(result.Changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(result.Changes))
	}
	if result.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRun_AnnotatesSensitiveModified(t *testing.T) {
	result := Run("a", "b", sampleChanges())
	var found bool
	for _, ann := range result.Annotations {
		if ann.Key == "DB_PASSWORD" && ann.Severity == "critical" {
			found = true
		}
	}
	if !found {
		t.Error("expected critical annotation for DB_PASSWORD")
	}
}

func TestRun_AnnotatesRemovedKey(t *testing.T) {
	result := Run("a", "b", sampleChanges())
	var found bool
	for _, ann := range result.Annotations {
		if ann.Key == "LEGACY_FLAG" && ann.Severity == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning annotation for LEGACY_FLAG")
	}
}

func TestRun_NoAnnotationForAddedNonSensitive(t *testing.T) {
	result := Run("a", "b", sampleChanges())
	for _, ann := range result.Annotations {
		if ann.Key == "APP_NAME" {
			t.Errorf("unexpected annotation for APP_NAME: %+v", ann)
		}
	}
}

func TestRun_EmptyChanges(t *testing.T) {
	result := Run("a", "b", []differ.Change{})
	if len(result.Annotations) != 0 {
		t.Errorf("expected no annotations, got %d", len(result.Annotations))
	}
}

func TestIsSensitiveKey(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"APP_NAME", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := isSensitiveKey(tc.key)
		if got != tc.want {
			t.Errorf("isSensitiveKey(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
