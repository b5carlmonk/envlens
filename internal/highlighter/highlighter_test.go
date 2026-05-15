package highlighter

import (
	"strings"
	"testing"

	"github.com/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
		{Key: "API_KEY", Type: differ.Removed, OldValue: "secret"},
		{Key: "PORT", Type: differ.Modified, OldValue: "3000", NewValue: "8080"},
	}
}

func TestApply_ReturnsSameCount(t *testing.T) {
	changes := sampleChanges()
	lines := Apply(changes, DefaultOptions())
	if len(lines) != len(changes) {
		t.Fatalf("expected %d lines, got %d", len(changes), len(lines))
	}
}

func TestApply_AddedLine_ContainsPlusPrefix(t *testing.T) {
	changes := []differ.Change{
		{Key: "NEW_KEY", Type: differ.Added, NewValue: "value"},
	}
	lines := Apply(changes, Options{NoColor: true})
	if !strings.HasPrefix(lines[0].Text, "+ ") {
		t.Errorf("expected '+ ' prefix, got: %s", lines[0].Text)
	}
}

func TestApply_RemovedLine_ContainsMinusPrefix(t *testing.T) {
	changes := []differ.Change{
		{Key: "OLD_KEY", Type: differ.Removed, OldValue: "gone"},
	}
	lines := Apply(changes, Options{NoColor: true})
	if !strings.HasPrefix(lines[0].Text, "- ") {
		t.Errorf("expected '- ' prefix, got: %s", lines[0].Text)
	}
}

func TestApply_ModifiedLine_ContainsTildePrefix(t *testing.T) {
	changes := []differ.Change{
		{Key: "PORT", Type: differ.Modified, OldValue: "3000", NewValue: "8080"},
	}
	lines := Apply(changes, Options{NoColor: true})
	if !strings.HasPrefix(lines[0].Text, "~ ") {
		t.Errorf("expected '~ ' prefix, got: %s", lines[0].Text)
	}
}

func TestApply_ModifiedLine_ShowsOldValue(t *testing.T) {
	changes := []differ.Change{
		{Key: "PORT", Type: differ.Modified, OldValue: "3000", NewValue: "8080"},
	}
	lines := Apply(changes, Options{NoColor: true, ShowOldValue: true})
	if !strings.Contains(lines[0].Text, "was: 3000") {
		t.Errorf("expected old value in output, got: %s", lines[0].Text)
	}
}

func TestApply_ModifiedLine_HidesOldValue(t *testing.T) {
	changes := []differ.Change{
		{Key: "PORT", Type: differ.Modified, OldValue: "3000", NewValue: "8080"},
	}
	lines := Apply(changes, Options{NoColor: true, ShowOldValue: false})
	if strings.Contains(lines[0].Text, "was:") {
		t.Errorf("expected old value hidden, got: %s", lines[0].Text)
	}
}

func TestApply_NoColor_NoEscapeCodes(t *testing.T) {
	changes := sampleChanges()
	lines := Apply(changes, Options{NoColor: true})
	for _, l := range lines {
		if strings.Contains(l.Text, "\033[") {
			t.Errorf("expected no ANSI codes in NoColor mode, got: %s", l.Text)
		}
	}
}

func TestApply_WithColor_ContainsEscapeCodes(t *testing.T) {
	changes := sampleChanges()
	lines := Apply(changes, Options{NoColor: false})
	found := false
	for _, l := range lines {
		if strings.Contains(l.Text, "\033[") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected ANSI escape codes when color is enabled")
	}
}

func TestRender_JoinsWithNewlines(t *testing.T) {
	changes := sampleChanges()
	lines := Apply(changes, Options{NoColor: true})
	out := Render(lines)
	parts := strings.Split(out, "\n")
	if len(parts) != len(changes) {
		t.Errorf("expected %d newline-separated parts, got %d", len(changes), len(parts))
	}
}
