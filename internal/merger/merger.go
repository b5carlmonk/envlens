package merger

import (
	"fmt"
	"maps"
)

// Strategy defines how conflicts are resolved when merging two env maps.
type Strategy int

const (
	// StrategySourceWins keeps the source value on conflict.
	StrategySourceWins Strategy = iota
	// StrategyTargetWins keeps the target value on conflict.
	StrategyTargetWins
	// StrategyError returns an error on conflict.
	StrategyError
)

// Result holds the merged environment map and metadata about the merge.
type Result struct {
	Merged    map[string]string
	Conflicts []Conflict
}

// Conflict describes a key that existed in both source and target with different values.
type Conflict struct {
	Key         string
	SourceValue string
	TargetValue string
	Resolved    string
}

// Options controls merger behaviour.
type Options struct {
	Strategy Strategy
}

// Merge combines source and target env maps according to the given options.
// Keys present only in source or only in target are included as-is.
// Keys present in both are handled according to the Strategy.
func Merge(source, target map[string]string, opts Options) (Result, error) {
	merged := make(map[string]string)
	maps.Copy(merged, target)

	var conflicts []Conflict

	for k, sv := range source {
		tv, exists := merged[k]
		if !exists {
			merged[k] = sv
			continue
		}
		if sv == tv {
			continue
		}
		// conflict
		switch opts.Strategy {
		case StrategySourceWins:
			merged[k] = sv
			conflicts = append(conflicts, Conflict{Key: k, SourceValue: sv, TargetValue: tv, Resolved: sv})
		case StrategyTargetWins:
			merged[k] = tv
			conflicts = append(conflicts, Conflict{Key: k, SourceValue: sv, TargetValue: tv, Resolved: tv})
		case StrategyError:
			return Result{}, fmt.Errorf("merge conflict on key %q: source=%q target=%q", k, sv, tv)
		}
	}

	return Result{Merged: merged, Conflicts: conflicts}, nil
}
