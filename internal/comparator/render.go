package comparator

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// RenderText writes a human-readable comparison summary to w.
func RenderText(w io.Writer, r *Result) error {
	fmt.Fprintf(w, "Comparing env files\n")
	fmt.Fprintf(w, "  Source : %s\n", r.SourceFile)
	fmt.Fprintf(w, "  Target : %s\n", r.TargetFile)
	fmt.Fprintf(w, "  Changes: %d\n", len(r.Changes))

	if len(r.Changes) == 0 {
		fmt.Fprintln(w, "  No differences found.")
		return nil
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, c := range r.Changes {
		symbol := changeSymbol(c.Type)
		switch c.Type {
		case differ.Added:
			fmt.Fprintf(w, "  %s %-30s = %s\n", symbol, c.Key, c.NewValue)
		case differ.Removed:
			fmt.Fprintf(w, "  %s %-30s (was: %s)\n", symbol, c.Key, c.OldValue)
		case differ.Modified:
			fmt.Fprintf(w, "  %s %-30s %s -> %s\n", symbol, c.Key, c.OldValue, c.NewValue)
		}
	}
	return nil
}

// RenderJSON writes a JSON-encoded comparison result to w.
func RenderJSON(w io.Writer, r *Result) error {
	type jsonResult struct {
		SourceFile string         `json:"source_file"`
		TargetFile string         `json:"target_file"`
		ChangeCount int           `json:"change_count"`
		Changes    []differ.Change `json:"changes"`
	}
	out := jsonResult{
		SourceFile:  r.SourceFile,
		TargetFile:  r.TargetFile,
		ChangeCount: len(r.Changes),
		Changes:     r.Changes,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func changeSymbol(t differ.ChangeType) string {
	switch t {
	case differ.Added:
		return "+"
	case differ.Removed:
		return "-"
	case differ.Modified:
		return "~"
	default:
		return "?"
	}
}
