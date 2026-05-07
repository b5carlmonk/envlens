package exporter_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/exporter"
)

func TestParseFormat_ValidFormats(t *testing.T) {
	cases := []struct {
		input    string
		expected exporter.Format
	}{
		{"text", exporter.FormatText},
		{"json", exporter.FormatJSON},
		{"markdown", exporter.FormatMarkdown},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := exporter.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestParseFormat_InvalidFormat_ReturnsError(t *testing.T) {
	_, err := exporter.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestSupportedFormats_ContainsAll(t *testing.T) {
	formats := exporter.SupportedFormats()
	expected := map[string]bool{"text": false, "json": false, "markdown": false}

	for _, f := range formats {
		if _, ok := expected[f]; ok {
			expected[f] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("expected format %q in SupportedFormats", name)
		}
	}
}
