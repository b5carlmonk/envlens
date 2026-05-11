package patcher_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/envlens/internal/patcher"
)

func buildResult() patcher.Result {
	return patcher.Result{
		Applied:   []string{"DB_HOST", "API_KEY"},
		Skipped:   []string{"OLD_KEY"},
		Conflicts: []string{"DB_HOST"},
		DryRun:    false,
	}
}

func TestRenderText_ContainsAppliedCount(t *testing.T) {
	out := patcher.RenderText(buildResult())
	if !strings.Contains(out, "Applied  : 2") {
		t.Errorf("expected applied count, got:\n%s", out)
	}
}

func TestRenderText_ContainsSkippedKeys(t *testing.T) {
	out := patcher.RenderText(buildResult())
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in skipped section, got:\n%s", out)
	}
}

func TestRenderText_DryRunLabel(t *testing.T) {
	r := buildResult()
	r.DryRun = true
	out := patcher.RenderText(r)
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run label, got:\n%s", out)
	}
}

func TestRenderText_ConflictSection(t *testing.T) {
	out := patcher.RenderText(buildResult())
	if !strings.Contains(out, "Conflict keys") {
		t.Errorf("expected conflict section, got:\n%s", out)
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	out, err := patcher.RenderJSON(buildResult())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestRenderJSON_ContainsAppliedKeys(t *testing.T) {
	out, err := patcher.RenderJSON(buildResult())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in JSON output, got:\n%s", out)
	}
}

func TestRenderJSON_EmptySlicesNotNull(t *testing.T) {
	r := patcher.Result{}
	out, err := patcher.RenderJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "null") {
		t.Errorf("expected no null slices in JSON, got:\n%s", out)
	}
}
