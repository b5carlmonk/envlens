package exporter

import "fmt"

// ParseFormat converts a string to a Format constant.
// Returns an error if the format is unrecognized.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText:
		return FormatText, nil
	case FormatJSON:
		return FormatJSON, nil
	case FormatMarkdown:
		return FormatMarkdown, nil
	default:
		return "", fmt.Errorf("exporter: unknown format %q; valid options are: text, json, markdown", s)
	}
}

// SupportedFormats returns a slice of all supported format strings.
func SupportedFormats() []string {
	return []string{
		string(FormatText),
		string(FormatJSON),
		string(FormatMarkdown),
	}
}
