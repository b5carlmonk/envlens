package tagger

import (
	"testing"

	"github.com/yourusername/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "JWT_SECRET", Type: differ.Modified, OldValue: "old", NewValue: "new"},
		{Key: "APP_PORT", Type: differ.Added, NewValue: "8080"},
		{Key: "FEATURE_DARK_MODE", Type: differ.Added, NewValue: "true"},
		{Key: "LOG_LEVEL", Type: differ.Modified, OldValue: "info", NewValue: "debug"},
		{Key: "SOME_RANDOM_KEY", Type: differ.Removed, OldValue: "value"},
	}
}

func TestApply_ReturnsSameCount(t *testing.T) {
	changes := sampleChanges()
	result := Apply(changes)
	if len(result.Tagged) != len(changes) {
		t.Errorf("expected %d tagged changes, got %d", len(changes), len(result.Tagged))
	}
}

func TestApply_DatabaseTag(t *testing.T) {
	changes := []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
	}
	result := Apply(changes)
	if !containsTag(result.Tagged[0].Tags, TagDatabase) {
		t.Errorf("expected tag %q for key DB_HOST, got %v", TagDatabase, result.Tagged[0].Tags)
	}
}

func TestApply_AuthTag(t *testing.T) {
	changes := []differ.Change{
		{Key: "JWT_SECRET", Type: differ.Modified, OldValue: "a", NewValue: "b"},
	}
	result := Apply(changes)
	tags := result.Tagged[0].Tags
	if !containsTag(tags, TagAuth) && !containsTag(tags, TagSecret) {
		t.Errorf("expected auth or secret tag for JWT_SECRET, got %v", tags)
	}
}

func TestApply_NetworkTag(t *testing.T) {
	changes := []differ.Change{
		{Key: "APP_PORT", Type: differ.Added, NewValue: "8080"},
	}
	result := Apply(changes)
	if !containsTag(result.Tagged[0].Tags, TagNetwork) {
		t.Errorf("expected tag %q for key APP_PORT, got %v", TagNetwork, result.Tagged[0].Tags)
	}
}

func TestApply_FeatureFlagTag(t *testing.T) {
	changes := []differ.Change{
		{Key: "FEATURE_DARK_MODE", Type: differ.Added, NewValue: "true"},
	}
	result := Apply(changes)
	if !containsTag(result.Tagged[0].Tags, TagFeatureFlag) {
		t.Errorf("expected tag %q for FEATURE_DARK_MODE, got %v", TagFeatureFlag, result.Tagged[0].Tags)
	}
}

func TestApply_UnknownTag_ForUnrecognizedKey(t *testing.T) {
	changes := []differ.Change{
		{Key: "SOME_RANDOM_KEY", Type: differ.Removed, OldValue: "value"},
	}
	result := Apply(changes)
	if !containsTag(result.Tagged[0].Tags, TagUnknown) {
		t.Errorf("expected tag %q for unrecognized key, got %v", TagUnknown, result.Tagged[0].Tags)
	}
}

func TestApply_LoggingTag(t *testing.T) {
	changes := []differ.Change{
		{Key: "LOG_LEVEL", Type: differ.Modified, OldValue: "info", NewValue: "debug"},
	}
	result := Apply(changes)
	if !containsTag(result.Tagged[0].Tags, TagLogging) {
		t.Errorf("expected tag %q for LOG_LEVEL, got %v", TagLogging, result.Tagged[0].Tags)
	}
}

func containsTag(tags []Tag, target Tag) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}
