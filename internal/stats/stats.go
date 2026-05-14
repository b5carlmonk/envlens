// Package stats computes aggregate statistics from a set of environment
// variable changes produced by the differ package.
package stats

import (
	"fmt"
	"sort"

	"github.com/yourusername/envlens/internal/differ"
)

// Result holds computed statistics for a diff result.
type Result struct {
	Total    int
	Added    int
	Removed  int
	Modified int
	// TopPrefixes lists the top N most-changed key prefixes (by underscore segment).
	TopPrefixes []PrefixCount
}

// PrefixCount holds a prefix and the number of changes associated with it.
type PrefixCount struct {
	Prefix string
	Count  int
}

// Compute derives a Result from the provided slice of changes.
func Compute(changes []differ.Change) Result {
	prefixMap := make(map[string]int)

	r := Result{Total: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case differ.Added:
			r.Added++
		case differ.Removed:
			r.Removed++
		case differ.Modified:
			r.Modified++
		}
		if prefix := keyPrefix(c.Key); prefix != "" {
			prefixMap[prefix]++
		}
	}
	r.TopPrefixes = topPrefixes(prefixMap, 5)
	return r
}

// keyPrefix returns the first underscore-delimited segment of a key,
// or an empty string if no underscore is present.
func keyPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return ""
}

// topPrefixes returns up to n prefixes sorted by descending count.
func topPrefixes(m map[string]int, n int) []PrefixCount {
	list := make([]PrefixCount, 0, len(m))
	for k, v := range m {
		list = append(list, PrefixCount{Prefix: k, Count: v})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Prefix < list[j].Prefix
	})
	if len(list) > n {
		list = list[:n]
	}
	return list
}

// Summary returns a human-readable one-line summary of the result.
func Summary(r Result) string {
	return fmt.Sprintf("total=%d added=%d removed=%d modified=%d",
		r.Total, r.Added, r.Removed, r.Modified)
}
