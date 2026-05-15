package scorer

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func buildResult(t *testing.T) ScoreResult {
	t.Helper()
	changes := []differ.Change{
		{Type: differ.Modified, Key: "SECRET_KEY", OldValue: "old", NewValue: "new"},
		{Type: differ.Removed, Key: "DB_PASSWORD", OldValue: "pass", NewValue: ""},
		{Type: differ.Added, Key: "APP_ENV", OldValue: "", NewValue: "production"},
	}
	return Evaluate(changes)
}

func TestRenderText_ContainsRiskLevel(t *testing.T) {
	result := buildResult(t)
	out := RenderText(result)
	if !strings.Contains(out, "Risk Level") {
		t.Errorf("expected output to contain 'Risk Level', got:\n%s", out)
	}
}

func TestRenderText_ContainsScore(t *testing.T) {
	result := buildResult(t)
	out := RenderText(result)
	if !strings.Contains(out, "Score") {
		t.Errorf("expected output to contain 'Score', got:\n%s", out)
	}
}

func TestRenderText_ContainsReasons(t *testing.T) {
	result := buildResult(t)
	out := RenderText(result)
	if !strings.Contains(out, "Reasons") {
		t.Errorf("expected output to contain 'Reasons', got:\n%s", out)
	}
}

func TestRenderText_NoChanges_ShowsNone(t *testing.T) {
	result := Evaluate(nil)
	out := RenderText(result)
	if !strings.Contains(out, "none") {
		t.Errorf("expected 'none' for empty reasons, got:\n%s", out)
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	result := buildResult(t)
	out, err := RenderJSON(result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
}

func TestRenderJSON_ContainsLevel(t *testing.T) {
	result := buildResult(t)
	out, err := RenderJSON(result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "level") && !strings.Contains(out, "Level") {
		t.Errorf("expected JSON to contain level field, got:\n%s", out)
	}
}
