package normalizer_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/normalizer"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"db_host":   "  localhost  ",
		"DB_PORT":   "5432",
		"api_key":   "  secret  ",
		"EMPTY_VAL": "",
	}
}

func TestApply_TrimValues_RemovesWhitespace(t *testing.T) {
	env := sampleEnv()
	opts := normalizer.DefaultOptions()
	opts.TrimValues = true

	res := normalizer.Apply(env, opts)

	if got := res.Env["db_host"]; got != "localhost" {
		t.Errorf("expected 'localhost', got %q", got)
	}
	if got := res.Env["api_key"]; got != "secret" {
		t.Errorf("expected 'secret', got %q", got)
	}
}

func TestApply_TrimValues_TracksTrimmedKeys(t *testing.T) {
	env := sampleEnv()
	opts := normalizer.DefaultOptions()
	opts.TrimValues = true

	res := normalizer.Apply(env, opts)

	trimmedSet := make(map[string]bool)
	for _, k := range res.Trimmed {
		trimmedSet[k] = true
	}
	if !trimmedSet["db_host"] {
		t.Error("expected db_host to be in Trimmed list")
	}
	if !trimmedSet["api_key"] {
		t.Error("expected api_key to be in Trimmed list")
	}
	if trimmedSet["DB_PORT"] {
		t.Error("DB_PORT should not be in Trimmed list")
	}
}

func TestApply_UppercaseKeys_NormalizesKeys(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
		"Api_Key": "abc",
	}
	opts := normalizer.Options{UppercaseKeys: true}

	res := normalizer.Apply(env, opts)

	if _, ok := res.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key after uppercase normalization")
	}
	if _, ok := res.Env["API_KEY"]; !ok {
		t.Error("expected API_KEY key after uppercase normalization")
	}
	if res.Renamed["db_host"] != "DB_HOST" {
		t.Errorf("expected Renamed[db_host]=DB_HOST, got %q", res.Renamed["db_host"])
	}
}

func TestApply_RemoveEmpty_DropsEmptyValues(t *testing.T) {
	env := sampleEnv()
	opts := normalizer.Options{TrimValues: true, RemoveEmpty: true}

	res := normalizer.Apply(env, opts)

	if _, ok := res.Env["EMPTY_VAL"]; ok {
		t.Error("expected EMPTY_VAL to be removed")
	}
	if len(res.Removed) == 0 {
		t.Error("expected Removed to contain at least one key")
	}
}

func TestApply_NoOptions_ReturnsOriginal(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	opts := normalizer.Options{}

	res := normalizer.Apply(env, opts)

	if res.Env["KEY"] != "value" {
		t.Errorf("expected value to be unchanged, got %q", res.Env["KEY"])
	}
	if len(res.Trimmed) != 0 {
		t.Error("expected no trimmed keys")
	}
	if len(res.Removed) != 0 {
		t.Error("expected no removed keys")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"key": "  spaced  "}
	opts := normalizer.DefaultOptions()

	normalizer.Apply(env, opts)

	if env["key"] != "  spaced  " {
		t.Error("Apply must not mutate the original map")
	}
}
