package baseline_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envlens/internal/baseline"
)

func TestNew_SetsFieldsCorrectly(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	b := baseline.New("prod-baseline", ".env.prod", env)

	if b.Name != "prod-baseline" {
		t.Errorf("expected name 'prod-baseline', got %q", b.Name)
	}
	if b.Source != ".env.prod" {
		t.Errorf("expected source '.env.prod', got %q", b.Source)
	}
	if b.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production")
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNew_CopiesEnvMap(t *testing.T) {
	env := map[string]string{"KEY": "original"}
	b := baseline.New("test", "src", env)
	env["KEY"] = "mutated"

	if b.Env["KEY"] != "original" {
		t.Error("expected baseline env to be independent copy")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	b := baseline.New("db-baseline", ".env", env)
	b.CreatedAt = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	path := filepath.Join(t.TempDir(), "baseline.json")
	if err := baseline.Save(b, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != b.Name {
		t.Errorf("name mismatch: got %q", loaded.Name)
	}
	if loaded.Env["DB_HOST"] != "localhost" {
		t.Errorf("env mismatch for DB_HOST")
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte("not-json"), 0644)
	_, err := baseline.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSave_ProducesValidJSON(t *testing.T) {
	b := baseline.New("check", ".env", map[string]string{"X": "1"})
	path := filepath.Join(t.TempDir(), "out.json")
	baseline.Save(b, path)

	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("saved file is not valid JSON: %v", err)
	}
}

func TestDriftKeys_DetectsChangedValues(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging", "PORT": "3000", "STABLE": "yes"}
	b := baseline.New("base", ".env", env)

	current := map[string]string{"APP_ENV": "production", "PORT": "3000", "STABLE": "yes"}
	drifted := baseline.DriftKeys(b, current)

	if len(drifted) != 1 || drifted[0] != "APP_ENV" {
		t.Errorf("expected [APP_ENV] drifted, got %v", drifted)
	}
}

func TestDriftKeys_NoDrift_ReturnsEmpty(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	b := baseline.New("base", ".env", env)
	drifted := baseline.DriftKeys(b, env)
	if len(drifted) != 0 {
		t.Errorf("expected no drift, got %v", drifted)
	}
}
