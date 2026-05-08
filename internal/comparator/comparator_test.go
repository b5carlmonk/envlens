package comparator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlens/internal/comparator"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestCompare_DetectsAddedKey(t *testing.T) {
	src := writeTempEnv(t, "APP=1\n")
	tgt := writeTempEnv(t, "APP=1\nNEW_KEY=hello\n")

	res, err := comparator.Compare(src, tgt, comparator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	if res.Changes[0].Key != "NEW_KEY" {
		t.Errorf("expected key NEW_KEY, got %s", res.Changes[0].Key)
	}
}

func TestCompare_DetectsRemovedKey(t *testing.T) {
	src := writeTempEnv(t, "APP=1\nOLD=gone\n")
	tgt := writeTempEnv(t, "APP=1\n")

	res, err := comparator.Compare(src, tgt, comparator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
}

func TestCompare_StrictMode_RemovedKeyReturnsError(t *testing.T) {
	src := writeTempEnv(t, "APP=1\nREQUIRED=yes\n")
	tgt := writeTempEnv(t, "APP=1\n")

	_, err := comparator.Compare(src, tgt, comparator.Options{StrictMode: true})
	if err == nil {
		t.Fatal("expected error in strict mode for removed key, got nil")
	}
}

func TestCompare_StrictMode_NoRemovals_Succeeds(t *testing.T) {
	src := writeTempEnv(t, "APP=1\n")
	tgt := writeTempEnv(t, "APP=2\nEXTRA=yes\n")

	res, err := comparator.Compare(src, tgt, comparator.Options{StrictMode: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCompare_InvalidSourceFile_ReturnsError(t *testing.T) {
	tgt := writeTempEnv(t, "APP=1\n")
	_, err := comparator.Compare("/nonexistent/.env", tgt, comparator.Options{})
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestCompare_PopulatesSourceAndTargetEnv(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\n")
	tgt := writeTempEnv(t, "FOO=baz\n")

	res, err := comparator.Compare(src, tgt, comparator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SourceEnv["FOO"] != "bar" {
		t.Errorf("expected source FOO=bar, got %s", res.SourceEnv["FOO"])
	}
	if res.TargetEnv["FOO"] != "baz" {
		t.Errorf("expected target FOO=baz, got %s", res.TargetEnv["FOO"])
	}
}
