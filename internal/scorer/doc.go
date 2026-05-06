// Package scorer provides risk assessment for environment variable diffs.
//
// It analyses a slice of differ.Change values and assigns a RiskLevel
// (none, low, medium, high) based on the nature of the changes:
//
//   - Added keys contribute a low base score.
//   - Removed keys contribute a medium score and are always noted as reasons.
//   - Modified keys contribute a higher score, with an additional penalty when
//     the key name matches a known sensitive pattern (e.g. PASSWORD, TOKEN,
//     SECRET, API_KEY).
//
// Usage:
//
//	changes := differ.Compare(before, after)
//	s := scorer.Evaluate(changes)
//	fmt.Printf("Risk: %s (%d pts)\n", s.Level, s.Points)
//	for _, r := range s.Reasons {
//		fmt.Println(" -", r)
//	}
package scorer
