// Package inspector enriches environment variable changes with deep
// metadata, including value length deltas, inferred type hints
// (string, numeric, boolean, URL), and value-equality checks.
//
// It is intended to be used after diffing with the differ package to
// provide richer context for reporting, auditing, and scoring pipelines.
//
// Example usage:
//
//	changes := differ.Compare(source, target)
//	inspections := inspector.Inspect(changes)
//	for _, ins := range inspections {
//		fmt.Printf("%s: %s -> %s (delta: %d)\n",
//			ins.Change.Key, ins.OldHint, ins.NewHint, ins.LengthDelta)
//	}
package inspector
