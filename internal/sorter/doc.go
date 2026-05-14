// Package sorter provides utilities for ordering environment variable changes
// returned by the differ package.
//
// Changes can be sorted by key name (alphabetically) or by change type
// (added → removed → modified → unchanged), with optional descending order.
//
// Example usage:
//
//	changes := differ.Compare(source, target)
//	sorted := sorter.Apply(changes, sorter.Options{
//		By:         sorter.SortByType,
//		Descending: false,
//	})
//
// The Apply function always returns a new slice and never mutates the input,
// making it safe to use in pipelines alongside filter, redactor, or linter.
package sorter
