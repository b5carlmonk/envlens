package redactor

import (
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Added, Key: "DB_PASSWORD", OldValue: "", NewValue: "supersecret"},
		{Type: differ.Modified, Key: "API_KEY", OldValue: "oldkey123", NewValue: "newkey456"},
		{Type: differ.Added, Key: "APP_PORT", OldValue: "", NewValue: "8080"},
		{Type: differ.Removed, Key: "SECRET_TOKEN", OldValue: "tok_abc", NewValue: ""},
		{Type: differ.Modified, Key: "LOG_LEVEL", OldValue: "info", NewValue: "debug"},
	}
}

func TestApply_MasksSensitiveKeys(t *testing.T) {
	changes := sampleChanges()
	result := Apply(changes, DefaultOptions())

	for _, c := range result {
		if c.Key == "DB_PASSWORD" || c.Key == "API_KEY" || c.Key == "SECRET_TOKEN" {
			if c.NewValue != "" && !strings.Contains(c.NewValue, "*") {
				t.Errorf("expected masked value for key %s, got %q", c.Key, c.NewValue)
			}
			if c.OldValue != "" && !strings.Contains(c.OldValue, "*") {
				t.Errorf("expected masked old value for key %s, got %q", c.Key, c.OldValue)
			}
		}
	}
}

func TestApply_PreservesNonSensitiveValues(t *testing.T) {
	changes := sampleChanges()
	result := Apply(changes, DefaultOptions())

	for _, c := range result {
		if c.Key == "APP_PORT" && c.NewValue != "8080" {
			t.Errorf("expected APP_PORT to remain unmasked, got %q", c.NewValue)
		}
		if c.Key == "LOG_LEVEL" && c.NewValue != "debug" {
			t.Errorf("expected LOG_LEVEL to remain unmasked, got %q", c.NewValue)
		}
	}
}

func TestApply_RedactAll_MasksEverything(t *testing.T) {
	changes := sampleChanges()
	opts := Options{RedactAll: true}
	result := Apply(changes, opts)

	for _, c := range result {
		if c.NewValue != "" && !strings.Contains(c.NewValue, "*") {
			t.Errorf("expected all values masked with RedactAll, key=%s got %q", c.Key, c.NewValue)
		}
	}
}

func TestApply_CustomSensitiveKeys(t *testing.T) {
	changes := []differ.Change{
		{Type: differ.Added, Key: "STRIPE_PUBKEY", OldValue: "", NewValue: "pk_live_abc123"},
	}
	opts := Options{CustomSensitiveKeys: []string{"stripe"}}
	result := Apply(changes, opts)

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if !strings.Contains(result[0].NewValue, "*") {
		t.Errorf("expected STRIPE_PUBKEY to be masked via custom key, got %q", result[0].NewValue)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := []differ.Change{
		{Type: differ.Added, Key: "DB_PASSWORD", OldValue: "", NewValue: "plaintext"},
	}
	Apply(original, DefaultOptions())

	if original[0].NewValue != "plaintext" {
		t.Errorf("Apply must not mutate original changes, got %q", original[0].NewValue)
	}
}
