package validator_test

import (
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/validator"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "APP_PORT", Type: differ.Added, NewValue: "8080"},
		{Key: "DB_HOST", Type: differ.Modified, OldValue: "localhost", NewValue: "prod.db"},
		{Key: "SECRET_KEY", Type: differ.Removed, OldValue: "abc123"},
	}
}

func TestValidate_NoIssues(t *testing.T) {
	result := validator.Validate(sampleChanges())
	if result.HasIssues() {
		t.Errorf("expected no issues, got %d: %+v", len(result.Issues), result.Issues)
	}
}

func TestValidate_NonConventionalKey(t *testing.T) {
	changes := []differ.Change{
		{Key: "my-lower-key", Type: differ.Added, NewValue: "value"},
	}
	result := validator.Validate(changes)
	if !result.HasIssues() {
		t.Fatal("expected issues for non-conventional key name")
	}
	if result.Issues[0].Severity != "warn" {
		t.Errorf("expected severity 'warn', got %q", result.Issues[0].Severity)
	}
}

func TestValidate_EmptyValueOnAdded(t *testing.T) {
	changes := []differ.Change{
		{Key: "APP_NAME", Type: differ.Added, NewValue: ""},
	}
	result := validator.Validate(changes)
	if !result.HasIssues() {
		t.Fatal("expected issue for empty added value")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "APP_NAME" && issue.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn issue for APP_NAME with empty value")
	}
}

func TestValidate_EmptyValueOnModified(t *testing.T) {
	changes := []differ.Change{
		{Key: "DB_URL", Type: differ.Modified, OldValue: "old", NewValue: "   "},
	}
	result := validator.Validate(changes)
	if !result.HasIssues() {
		t.Fatal("expected issue for whitespace-only modified value")
	}
}

func TestValidate_Summary_WithIssues(t *testing.T) {
	changes := []differ.Change{
		{Key: "bad-key", Type: differ.Added, NewValue: ""},
	}
	result := validator.Validate(changes)
	summary := result.Summary()
	if summary == "No validation issues found." {
		t.Error("expected non-empty summary when issues exist")
	}
}

func TestValidate_Summary_NoIssues(t *testing.T) {
	result := validator.Validate(sampleChanges())
	if result.Summary() != "No validation issues found." {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestValidate_BlankKey_ReturnsError(t *testing.T) {
	changes := []differ.Change{
		{Key: "", Type: differ.Added, NewValue: "something"},
	}
	result := validator.Validate(changes)
	if !result.HasIssues() {
		t.Fatal("expected error issue for blank key")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one error-severity issue for blank key")
	}
}
