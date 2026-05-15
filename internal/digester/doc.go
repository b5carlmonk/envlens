// Package digester provides deterministic SHA-256 hashing for environment
// variable maps, enabling fast equality checks and key-level change detection
// without a full diff pass.
//
// # Overview
//
// Compute produces a Result containing:
//   - A top-level digest of the entire env map (order-independent)
//   - Per-key digests for granular comparison
//
// # Usage
//
//	src := map[string]string{"APP_ENV": "staging", "DB_HOST": "db.local"}
//	dst := map[string]string{"APP_ENV": "production", "DB_HOST": "db.local"}
//
//	ra := digester.Compute(src)
//	rb := digester.Compute(dst)
//
//	if !digester.Equal(ra, rb) {
//		changed := digester.DiffDigests(ra, rb)
//		fmt.Println("Changed keys:", changed)
//	}
package digester
