// Package grouper organises a slice of differ.Change entries into named
// groups based on key structure.
//
// Two strategies are provided:
//
//   - ByPrefix splits keys on a separator character (default "_") and uses
//     the leading segment as the group name. For example DB_HOST and DB_PORT
//     both land in the "DB" group. Keys that contain no separator are
//     collected in Result.Ungrouped.
//
//   - ByCustom accepts a caller-supplied map of label → prefixes so that
//     domain-meaningful names ("database", "auth") can be assigned
//     regardless of the raw key prefix. A change is placed in the first
//     matching group; unmatched changes are collected in Result.Ungrouped.
//
// Groups within a Result are always returned in alphabetical order for
// deterministic output and reproducible tests.
package grouper
