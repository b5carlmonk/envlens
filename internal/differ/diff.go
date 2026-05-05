// Package differ provides functionality to compare two sets of environment
// variables and produce a structured diff result.
package differ

// ChangeType represents the kind of change detected for an environment variable.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Entry represents a single diff entry for an environment variable.
type Entry struct {
	Key      string
	OldValue string
	NewValue string
	Change   ChangeType
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// Added returns only the entries that were added.
func (r *Result) Added() []Entry {
	return r.filter(Added)
}

// Removed returns only the entries that were removed.
func (r *Result) Removed() []Entry {
	return r.filter(Removed)
}

// Modified returns only the entries that were modified.
func (r *Result) Modified() []Entry {
	return r.filter(Modified)
}

func (r *Result) filter(ct ChangeType) []Entry {
	var out []Entry
	for _, e := range r.Entries {
		if e.Change == ct {
			out = append(out, e)
		}
	}
	return out
}

// Compare takes two env maps (old, new) and returns a Result describing
// the differences between them.
func Compare(oldEnv, newEnv map[string]string) *Result {
	result := &Result{}

	// Check for removed or modified keys
	for key, oldVal := range oldEnv {
		if newVal, exists := newEnv[key]; !exists {
			result.Entries = append(result.Entries, Entry{
				Key:      key,
				OldValue: oldVal,
				Change:   Removed,
			})
		} else if oldVal != newVal {
			result.Entries = append(result.Entries, Entry{
				Key:      key,
				OldValue: oldVal,
				NewValue: newVal,
				Change:   Modified,
			})
		}
	}

	// Check for added keys
	for key, newVal := range newEnv {
		if _, exists := oldEnv[key]; !exists {
			result.Entries = append(result.Entries, Entry{
				Key:      key,
				NewValue: newVal,
				Change:   Added,
			})
		}
	}

	return result
}
