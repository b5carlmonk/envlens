package inspector_test

import (
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/inspector"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Added, Key: "PORT", OldValue: "", NewValue: "8080"},
		{Type: differ.Modified, Key: "DB_URL", OldValue: "postgres://old", NewValue: "postgres://new"},
		{Type: differ.Modified, Key: "DEBUG", OldValue: "false", NewValue: "true"},
		{Type: differ.Removed, Key: "SECRET", OldValue: "abc123", NewValue: ""},
	}
}

func TestInspect_ReturnsSameCount(t *testing.T) {
	changes := sampleChanges()
	result := inspector.Inspect(changes)
	if len(result) != len(changes) {
		t.Fatalf("expected %d inspections, got %d", len(changes), len(result))
	}
}

func TestInspect_NumericHint(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Added, Key: "PORT", OldValue: "", NewValue: "8080"},
	}
	result := inspector.Inspect(changes)
	if result[0].NewHint != inspector.TypeNumeric {
		t.Errorf("expected TypeNumeric, got %s", result[0].NewHint)
	}
}

func TestInspect_BooleanHint(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Modified, Key: "DEBUG", OldValue: "false", NewValue: "true"},
	}
	result := inspector.Inspect(changes)
	if result[0].OldHint != inspector.TypeBoolean {
		t.Errorf("expected TypeBoolean for old, got %s", result[0].OldHint)
	}
	if result[0].NewHint != inspector.TypeBoolean {
		t.Errorf("expected TypeBoolean for new, got %s", result[0].NewHint)
	}
}

func TestInspect_URLHint(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Modified, Key: "DB_URL", OldValue: "postgres://old", NewValue: "https://example.com"},
	}
	result := inspector.Inspect(changes)
	if result[0].NewHint != inspector.TypeURL {
		t.Errorf("expected TypeURL, got %s", result[0].NewHint)
	}
}

func TestInspect_LengthDelta(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Modified, Key: "TOKEN", OldValue: "short", NewValue: "muchlongervalue"},
	}
	result := inspector.Inspect(changes)
	expected := len("muchlongervalue") - len("short")
	if result[0].LengthDelta != expected {
		t.Errorf("expected delta %d, got %d", expected, result[0].LengthDelta)
	}
}

func TestInspect_ValueSame_WhenUnchanged(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Modified, Key: "KEY", OldValue: "same", NewValue: "same"},
	}
	result := inspector.Inspect(changes)
	if !result[0].ValueSame {
		t.Error("expected ValueSame to be true")
	}
}

func TestInspect_StringHint_Default(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Added, Key: "APP_NAME", OldValue: "", NewValue: "myapp"},
	}
	result := inspector.Inspect(changes)
	if result[0].NewHint != inspector.TypeString {
		t.Errorf("expected TypeString, got %s", result[0].NewHint)
	}
}
