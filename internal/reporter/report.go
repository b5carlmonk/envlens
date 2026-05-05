package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the diff result and metadata for rendering.
type Report struct {
	FromFile string
	ToFile   string
	Diff     []differ.Change
}

// NewReport creates a new Report from two file paths and a diff result.
func NewReport(fromFile, toFile string, diff []differ.Change) *Report {
	return &Report{
		FromFile: fromFile,
		ToFile:   toFile,
		Diff:     diff,
	}
}

// HasChanges returns true if the report contains any changes.
func (r *Report) HasChanges() bool {
	return len(r.Diff) > 0
}

// Summary returns a short human-readable summary of the diff.
func (r *Report) Summary() string {
	var added, removed, modified int
	for _, c := range r.Diff {
		switch c.Type {
		case differ.Added:
			added++
		case differ.Removed:
			removed++
		case differ.Modified:
			modified++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d modified", added, removed, modified)
}

// Render writes the report to w in the given format.
func (r *Report) Render(w io.Writer, format Format) error {
	switch format {
	case FormatJSON:
		return renderJSON(w, r)
	case FormatText:
		return renderText(w, r)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func changeSymbol(ct differ.ChangeType) string {
	switch ct {
	case differ.Added:
		return "+"
	case differ.Removed:
		return "-"
	case differ.Modified:
		return "~"
	}
	return " "
}

func renderText(w io.Writer, r *Report) error {
	fmt.Fprintf(w, "envlens diff: %s → %s\n", r.FromFile, r.ToFile)
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 40))
	if !r.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return nil
	}
	for _, c := range r.Diff {
		sym := changeSymbol(c.Type)
		switch c.Type {
		case differ.Added:
			fmt.Fprintf(w, "%s %s=%s\n", sym, c.Key, c.NewValue)
		case differ.Removed:
			fmt.Fprintf(w, "%s %s=%s\n", sym, c.Key, c.OldValue)
		case differ.Modified:
			fmt.Fprintf(w, "%s %s: %q → %q\n", sym, c.Key, c.OldValue, c.NewValue)
		}
	}
	fmt.Fprintf(w, "\n%s\n", r.Summary())
	return nil
}
