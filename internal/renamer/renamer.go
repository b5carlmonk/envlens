// Package renamer provides utilities for detecting and applying key rename
// operations between two environment variable maps. It identifies cases where
// a key was removed in one env and a matching value appears under a new key
// in the other, suggesting a rename rather than an independent add/remove.
package renamer

import "github.com/user/envlens/internal/differ"

// RenameHint represents a suspected key rename between two environments.
type RenameHint struct {
	OldKey string
	NewKey string
	Value  string
}

// Result holds the output of a rename detection pass.
type Result struct {
	Hints     []RenameHint
	Remaining []differ.Change
}

// Detect analyses a slice of differ.Change values and attempts to pair removed
// keys with added keys that share the same value, treating them as renames.
// Changes that cannot be paired are returned in Result.Remaining.
func Detect(changes []differ.Change) Result {
	var added, removed []differ.Change
	var other []differ.Change

	for _, c := range changes {
		switch c.Type {
		case differ.Added:
			added = append(added, c)
		case differ.Removed:
			removed = append(removed, c)
		default:
			other = append(other, c)
		}
	}

	matched := make(map[string]bool)
	var hints []RenameHint

	for _, r := range removed {
		for _, a := range added {
			if matched[a.Key] {
				continue
			}
			if r.OldValue == a.NewValue && r.OldValue != "" {
				hints = append(hints, RenameHint{
					OldKey: r.Key,
					NewKey: a.Key,
					Value:  r.OldValue,
				})
				matched[a.Key] = true
				matched[r.Key] = true
				break
			}
		}
	}

	var remaining []differ.Change
	for _, c := range append(added, removed...) {
		if !matched[c.Key] {
			remaining = append(remaining, c)
		}
	}
	remaining = append(remaining, other...)

	return Result{Hints: hints, Remaining: remaining}
}
