package baseline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/envlens/internal/baseline"
)

func buildBaseline() *baseline.Baseline {
	b := baseline.New("test-base", ".env.test", map[string]string{
		"APP_ENV": "staging",
		"PORT":    "8080",
		"DB_URL":  "postgres://localhost/db",
	})
	b.CreatedAt = time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)
	return b
}

func TestRenderText_ContainsName(t *testing.T) {
	var buf bytes.Buffer
	baseline.RenderText(buildBaseline(), nil, &buf)
	if !strings.Contains(buf.String(), "test-base") {
		t.Errorf("expected baseline name in output")
	}
}

func TestRenderText_ContainsSource(t *testing.T) {
	var buf bytes.Buffer
	baseline.RenderText(buildBaseline(), nil, &buf)
	if !strings.Contains(buf.String(), ".env.test") {
		t.Errorf("expected source in output")
	}
}

func TestRenderText_NoDrift_ShowsNone(t *testing.T) {
	var buf bytes.Buffer
	baseline.RenderText(buildBaseline(), []string{}, &buf)
	if !strings.Contains(buf.String(), "none") {
		t.Errorf("expected 'none' when no drift")
	}
}

func TestRenderText_WithDrift_ShowsKeys(t *testing.T) {
	var buf bytes.Buffer
	baseline.RenderText(buildBaseline(), []string{"APP_ENV", "PORT"}, &buf)
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in drift output")
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in drift output")
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := baseline.RenderJSON(buildBaseline(), []string{"DB_URL"}, &buf)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}
}

func TestRenderJSON_HasDrift_True(t *testing.T) {
	var buf bytes.Buffer
	baseline.RenderJSON(buildBaseline(), []string{"APP_ENV"}, &buf)
	var raw map[string]interface{}
	json.Unmarshal(buf.Bytes(), &raw)
	if raw["has_drift"] != true {
		t.Errorf("expected has_drift=true")
	}
}

func TestRenderJSON_NoDrift_EmptyArray(t *testing.T) {
	var buf bytes.Buffer
	baseline.RenderJSON(buildBaseline(), nil, &buf)
	var raw map[string]interface{}
	json.Unmarshal(buf.Bytes(), &raw)
	drift, ok := raw["drift"].([]interface{})
	if !ok || len(drift) != 0 {
		t.Errorf("expected empty drift array")
	}
}
