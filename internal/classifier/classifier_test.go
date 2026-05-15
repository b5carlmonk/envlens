package classifier

import (
	"testing"

	"github.com/user/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", OldValue: "", NewValue: "localhost", Type: differ.Added},
		{Key: "AUTH_TOKEN", OldValue: "old", NewValue: "new", Type: differ.Modified},
		{Key: "APP_PORT", OldValue: "8080", NewValue: "", Type: differ.Removed},
		{Key: "LOG_LEVEL", OldValue: "info", NewValue: "debug", Type: differ.Modified},
		{Key: "FEATURE_DARK_MODE", OldValue: "", NewValue: "true", Type: differ.Added},
		{Key: "OTEL_EXPORTER_ENDPOINT", OldValue: "", NewValue: "http://otel:4317", Type: differ.Added},
		{Key: "SOME_RANDOM_VAR", OldValue: "a", NewValue: "b", Type: differ.Modified},
	}
}

func TestApply_ReturnsSameCount(t *testing.T) {
	changes := sampleChanges()
	results := Apply(changes)
	if len(results) != len(changes) {
		t.Fatalf("expected %d results, got %d", len(changes), len(results))
	}
}

func TestApply_DatabaseCategory(t *testing.T) {
	results := Apply([]differ.Change{
		{Key: "DB_HOST", NewValue: "localhost", Type: differ.Added},
	})
	if results[0].Category != CategoryDatabase {
		t.Errorf("expected database, got %s", results[0].Category)
	}
}

func TestApply_AuthCategory(t *testing.T) {
	results := Apply([]differ.Change{
		{Key: "AUTH_TOKEN", OldValue: "x", NewValue: "y", Type: differ.Modified},
	})
	if results[0].Category != CategoryAuth {
		t.Errorf("expected auth, got %s", results[0].Category)
	}
}

func TestApply_NetworkCategory(t *testing.T) {
	results := Apply([]differ.Change{
		{Key: "APP_PORT", OldValue: "8080", NewValue: "", Type: differ.Removed},
	})
	if results[0].Category != CategoryNetwork {
		t.Errorf("expected network, got %s", results[0].Category)
	}
}

func TestApply_LoggingCategory(t *testing.T) {
	results := Apply([]differ.Change{
		{Key: "LOG_LEVEL", OldValue: "info", NewValue: "debug", Type: differ.Modified},
	})
	if results[0].Category != CategoryLogging {
		t.Errorf("expected logging, got %s", results[0].Category)
	}
}

func TestApply_FeatureCategory(t *testing.T) {
	results := Apply([]differ.Change{
		{Key: "FEATURE_DARK_MODE", NewValue: "true", Type: differ.Added},
	})
	if results[0].Category != CategoryFeature {
		t.Errorf("expected feature, got %s", results[0].Category)
	}
}

func TestApply_ObservabilityCategory(t *testing.T) {
	results := Apply([]differ.Change{
		{Key: "OTEL_EXPORTER_ENDPOINT", NewValue: "http://otel:4317", Type: differ.Added},
	})
	if results[0].Category != CategoryObservability {
		t.Errorf("expected observability, got %s", results[0].Category)
	}
}

func TestApply_UnknownCategory(t *testing.T) {
	results := Apply([]differ.Change{
		{Key: "SOME_RANDOM_VAR", OldValue: "a", NewValue: "b", Type: differ.Modified},
	})
	if results[0].Category != CategoryUnknown {
		t.Errorf("expected unknown, got %s", results[0].Category)
	}
}

func TestApply_EmptyChanges(t *testing.T) {
	results := Apply([]differ.Change{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
