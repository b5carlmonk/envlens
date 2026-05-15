// Package scorer evaluates the risk level of a set of environment variable
// changes produced by the differ package.
//
// # Overview
//
// The scorer assigns a numeric risk score to a diff result based on heuristics
// such as whether sensitive keys were modified or removed, whether values were
// emptied, and how many total changes occurred.
//
// # Risk Levels
//
// Scores map to one of four risk levels:
//
//   - Low    – minor or additive changes only
//   - Medium – non-sensitive modifications present
//   - High   – sensitive keys modified or non-trivial removals detected
//   - Critical – sensitive keys removed or multiple high-risk signals combined
//
// # Usage
//
//	result := scorer.Evaluate(changes)
//	fmt.Println(scorer.RenderText(result))
package scorer
