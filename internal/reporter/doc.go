// Package reporter provides functionality for rendering environment diff
// results in multiple output formats.
//
// It consumes the output produced by the differ package and formats it
// for human-readable text or machine-parseable JSON.
//
// Supported formats:
//
//   - FormatText ("text") — coloured, line-oriented diff suitable for terminals.
//   - FormatJSON ("json") — structured JSON output for tooling integration.
//
// Basic usage:
//
//	report := reporter.NewReport("old.env", "new.env", changes)
//	err := report.Render(os.Stdout, reporter.FormatText)
//
// Use Report.HasChanges() to detect whether any diff was found, and
// Report.Summary() for a one-line change count.
package reporter
