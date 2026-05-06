// Package filter provides key-based and type-based filtering for environment
// variable diff results produced by the differ package.
//
// It allows callers to narrow down a list of [differ.Change] entries using
// one or more criteria:
//
//   - Prefix: only include keys that start with a given string (e.g. "DB_")
//   - KeyContains: only include keys that contain a given substring
//   - Types: only include changes of specific types ("added", "removed", "modified")
//
// Multiple criteria are combined with AND semantics — a change must satisfy
// every non-empty option to be included in the output.
//
// Example usage:
//
//	filtered := filter.Apply(changes, filter.Options{
//		Prefix: "DB_",
//		Types:  []string{"modified"},
//	})
package filter
