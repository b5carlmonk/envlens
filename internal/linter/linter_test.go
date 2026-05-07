package linter

import (
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "API_KEY", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "REMOVED_VAR", Type: differ.Removed, OldValue: "gone"},
	}
}

func TestLint_NoIssues(t *testing.T) {
	result := Lint(sampleChanges())
	if result.HasIssues() {
		t.Errorf("expected no issues, got %d", len(result.Issues))
	}
}

func TestLint_EmptyValueWarning(t *testing.T) {
	changes := []differ.Change{
		{Key: "EMPTY_VAR", Type: differ.Added, NewValue: ""},
	}
	result := Lint(changes)
	if !result.HasIssues() {
		t.Fatal("expected issues for empty value")
	}
	if result.Issues[0].Key != "EMPTY_VAR" {
		t.Errorf("expected key EMPTY_VAR, got %s", result.Issues[0].Key)
	}
	if result.Issues[0].Level != LevelWarning {
		t.Errorf("expected warning level, got %s", result.Issues[0].Level)
	}
}

func TestLint_LowercaseKeyWarning(t *testing.T) {
	changes := []differ.Change{
		{Key: "bad_key", Type: differ.Added, NewValue: "value"},
	}
	result := Lint(changes)
	if !result.HasIssues() {
		t.Fatal("expected issue for lowercase key")
	}
	found := false
	for _, iss := range result.Issues {
		if iss.Key == "bad_key" && iss.Level == LevelWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for bad_key")
	}
}

func TestLint_WhitespaceInValue(t *testing.T) {
	changes := []differ.Change{
		{Key: "MY_VAR", Type: differ.Modified, OldValue: "a", NewValue: "hello world"},
	}
	result := Lint(changes)
	if !result.HasIssues() {
		t.Fatal("expected issue for whitespace in value")
	}
	if result.Issues[0].Message != "value contains unquoted whitespace" {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	changes := []differ.Change{
		{Key: "bad_key", Type: differ.Added, NewValue: ""},
	}
	result := Lint(changes)
	// Expect at least: empty value + lowercase key
	if len(result.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(result.Issues))
	}
}

func TestLint_RemovedKeyNoEmptyValueCheck(t *testing.T) {
	changes := []differ.Change{
		{Key: "GONE_VAR", Type: differ.Removed, OldValue: "val"},
	}
	result := Lint(changes)
	for _, iss := range result.Issues {
		if iss.Message == "value is empty" {
			t.Error("removed keys should not trigger empty value warning")
		}
	}
}
