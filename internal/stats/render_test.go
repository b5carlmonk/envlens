package stats_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/stats"
)

func buildResult() stats.Result {
	return stats.Result{
		Total:    5,
		Added:    2,
		Removed:  1,
		Modified: 2,
		TopPrefixes: []stats.PrefixCount{
			{Prefix: "DB", Count: 2},
			{Prefix: "AUTH", Count: 1},
		},
	}
}

func TestRenderText_ContainsHeader(t *testing.T) {
	out := stats.RenderText(buildResult())
	if !strings.Contains(out, "Statistics") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestRenderText_ContainsCounts(t *testing.T) {
	out := stats.RenderText(buildResult())
	for _, want := range []string{"Total", "Added", "Removed", "Modified"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in text output", want)
		}
	}
}

func TestRenderText_ContainsTopPrefixes(t *testing.T) {
	out := stats.RenderText(buildResult())
	if !strings.Contains(out, "DB") {
		t.Errorf("expected prefix DB in output, got:\n%s", out)
	}
}

func TestRenderText_EmptyPrefixes_NoPrefixSection(t *testing.T) {
	r := stats.Result{Total: 1, Added: 1}
	out := stats.RenderText(r)
	if strings.Contains(out, "Top key prefixes") {
		t.Error("expected no prefix section when TopPrefixes is empty")
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	out, err := stats.RenderJSON(buildResult())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
}

func TestRenderJSON_ContainsFields(t *testing.T) {
	out, err := stats.RenderJSON(buildResult())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, field := range []string{"total", "added", "removed", "modified", "top_prefixes"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in JSON output", field)
		}
	}
}
