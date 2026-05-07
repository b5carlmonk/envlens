package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/differ"
	"github.com/yourusername/envlens/internal/exporter"
)

func sampleChanges() []differ.Change {
	return []differ.Change{
		{Type: differ.Added, Key: "NEW_KEY", OldValue: "", NewValue: "hello"},
		{Type: differ.Removed, Key: "OLD_KEY", OldValue: "bye", NewValue: ""},
		{Type: differ.Modified, Key: "MOD_KEY", OldValue: "v1", NewValue: "v2"},
	}
}

func TestExport_TextToStdout(t *testing.T) {
	changes := sampleChanges()
	opts := exporter.Options{Format: exporter.FormatText}
	if err := exporter.Export(changes, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExport_JSONToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	changes := sampleChanges()

	opts := exporter.Options{Format: exporter.FormatJSON, OutputPath: path}
	if err := exporter.Export(changes, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var result []differ.Change
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 changes, got %d", len(result))
	}
}

func TestExport_MarkdownToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.md")
	changes := sampleChanges()

	opts := exporter.Options{Format: exporter.FormatMarkdown, OutputPath: path}
	if err := exporter.Export(changes, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "| Type |") {
		t.Error("expected markdown table header")
	}
	if !strings.Contains(content, "NEW_KEY") {
		t.Error("expected key NEW_KEY in markdown output")
	}
}

func TestExport_InvalidPath_ReturnsError(t *testing.T) {
	changes := sampleChanges()
	opts := exporter.Options{Format: exporter.FormatText, OutputPath: "/nonexistent/dir/out.txt"}
	if err := exporter.Export(changes, opts); err == nil {
		t.Error("expected error for invalid output path")
	}
}
