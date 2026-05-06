package auditor

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/differ"
)

func buildResult() AuditResult {
	changes := []differ.Change{
		{Key: "DB_PASSWORD", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "LEGACY", Type: differ.Removed, OldValue: "1"},
	}
	return Run("base.env", "next.env", changes)
}

func TestRenderText_ContainsSourceAndTarget(t *testing.T) {
	var buf bytes.Buffer
	RenderText(&buf, buildResult())
	out := buf.String()
	if !strings.Contains(out, "base.env") {
		t.Error("expected source in output")
	}
	if !strings.Contains(out, "next.env") {
		t.Error("expected target in output")
	}
}

func TestRenderText_ContainsRiskLevel(t *testing.T) {
	var buf bytes.Buffer
	RenderText(&buf, buildResult())
	out := buf.String()
	if !strings.Contains(out, "Risk") {
		t.Error("expected Risk line in output")
	}
}

func TestRenderText_ContainsAnnotations(t *testing.T) {
	var buf bytes.Buffer
	RenderText(&buf, buildResult())
	out := buf.String()
	if !strings.Contains(out, "CRITICAL") && !strings.Contains(out, "WARNING") {
		t.Error("expected annotation severity tags in output")
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := RenderJSON(&buf, buildResult())
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestRenderJSON_ContainsSourceKey(t *testing.T) {
	var buf bytes.Buffer
	_ = RenderJSON(&buf, buildResult())
	var out map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)
	if _, ok := out["Source"]; !ok {
		t.Error("expected Source key in JSON output")
	}
}
