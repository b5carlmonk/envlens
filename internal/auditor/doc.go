// Package auditor orchestrates the full audit pipeline for envlens.
//
// It combines diffing, risk scoring, and annotation generation into a
// single AuditResult that can be consumed by reporters or CLI output.
//
// Typical usage:
//
//	// Parse both env files, compute diff, then audit.
//	// changes := differ.Compare(oldEnv, newEnv)
//	// result := auditor.Run("old.env", "new.env", changes)
//	// fmt.Println(result.RiskScore.Level)
//
// Annotations are automatically generated for:
//   - Removed keys (warning severity)
//   - Modified sensitive keys such as passwords, tokens, secrets (critical severity)
package auditor
