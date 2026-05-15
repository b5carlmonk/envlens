// Package changelog provides a lightweight in-memory changelog for tracking
// environment variable diff history across multiple comparisons.
//
// It is useful when running envlens in watch or pipeline mode, where multiple
// diffs are performed in sequence and a cumulative audit trail is desired.
//
// Usage:
//
//	cl := changelog.New()
//	cl.Add("dev.env", "prod.env", changes)
//	fmt.Print(changelog.RenderText(cl))
package changelog
