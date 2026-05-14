package sorter

import (
	"sort"

	"github.com/yourusername/envlens/internal/differ"
)

// SortBy defines the field to sort changes by.
type SortBy string

const (
	// SortByKey sorts changes alphabetically by key name.
	SortByKey SortBy = "key"
	// SortByType sorts changes by change type: added, removed, modified, unchanged.
	SortByType SortBy = "type"
)

// Options configures how changes are sorted.
type Options struct {
	// By specifies the sort field. Defaults to SortByKey.
	By SortBy
	// Descending reverses the sort order when true.
	Descending bool
}

// typeOrder assigns a numeric rank to each change type for consistent ordering.
var typeOrder = map[differ.ChangeType]int{
	differ.Added:     0,
	differ.Removed:   1,
	differ.Modified:  2,
	differ.Unchanged: 3,
}

// Apply returns a sorted copy of the provided changes slice.
// The original slice is not modified.
func Apply(changes []differ.Change, opts Options) []differ.Change {
	if len(changes) == 0 {
		return changes
	}

	if opts.By == "" {
		opts.By = SortByKey
	}

	result := make([]differ.Change, len(changes))
	copy(result, changes)

	sort.SliceStable(result, func(i, j int) bool {
		var less bool

		switch opts.By {
		case SortByType:
			ri := typeOrder[result[i].Type]
			rj := typeOrder[result[j].Type]
			if ri != rj {
				less = ri < rj
				break
			}
			// Fall through to key sort for stable ordering within same type.
			fallthrough
		default:
			less = result[i].Key < result[j].Key
		}

		if opts.Descending {
			return !less
		}
		return less
	})

	return result
}
