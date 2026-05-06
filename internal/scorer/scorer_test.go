package scorer_test

import (
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/scorer"
)

func TestEvaluate_NoChanges(t *testing.T) {
	s := scorer.Evaluate([]differ.Change{})
	if s.Level != scorer.RiskNone {
		t.Errorf("expected none, got %s", s.Level)
	}
	if s.Points != 0 {
		t.Errorf("expected 0 points, got %d", s.Points)
	}
}

func TestEvaluate_OnlyAdded(t *testing.T) {
	changes := []differ.Change{
		{Key: "NEW_VAR", Type: differ.Added, NewValue: "val"},
		{Key: "ANOTHER", Type: differ.Added, NewValue: "val2"},
	}
	s := scorer.Evaluate(changes)
	if s.Level != scorer.RiskLow {
		t.Errorf("expected low, got %s", s.Level)
	}
	if s.Points != 2 {
		t.Errorf("expected 2 points, got %d", s.Points)
	}
}

func TestEvaluate_RemovedKey_RaisesRisk(t *testing.T) {
	changes := []differ.Change{
		{Key: "OLD_VAR", Type: differ.Removed, OldValue: "val"},
	}
	s := scorer.Evaluate(changes)
	if s.Level != scorer.RiskLow {
		t.Errorf("expected low, got %s", s.Level)
	}
	if len(s.Reasons) == 0 {
		t.Error("expected at least one reason for removed key")
	}
}

func TestEvaluate_SensitiveModified_HighRisk(t *testing.T) {
	changes := []differ.Change{
		{Key: "DB_PASSWORD", Type: differ.Modified, OldValue: "old", NewValue: "new"},
	}
	s := scorer.Evaluate(changes)
	if s.Level != scorer.RiskHigh {
		t.Errorf("expected high, got %s", s.Level)
	}
	if len(s.Reasons) == 0 {
		t.Error("expected reason for sensitive key modification")
	}
}

func TestEvaluate_NonSensitiveModified_MediumRisk(t *testing.T) {
	changes := []differ.Change{
		{Key: "APP_PORT", Type: differ.Modified, OldValue: "8080", NewValue: "9090"},
		{Key: "LOG_LEVEL", Type: differ.Modified, OldValue: "info", NewValue: "debug"},
	}
	s := scorer.Evaluate(changes)
	if s.Level != scorer.RiskMedium {
		t.Errorf("expected medium, got %s", s.Level)
	}
}

func TestEvaluate_MixedChanges(t *testing.T) {
	changes := []differ.Change{
		{Key: "API_KEY", Type: differ.Modified, OldValue: "abc", NewValue: "xyz"},
		{Key: "REMOVED_VAR", Type: differ.Removed, OldValue: "val"},
		{Key: "NEW_VAR", Type: differ.Added, NewValue: "val"},
	}
	s := scorer.Evaluate(changes)
	if s.Level != scorer.RiskHigh {
		t.Errorf("expected high, got %s", s.Level)
	}
}
