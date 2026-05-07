// Package exporter provides functionality to export environment variable diffs
// to various output formats and destinations.
//
// Supported formats:
//
//   - text     — plain human-readable lines, one change per line
//   - json     — structured JSON array of change objects
//   - markdown — GitHub-flavored markdown table
//
// Output can be directed to stdout (default) or written to a file by setting
// the OutputPath field in Options.
//
// Example usage:
//
//	changes := differ.Compare(source, target)
//	err := exporter.Export(changes, exporter.Options{
//		Format:     exporter.FormatMarkdown,
//		OutputPath: "diff-report.md",
//	})
package exporter
