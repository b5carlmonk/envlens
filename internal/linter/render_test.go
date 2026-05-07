package linter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

func buildResult() Result {
	changes := []differ.Change{
		{Key: "bad_key", Type: differ.Added, NewValue: ""},
	}
	return Lint(changes)
}

func TestRenderText_NoIssues(t *testing.T) {
	r := Result{}
	out := RenderText(r)
	if !strings.Contains(out, "No lint issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestRenderText_ContainsKey(t *testing.T) {
	r := buildResult()
	out := RenderText(r)
	if !strings.Contains(out, "bad_key") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestRenderText_ContainsIssueCount(t *testing.T) {
	r := buildResult()
	out := RenderText(r)
	if !strings.Contains(out, "issue") {
		t.Errorf("expected issue count in output, got: %s", out)
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	r := buildResult()
	out, err := RenderJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestRenderJSON_TotalField(t *testing.T) {
	r := buildResult()
	out, _ := RenderJSON(r)
	var parsed map[string]interface{}
	json.Unmarshal([]byte(out), &parsed)
	total, ok := parsed["total"].(float64)
	if !ok || total == 0 {
		t.Errorf("expected non-zero total in JSON, got: %v", parsed["total"])
	}
}

func TestRenderJSON_EmptyIssues(t *testing.T) {
	r := Result{}
	out, err := RenderJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"issues\": []") {
		t.Errorf("expected empty issues array, got: %s", out)
	}
}
