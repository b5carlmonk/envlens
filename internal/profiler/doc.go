// Package profiler categorizes environment variable changes against
// named profiles (e.g. "web", "cloud") to help operators quickly
// understand which subsystems are affected by a deployment diff.
//
// Built-in profiles group keys by common prefixes:
//
//	"web"   — database, auth, server, cache, logging
//	"cloud" — aws, gcp, azure, k8s
//
// Usage:
//
//	result, err := profiler.Run(changes, "web")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Matched["database"])
package profiler
