// Package differ implements environment variable diffing for envlens.
//
// It compares two maps of environment variables (old vs new) and produces
// a structured Result containing categorised change entries.
//
// Usage:
//
//	oldEnv, _ := parser.ParseFile(".env.old")
//	newEnv, _ := parser.ParseFile(".env.new")
//
//	result := differ.Compare(oldEnv, newEnv)
//
//	for _, e := range result.Added() {
//		fmt.Printf("+ %s=%s\n", e.Key, e.NewValue)
//	}
//	for _, e := range result.Removed() {
//		fmt.Printf("- %s=%s\n", e.Key, e.OldValue)
//	}
//	for _, e := range result.Modified() {
//		fmt.Printf("~ %s: %s -> %s\n", e.Key, e.OldValue, e.NewValue)
//	}
//
// ChangeType constants: Added, Removed, Modified, Unchanged.
package differ
