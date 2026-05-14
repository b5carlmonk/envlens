package duplicator_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/duplicator"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"DB_HOST":      "localhost",
		"REDIS_HOST":   "localhost",
		"API_KEY":      "secret123",
		"SERVICE_KEY":  "secret123",
		"UNIQUE_KEY":   "only-one",
		"EMPTY_KEY":    "",
	}
}

func TestDetect_FindsDuplicateValues(t *testing.T) {
	result := duplicator.Detect("test.env", sampleEnv())
	if len(result.Matches) != 2 {
		t.Fatalf("expected 2 duplicate groups, got %d", len(result.Matches))
	}
}

func TestDetect_SkipsEmptyValues(t *testing.T) {
	env := map[string]string{
		"A": "",
		"B": "",
	}
	result := duplicator.Detect("test.env", env)
	if len(result.Matches) != 0 {
		t.Errorf("expected no matches for empty values, got %d", len(result.Matches))
	}
}

func TestDetect_NoChanges_WhenAllUnique(t *testing.T) {
	env := map[string]string{
		"FOO": "alpha",
		"BAR": "beta",
		"BAZ": "gamma",
	}
	result := duplicator.Detect("prod.env", env)
	if len(result.Matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(result.Matches))
	}
}

func TestDetect_MatchContainsCorrectKeys(t *testing.T) {
	env := map[string]string{
		"DB_HOST":    "localhost",
		"CACHE_HOST": "localhost",
	}
	result := duplicator.Detect("dev.env", env)
	if len(result.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(result.Matches))
	}
	m := result.Matches[0]
	if len(m.Keys) != 2 {
		t.Errorf("expected 2 keys in match, got %d", len(m.Keys))
	}
	if m.Value != "localhost" {
		t.Errorf("expected value 'localhost', got %q", m.Value)
	}
}

func TestDetect_SetsSourceField(t *testing.T) {
	result := duplicator.Detect("staging.env", map[string]string{})
	if result.Source != "staging.env" {
		t.Errorf("expected source 'staging.env', got %q", result.Source)
	}
}

func TestRenderText_ContainsSource(t *testing.T) {
	result := duplicator.Detect("prod.env", sampleEnv())
	out := duplicator.RenderText(result)
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected output to contain source name")
	}
}

func TestRenderText_NoMatches_ShowsNone(t *testing.T) {
	result := duplicator.Detect("clean.env", map[string]string{"FOO": "bar"})
	out := duplicator.RenderText(result)
	if !strings.Contains(out, "No duplicate") {
		t.Errorf("expected 'No duplicate' message in output")
	}
}

func TestRenderText_ContainsDuplicateKeys(t *testing.T) {
	env := map[string]string{
		"DB_HOST":    "localhost",
		"CACHE_HOST": "localhost",
	}
	result := duplicator.Detect("dev.env", env)
	out := duplicator.RenderText(result)
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected output to contain duplicate value")
	}
	if !strings.Contains(out, "CACHE_HOST") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected output to contain both duplicate keys")
	}
}
