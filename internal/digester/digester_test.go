package digester_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/digester"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PASS":  "secret",
		"LOG_LEVEL": "info",
	}
}

func TestCompute_ReturnsDeterministicDigest(t *testing.T) {
	env := sampleEnv()
	r1 := digester.Compute(env)
	r2 := digester.Compute(env)
	if r1.Digest != r2.Digest {
		t.Errorf("expected same digest on repeated calls, got %q and %q", r1.Digest, r2.Digest)
	}
}

func TestCompute_KeyCountMatchesEnv(t *testing.T) {
	env := sampleEnv()
	r := digester.Compute(env)
	if r.KeyCount != len(env) {
		t.Errorf("expected KeyCount=%d, got %d", len(env), r.KeyCount)
	}
}

func TestCompute_KeyDigestsPresent(t *testing.T) {
	env := sampleEnv()
	r := digester.Compute(env)
	for k := range env {
		if _, ok := r.KeyDigests[k]; !ok {
			t.Errorf("expected key digest for %q", k)
		}
	}
}

func TestCompute_DigestChangesWhenValueChanges(t *testing.T) {
	env1 := sampleEnv()
	env2 := sampleEnv()
	env2["DB_PASS"] = "changed"

	r1 := digester.Compute(env1)
	r2 := digester.Compute(env2)
	if r1.Digest == r2.Digest {
		t.Error("expected digest to differ after value change")
	}
}

func TestEqual_SameEnv_ReturnsTrue(t *testing.T) {
	env := sampleEnv()
	if !digester.Equal(digester.Compute(env), digester.Compute(env)) {
		t.Error("expected Equal to return true for identical envs")
	}
}

func TestEqual_DifferentEnv_ReturnsFalse(t *testing.T) {
	env1 := sampleEnv()
	env2 := sampleEnv()
	env2["NEW_KEY"] = "value"
	if digester.Equal(digester.Compute(env1), digester.Compute(env2)) {
		t.Error("expected Equal to return false for different envs")
	}
}

func TestDiffDigests_ReturnsChangedKeys(t *testing.T) {
	env1 := sampleEnv()
	env2 := sampleEnv()
	env2["DB_PASS"] = "newpass"
	env2["EXTRA_KEY"] = "extra"

	changed := digester.DiffDigests(digester.Compute(env1), digester.Compute(env2))

	found := map[string]bool{}
	for _, k := range changed {
		found[k] = true
	}
	if !found["DB_PASS"] {
		t.Error("expected DB_PASS in diff")
	}
	if !found["EXTRA_KEY"] {
		t.Error("expected EXTRA_KEY in diff")
	}
}

func TestDiffDigests_NoChanges_ReturnsEmpty(t *testing.T) {
	env := sampleEnv()
	changed := digester.DiffDigests(digester.Compute(env), digester.Compute(env))
	if len(changed) != 0 {
		t.Errorf("expected no diff, got %v", changed)
	}
}
