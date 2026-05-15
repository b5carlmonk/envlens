// Package differ provides functionality for comparing two environment variable
// maps and producing a structured list of changes.
//
// # Overview
//
// The differ package accepts two maps of string key-value pairs (representing
// parsed .env files or deployment environment snapshots) and returns a slice
// of Change values describing what was added, removed, or modified.
//
// # Usage
//
//	changes := differ.Compare(source, target)
//	for _, c := range changes {
//		fmt.Println(c.Type, c.Key)
//	}
package differ
