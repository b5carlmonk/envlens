package tracer

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/envlens/internal/differ"
)

func buildTrace() *Trace {
	tr := New()
	tr.Add("v1.0", "dev.env", "prod.env", []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "db.prod.local"},
		{Key: "OLD_KEY", Type: differ.Removed, OldValue: "legacy"},
	})
	return tr
}

func TestRenderText_ContainsHeader(t *testing.T) {
	tr := buildTrace()
	out := RenderText(tr)
	if !strings.Contains(out, "Trace Log") {
		t.Error("expected 'Trace Log' header in output")
	}
}

func TestRenderText_ContainsLabel(t *testing.T) {
	tr := buildTrace()
	out := RenderText(tr)
	if !strings.Contains(out, "v1.0") {
		t.Error("expected label 'v1.0' in output")
	}
}

func TestRenderText_ContainsKeys(t *testing.T) {
	tr := buildTrace()
	out := RenderText(tr)
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key 'DB_HOST' in output")
	}
	if !strings.Contains(out, "OLD_KEY") {
		t.Error("expected key 'OLD_KEY' in output")
	}
}

func TestRenderText_EmptyTrace(t *testing.T) {
	tr := New()
	out := RenderText(tr)
	if !strings.Contains(out, "No trace entries") {
		t.Error("expected 'No trace entries' message")
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	tr := buildTrace()
	out, err := RenderJSON(tr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(parsed) != 1 {
		t.Errorf("expected 1 entry, got %d", len(parsed))
	}
}

func TestRenderJSON_ContainsLabel(t *testing.T) {
	tr := buildTrace()
	out, err := RenderJSON(tr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "v1.0") {
		t.Error("expected label 'v1.0' in JSON output")
	}
}

func TestRenderJSON_EmptyTrace_ReturnsArray(t *testing.T) {
	tr := New()
	out, err := RenderJSON(tr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[") {
		t.Error("expected JSON array in output")
	}
}
