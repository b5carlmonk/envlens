// Package snapshot provides functionality for capturing, saving, and loading
// environment variable snapshots to disk.
//
// A Snapshot records the full state of an environment (as a key-value map)
// along with a human-readable label and a UTC timestamp. Snapshots can be
// persisted to JSON files and reloaded later, enabling point-in-time
// comparisons between deployments or configuration changes.
//
// Typical usage:
//
//	env, _ := parser.ParseFile(".env.production")
//	snap := snapshot.New("production-2024-01-15", env)
//	snapshot.Save(snap, "snapshots/prod.json")
//
//	// Later, load and compare:
//	old, _ := snapshot.Load("snapshots/prod.json")
//	new, _ := parser.ParseFile(".env.production")
//	changes := differ.Compare(old.Env, new)
package snapshot
