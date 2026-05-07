package profiler

import (
	"testing"

	"github.com/your-org/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "JWT_SECRET", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "PORT", Type: differ.Added, NewValue: "8080"},
		{Key: "REDIS_URL", Type: differ.Added, NewValue: "redis://localhost"},
		{Key: "LOG_LEVEL", Type: differ.Modified, OldValue: "info", NewValue: "debug"},
		{Key: "CUSTOM_VAR", Type: differ.Added, NewValue: "value"},
	}
}

func TestRun_DefaultProfile(t *testing.T) {
	result, err := Run(sampleChanges(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Profile.Name != "web" {
		t.Errorf("expected profile 'web', got %q", result.Profile.Name)
	}
}

func TestRun_MatchesDatabaseCategory(t *testing.T) {
	result, err := Run(sampleChanges(), "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dbKeys, ok := result.Matched["database"]
	if !ok {
		t.Fatal("expected 'database' category in matched")
	}
	if len(dbKeys) == 0 || dbKeys[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in database category, got %v", dbKeys)
	}
}

func TestRun_MatchesAuthCategory(t *testing.T) {
	result, err := Run(sampleChanges(), "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	authKeys := result.Matched["auth"]
	if len(authKeys) == 0 {
		t.Error("expected at least one key in auth category")
	}
}

func TestRun_UnmatchedKeys(t *testing.T) {
	result, err := Run(sampleChanges(), "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, k := range result.Unmatched {
		if k == "CUSTOM_VAR" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected CUSTOM_VAR in unmatched, got %v", result.Unmatched)
	}
}

func TestRun_TotalAndMatchedCounts(t *testing.T) {
	result, err := Run(sampleChanges(), "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalKeys != 6 {
		t.Errorf("expected TotalKeys=6, got %d", result.TotalKeys)
	}
	if result.MatchedKeys+len(result.Unmatched) != result.TotalKeys {
		t.Errorf("matched + unmatched should equal total")
	}
}

func TestRun_UnknownProfile_ReturnsError(t *testing.T) {
	_, err := Run(sampleChanges(), "nonexistent")
	if err == nil {
		t.Error("expected error for unknown profile, got nil")
	}
}

func TestRun_CloudProfile(t *testing.T) {
	changes := []differ.Change{
		{Key: "AWS_ACCESS_KEY_ID", Type: differ.Added, NewValue: "AKIA..."},
		{Key: "GCP_PROJECT", Type: differ.Added, NewValue: "my-project"},
		{Key: "UNRELATED", Type: differ.Added, NewValue: "val"},
	}
	result, err := Run(changes, "cloud")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.Matched["aws"]; !ok {
		t.Error("expected 'aws' category to be matched")
	}
	if _, ok := result.Matched["gcp"]; !ok {
		t.Error("expected 'gcp' category to be matched")
	}
}
