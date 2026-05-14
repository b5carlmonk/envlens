package grouper

import (
	"sort"
	"strings"

	"github.com/envlens/internal/differ"
)

// Group represents a named collection of changes sharing a common prefix.
type Group struct {
	Name    string
	Changes []differ.Change
}

// Result holds all groups produced by a grouping operation.
type Result struct {
	Groups   []Group
	Ungrouped []differ.Change
}

// ByPrefix groups changes by their key prefix, splitting on the given
// separator (typically "_"). Keys with no separator land in Ungrouped.
func ByPrefix(changes []differ.Change, sep string) Result {
	if sep == "" {
		sep = "_"
	}

	buckets := make(map[string][]differ.Change)
	var ungrouped []differ.Change

	for _, c := range changes {
		idx := strings.Index(c.Key, sep)
		if idx <= 0 {
			ungrouped = append(ungrouped, c)
			continue
		}
		prefix := c.Key[:idx]
		buckets[prefix] = append(buckets[prefix], c)
	}

	// Sort group names for deterministic output.
	names := make([]string, 0, len(buckets))
	for name := range buckets {
		names = append(names, name)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, name := range names {
		groups = append(groups, Group{
			Name:    name,
			Changes: buckets[name],
		})
	}

	return Result{
		Groups:    groups,
		Ungrouped: ungrouped,
	}
}

// ByCustom groups changes using caller-supplied label→prefix mappings.
// A change is placed in the first matching group; unmatched changes go to Ungrouped.
func ByCustom(changes []differ.Change, labels map[string][]string) Result {
	buckets := make(map[string][]differ.Change)
	matched := make(map[string]bool)

	for _, c := range changes {
		placed := false
		for label, prefixes := range labels {
			for _, p := range prefixes {
				if strings.HasPrefix(strings.ToUpper(c.Key), strings.ToUpper(p)) {
					buckets[label] = append(buckets[label], c)
					matched[c.Key] = true
					placed = true
					break
				}
			}
			if placed {
				break
			}
		}
	}

	names := make([]string, 0, len(buckets))
	for name := range buckets {
		names = append(names, name)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, name := range names {
		groups = append(groups, Group{Name: name, Changes: buckets[name]})
	}

	var ungrouped []differ.Change
	for _, c := range changes {
		if !matched[c.Key] {
			ungrouped = append(ungrouped, c)
		}
	}

	return Result{Groups: groups, Ungrouped: ungrouped}
}
