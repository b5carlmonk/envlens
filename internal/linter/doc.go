// Package linter provides static analysis rules for environment variable changes.
//
// It inspects a slice of [differ.Change] values and applies a set of built-in
// rules to surface potential configuration issues before deployment.
//
// # Built-in Rules
//
// The following rules are applied automatically by [Lint]:
//
//   - Empty value: warns when an added or modified key has an empty value.
//   - UPPER_SNAKE_CASE: warns when a key contains lowercase letters.
//   - Whitespace in value: warns when a value contains unquoted spaces or tabs.
//
// # Usage
//
//	result := linter.Lint(changes)
//	if result.HasIssues() {
//		fmt.Print(linter.RenderText(result))
//	}
//
// Results can also be serialised to JSON via [RenderJSON].
package linter
