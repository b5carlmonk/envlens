package templater_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/templater"
)

func sampleTemplate() *templater.Template {
	return &templater.Template{
		Required: []string{"DB_HOST", "DB_PORT"},
		Optional: []string{"LOG_LEVEL"},
	}
}

func TestValidate_AllPresent_NoIssues(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"LOG_LEVEL": "info",
	}
	result := templater.Validate(sampleTemplate(), env)
	if result.HasIssues() {
		t.Errorf("expected no issues, got missing=%v unexpected=%v", result.MissingRequired, result.UnexpectedKeys)
	}
	if len(result.Present) != 3 {
		t.Errorf("expected 3 present keys, got %d", len(result.Present))
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
	}
	result := templater.Validate(sampleTemplate(), env)
	if len(result.MissingRequired) != 1 || result.MissingRequired[0] != "DB_PORT" {
		t.Errorf("expected DB_PORT missing, got %v", result.MissingRequired)
	}
}

func TestValidate_UnexpectedKey(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"UNKNOWN_KEY": "value",
	}
	result := templater.Validate(sampleTemplate(), env)
	if len(result.UnexpectedKeys) != 1 || result.UnexpectedKeys[0] != "UNKNOWN_KEY" {
		t.Errorf("expected UNKNOWN_KEY unexpected, got %v", result.UnexpectedKeys)
	}
}

func TestValidate_EmptyEnv_AllRequiredMissing(t *testing.T) {
	result := templater.Validate(sampleTemplate(), map[string]string{})
	if len(result.MissingRequired) != 2 {
		t.Errorf("expected 2 missing required keys, got %d", len(result.MissingRequired))
	}
}

func TestFromMap_ValidDirectives(t *testing.T) {
	m := map[string]string{
		"API_KEY":   "required",
		"LOG_LEVEL": "optional",
	}
	tmpl, err := templater.FromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Required) != 1 || tmpl.Required[0] != "API_KEY" {
		t.Errorf("expected API_KEY in Required, got %v", tmpl.Required)
	}
	if len(tmpl.Optional) != 1 || tmpl.Optional[0] != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL in Optional, got %v", tmpl.Optional)
	}
}

func TestFromMap_InvalidDirective_ReturnsError(t *testing.T) {
	m := map[string]string{
		"SOME_KEY": "mandatory",
	}
	_, err := templater.FromMap(m)
	if err == nil {
		t.Error("expected error for unknown directive, got nil")
	}
}

func TestHasIssues_FalseWhenClean(t *testing.T) {
	r := templater.ValidationResult{}
	if r.HasIssues() {
		t.Error("expected HasIssues() to be false for empty result")
	}
}
