package merger

import (
	"encoding/json"
	"strings"
	"testing"
)

func buildResult() Result {
	return Result{
		Merged: map[string]string{
			"APP_ENV":  "production",
			"DB_HOST":  "localhost",
			"API_KEY":  "resolved-key",
		},
		Conflicts: []Conflict{
			{Key: "API_KEY", SourceValue: "src-key", TargetValue: "tgt-key", Resolved: "resolved-key"},
		},
	}
}

func TestRenderText_ContainsMergedCount(t *testing.T) {
	out := RenderText(buildResult())
	if !strings.Contains(out, "Merged keys: 3") {
		t.Errorf("expected merged key count, got:\n%s", out)
	}
}

func TestRenderText_ContainsConflictCount(t *testing.T) {
	out := RenderText(buildResult())
	if !strings.Contains(out, "Conflicts resolved: 1") {
		t.Errorf("expected conflict count, got:\n%s", out)
	}
}

func TestRenderText_ContainsConflictDetail(t *testing.T) {
	out := RenderText(buildResult())
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in conflict detail, got:\n%s", out)
	}
}

func TestRenderText_NoConflicts_NoConflictSection(t *testing.T) {
	r := Result{Merged: map[string]string{"FOO": "bar"}, Conflicts: nil}
	out := RenderText(r)
	if strings.Contains(out, "Conflicts:\n") {
		t.Errorf("unexpected conflict section in output:\n%s", out)
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	out, err := RenderJSON(buildResult())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
}

func TestRenderJSON_ContainsMergedAndConflicts(t *testing.T) {
	out, err := RenderJSON(buildResult())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"merged\"") {
		t.Errorf("expected merged key in JSON, got:\n%s", out)
	}
	if !strings.Contains(out, "\"conflicts\"") {
		t.Errorf("expected conflicts key in JSON, got:\n%s", out)
	}
}
