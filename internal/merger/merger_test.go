package merger

import (
	"testing"
)

func TestMerge_SourceOnlyKeys(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	tgt := map[string]string{}
	r, err := Merge(src, tgt, Options{Strategy: StrategySourceWins})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", r.Merged["FOO"])
	}
}

func TestMerge_TargetOnlyKeys(t *testing.T) {
	src := map[string]string{}
	tgt := map[string]string{"BAZ": "qux"}
	r, err := Merge(src, tgt, Options{Strategy: StrategySourceWins})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", r.Merged["BAZ"])
	}
}

func TestMerge_ConflictSourceWins(t *testing.T) {
	src := map[string]string{"KEY": "from-source"}
	tgt := map[string]string{"KEY": "from-target"}
	r, err := Merge(src, tgt, Options{Strategy: StrategySourceWins})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["KEY"] != "from-source" {
		t.Errorf("expected from-source, got %q", r.Merged["KEY"])
	}
	if len(r.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(r.Conflicts))
	}
}

func TestMerge_ConflictTargetWins(t *testing.T) {
	src := map[string]string{"KEY": "from-source"}
	tgt := map[string]string{"KEY": "from-target"}
	r, err := Merge(src, tgt, Options{Strategy: StrategyTargetWins})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["KEY"] != "from-target" {
		t.Errorf("expected from-target, got %q", r.Merged["KEY"])
	}
	if r.Conflicts[0].Resolved != "from-target" {
		t.Errorf("resolved should be from-target")
	}
}

func TestMerge_ConflictStrategyError(t *testing.T) {
	src := map[string]string{"KEY": "a"}
	tgt := map[string]string{"KEY": "b"}
	_, err := Merge(src, tgt, Options{Strategy: StrategyError})
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_NoConflict_SameValue(t *testing.T) {
	src := map[string]string{"KEY": "same"}
	tgt := map[string]string{"KEY": "same"}
	r, err := Merge(src, tgt, Options{Strategy: StrategySourceWins})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(r.Conflicts))
	}
}
