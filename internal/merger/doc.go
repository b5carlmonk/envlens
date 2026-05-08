// Package merger provides functionality for combining two environment variable
// maps (source and target) into a single merged result.
//
// When a key exists in both maps with different values, a conflict is detected.
// The caller controls conflict resolution via a Strategy:
//
//   - StrategySourceWins — the source value takes precedence.
//   - StrategyTargetWins — the target value takes precedence.
//   - StrategyError      — an error is returned immediately on the first conflict.
//
// The Result type captures the final merged map together with a list of all
// conflicts that were encountered, including which value was ultimately chosen.
//
// Rendering helpers (RenderText, RenderJSON) are provided for reporting merge
// outcomes in human-readable or machine-readable formats.
package merger
