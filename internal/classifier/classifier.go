// Package classifier categorizes environment variable changes by their
// inferred purpose based on key naming patterns and value heuristics.
package classifier

import (
	"strings"

	"github.com/user/envlens/internal/differ"
)

// Category represents the inferred purpose of an environment variable.
type Category string

const (
	CategoryDatabase  Category = "database"
	CategoryAuth      Category = "auth"
	CategoryNetwork   Category = "network"
	CategoryLogging   Category = "logging"
	CategoryFeature   Category = "feature"
	CategoryObservability Category = "observability"
	CategoryUnknown   Category = "unknown"
)

// Result holds a classified change with its inferred category.
type Result struct {
	Change   differ.Change
	Category Category
	Reason   string
}

var categoryRules = []struct {
	category Category
	reason   string
	keywords []string
}{
	{CategoryDatabase, "key matches database pattern", []string{"DB", "DATABASE", "POSTGRES", "MYSQL", "MONGO", "REDIS", "DSN", "SQL"}},
	{CategoryAuth, "key matches auth pattern", []string{"AUTH", "TOKEN", "SECRET", "PASSWORD", "API_KEY", "OAUTH", "JWT", "CREDENTIAL"}},
	{CategoryNetwork, "key matches network pattern", []string{"HOST", "PORT", "URL", "ADDR", "ENDPOINT", "PROXY", "TLS", "SSL", "DOMAIN"}},
	{CategoryLogging, "key matches logging pattern", []string{"LOG", "LOGGER", "LOG_LEVEL", "DEBUG", "VERBOSE", "TRACE"}},
	{CategoryFeature, "key matches feature flag pattern", []string{"FEATURE", "FLAG", "ENABLE", "DISABLE", "TOGGLE"}},
	{CategoryObservability, "key matches observability pattern", []string{"METRIC", "TRACE", "SENTRY", "DATADOG", "OTEL", "PROMETHEUS", "GRAFANA"}},
}

// Apply classifies each change in the provided slice and returns Results.
func Apply(changes []differ.Change) []Result {
	results := make([]Result, 0, len(changes))
	for _, c := range changes {
		cat, reason := classify(c.Key)
		results = append(results, Result{
			Change:   c,
			Category: cat,
			Reason:   reason,
		})
	}
	return results
}

func classify(key string) (Category, string) {
	upper := strings.ToUpper(key)
	for _, rule := range categoryRules {
		for _, kw := range rule.keywords {
			if strings.Contains(upper, kw) {
				return rule.category, rule.reason
			}
		}
	}
	return CategoryUnknown, "no matching pattern found"
}
