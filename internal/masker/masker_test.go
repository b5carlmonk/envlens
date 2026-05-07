package masker_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/masker"
)

func TestIsSensitive_DetectsSecretKeys(t *testing.T) {
	sensitiveKeys := []string{
		"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY",
		"SECRET", "aws_secret_access_key", "CERT_DATA", "USER_PWD",
	}
	for _, k := range sensitiveKeys {
		if !masker.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_AllowsNonSensitiveKeys(t *testing.T) {
	safeKeys := []string{"APP_ENV", "PORT", "LOG_LEVEL", "REGION", "TIMEOUT"}
	for _, k := range safeKeys {
		if masker.IsSensitive(k) {
			t.Errorf("expected %q to be non-sensitive", k)
		}
	}
}

func TestMask_ReplacesWithStars(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.Mask("supersecret", opts)
	if !strings.HasPrefix(result, "******") {
		t.Errorf("expected result to start with '******', got %q", result)
	}
}

func TestMask_PreservesVisibleSuffix(t *testing.T) {
	opts := masker.DefaultOptions() // VisibleSuffix = 4
	result := masker.Mask("supersecret", opts)
	if !strings.HasSuffix(result, "cret") {
		t.Errorf("expected suffix 'cret', got %q", result)
	}
}

func TestMask_ShortValue_FullyMasked(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.Mask("abc", opts)
	if strings.Contains(result, "abc") {
		t.Errorf("expected short value to be fully masked, got %q", result)
	}
}

func TestMask_ZeroVisibleSuffix_FullyMasked(t *testing.T) {
	opts := masker.Options{MaskChar: "*", VisibleSuffix: 0}
	result := masker.Mask("topsecret", opts)
	if strings.Contains(result, "topsecret") {
		t.Errorf("expected value to be fully masked, got %q", result)
	}
}

func TestMaskIfSensitive_MasksSensitiveKey(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.MaskIfSensitive("DB_PASSWORD", "hunter2", opts)
	if result == "hunter2" {
		t.Error("expected sensitive value to be masked")
	}
}

func TestMaskIfSensitive_LeavesNonSensitiveIntact(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.MaskIfSensitive("APP_ENV", "production", opts)
	if result != "production" {
		t.Errorf("expected non-sensitive value unchanged, got %q", result)
	}
}
