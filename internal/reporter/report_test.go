package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/reporter"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Added, Key: "NEW_KEY", NewValue: "hello"},
		{Type: differ.Removed, Key: "OLD_KEY", OldValue: "bye"},
		{Type: differ.Modified, Key: "MOD_KEY", OldValue: "v1", NewValue: "v2"},
	}
}

func TestReport_HasChanges(t *testing.T) {
	r := reporter.NewReport("a.env", "b.env", sampleChanges())
	if !r.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestReport_NoChanges(t *testing.T) {
	r := reporter.NewReport("a.env", "b.env", nil)
	if r.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestReport_Summary(t *testing.T) {
	r := reporter.NewReport("a.env", "b.env", sampleChanges())
	summary := r.Summary()
	if !strings.Contains(summary, "+1") || !strings.Contains(summary, "-1") || !strings.Contains(summary, "~1") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestReport_RenderText(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewReport("a.env", "b.env", sampleChanges())
	if err := r.Render(&buf, reporter.FormatText); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"NEW_KEY", "OLD_KEY", "MOD_KEY", "a.env", "b.env"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestReport_RenderTextNoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewReport("a.env", "b.env", nil)
	if err := r.Render(&buf, reporter.FormatText); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Error("expected 'No changes' message")
	}
}

func TestReport_RenderJSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewReport("a.env", "b.env", sampleChanges())
	if err := r.Render(&buf, reporter.FormatJSON); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if payload["from_file"] != "a.env" {
		t.Errorf("expected from_file=a.env, got %v", payload["from_file"])
	}
	changes, ok := payload["changes"].([]interface{})
	if !ok || len(changes) != 3 {
		t.Errorf("expected 3 changes, got %v", payload["changes"])
	}
}

func TestReport_RenderUnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewReport("a.env", "b.env", sampleChanges())
	if err := r.Render(&buf, reporter.Format("xml")); err == nil {
		t.Error("expected error for unknown format")
	}
}
