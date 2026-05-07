package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envlens/internal/differ"
)

// Format represents the output format for exporting diffs.
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// Options configures the export behavior.
type Options struct {
	Format     Format
	OutputPath string // empty means stdout
	MaskValues bool
}

// Export writes the given changes to the configured output in the specified format.
func Export(changes []differ.Change, opts Options) error {
	var content string
	var err error

	switch opts.Format {
	case FormatJSON:
		content, err = toJSON(changes)
	case FormatMarkdown:
		content = toMarkdown(changes)
	default:
		content = toText(changes)
	}

	if err != nil {
		return fmt.Errorf("exporter: failed to render %s: %w", opts.Format, err)
	}

	if opts.OutputPath == "" {
		fmt.Print(content)
		return nil
	}

	if err := os.WriteFile(opts.OutputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("exporter: failed to write file %s: %w", opts.OutputPath, err)
	}
	return nil
}

func toText(changes []differ.Change) string {
	var sb strings.Builder
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", c.Type, c.Key))
	}
	return sb.String()
}

func toJSON(changes []differ.Change) (string, error) {
	data, err := json.MarshalIndent(changes, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func toMarkdown(changes []differ.Change) string {
	var sb strings.Builder
	sb.WriteString("| Type | Key | Old Value | New Value |\n")
	sb.WriteString("|------|-----|-----------|-----------|\n")
	for _, c := range changes {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
			c.Type, c.Key, c.OldValue, c.NewValue))
	}
	return sb.String()
}
