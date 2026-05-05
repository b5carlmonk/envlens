package differ_test

import (
	"testing"

	"github.com/user/envlens/internal/differ"
)

func TestCompare_AddedKeys(t *testing.T) {
	oldEnv := map[string]string{"FOO": "bar"}
	newEnv := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := differ.Compare(oldEnv, newEnv)
	added := result.Added()

	if len(added) != 1 {
		t.Fatalf("expected 1 added entry, got %d", len(added))
	}
	if added[0].Key != "BAZ" || added[0].NewValue != "qux" {
		t.Errorf("unexpected added entry: %+v", added[0])
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	oldEnv := map[string]string{"FOO": "bar", "GONE": "bye"}
	newEnv := map[string]string{"FOO": "bar"}

	result := differ.Compare(oldEnv, newEnv)
	removed := result.Removed()

	if len(removed) != 1 {
		t.Fatalf("expected 1 removed entry, got %d", len(removed))
	}
	if removed[0].Key != "GONE" || removed[0].OldValue != "bye" {
		t.Errorf("unexpected removed entry: %+v", removed[0])
	}
}

func TestCompare_ModifiedKeys(t *testing.T) {
	oldEnv := map[string]string{"DB_URL": "localhost:5432"}
	newEnv := map[string]string{"DB_URL": "prod-db:5432"}

	result := differ.Compare(oldEnv, newEnv)
	modified := result.Modified()

	if len(modified) != 1 {
		t.Fatalf("expected 1 modified entry, got %d", len(modified))
	}
	e := modified[0]
	if e.Key != "DB_URL" || e.OldValue != "localhost:5432" || e.NewValue != "prod-db:5432" {
		t.Errorf("unexpected modified entry: %+v", e)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	result := differ.Compare(env, env)

	if len(result.Entries) != 0 {
		t.Errorf("expected no diff entries, got %d", len(result.Entries))
	}
}

func TestCompare_EmptyEnvs(t *testing.T) {
	result := differ.Compare(map[string]string{}, map[string]string{})
	if len(result.Entries) != 0 {
		t.Errorf("expected empty result for empty envs, got %d entries", len(result.Entries))
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	oldEnv := map[string]string{"KEEP": "same", "CHANGE": "old", "DROP": "gone"}
	newEnv := map[string]string{"KEEP": "same", "CHANGE": "new", "NEW": "here"}

	result := differ.Compare(oldEnv, newEnv)

	if len(result.Added()) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Added()))
	}
	if len(result.Removed()) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed()))
	}
	if len(result.Modified()) != 1 {
		t.Errorf("expected 1 modified, got %d", len(result.Modified()))
	}
}
