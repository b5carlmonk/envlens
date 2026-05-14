// Package differ provides functionality for comparing two environment variable
// maps and producing a structured list of changes.
//
// A change can be one of three types:
//   - Added: a key present in the target but not in the source
//   - Removed: a key present in the source but not in the target
//   - Modified: a key present in both but with a different value
//
// Usage:
//
//	changes := differ.Compare(sourceEnv, targetEnv)
//	for _, c := range changes {
//		fmt.Println(c.Key, c.Type)
//	}
package differ
