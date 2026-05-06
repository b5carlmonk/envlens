package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Added, Key: "NEW_KEY", OldValue: "", NewValue: "new_value"},
		{Type: differ.Removed, Key: "OLD_KEY", OldValue: "old_value", NewValue: ""},
		{Type: differ.Modified, Key: "CHANGED_KEY", OldValue: "before", NewValue: "after"},
	}
}

func TestRenderTable_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTable(&buf, sampleChanges()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"TYPE", "KEY", "OLD VALUE", "NEW VALUE"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected header column %q in output", col)
		}
	}
}

func TestRenderTable_ContainsChangeTypes(t *testing.T) {
	var buf bytes.Buffer
	_ = RenderTable(&buf, sampleChanges())
	out := buf.String()

	for _, typ := range []string{"ADDED", "REMOVED", "MODIFIED"} {
		if !strings.Contains(out, typ) {
			t.Errorf("expected change type %q in output", typ)
		}
	}
}

func TestRenderTable_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	_ = RenderTable(&buf, sampleChanges())
	out := buf.String()

	for _, key := range []string{"NEW_KEY", "OLD_KEY", "CHANGED_KEY"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output", key)
		}
	}
}

func TestRenderTable_EmptyChanges(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTable(&buf, []differ.Change{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TYPE") {
		t.Error("expected header even with empty changes")
	}
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10c", 10, "exactly10c"},
		{"this is a very long string", 10, "this is..."},
		{"abc", 3, "abc"},
		{"abcd", 2, "ab"},
	}
	for _, tc := range cases {
		got := truncate(tc.input, tc.maxLen)
		if got != tc.expected {
			t.Errorf("truncate(%q, %d) = %q; want %q", tc.input, tc.maxLen, got, tc.expected)
		}
	}
}
