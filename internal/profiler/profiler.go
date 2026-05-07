package profiler

import (
	"fmt"
	"strings"

	"github.com/your-org/envlens/internal/differ"
)

// Profile represents a named environment profile with categorized keys.
type Profile struct {
	Name       string
	Categories map[string][]string
}

// Result holds the profiling output for a set of changes.
type Result struct {
	Profile     Profile
	Matched     map[string][]string
	Unmatched   []string
	TotalKeys   int
	MatchedKeys int
}

// BuiltinProfiles contains predefined key category patterns.
var BuiltinProfiles = map[string]Profile{
	"web": {
		Name: "web",
		Categories: map[string][]string{
			"database": {"DB_", "DATABASE_", "POSTGRES_", "MYSQL_", "MONGO_"},
			"auth":     {"AUTH_", "JWT_", "OAUTH_", "SECRET_", "TOKEN_"},
			"server":   {"HOST", "PORT", "ADDR", "BIND_"},
			"cache":    {"REDIS_", "CACHE_", "MEMCACHE_"},
			"logging":  {"LOG_", "LOGGER_", "DEBUG", "VERBOSE"},
		},
	},
	"cloud": {
		Name: "cloud",
		Categories: map[string][]string{
			"aws":   {"AWS_"},
			"gcp":   {"GCP_", "GOOGLE_"},
			"azure": {"AZURE_"},
			"k8s":   {"K8S_", "KUBE_", "KUBERNETES_"},
		},
	},
}

// Run profiles the given changes against the named profile.
// If profileName is empty, the "web" profile is used.
func Run(changes []differ.Change, profileName string) (Result, error) {
	if profileName == "" {
		profileName = "web"
	}
	profile, ok := BuiltinProfiles[profileName]
	if !ok {
		return Result{}, fmt.Errorf("profiler: unknown profile %q; available: %s",
			profileName, strings.Join(availableProfiles(), ", "))
	}

	matched := make(map[string][]string)
	unmatched := []string{}

	for _, c := range changes {
		cat := categorize(c.Key, profile)
		if cat == "" {
			unmatched = append(unmatched, c.Key)
		} else {
			matched[cat] = append(matched[cat], c.Key)
		}
	}

	total := len(changes)
	matchedCount := total - len(unmatched)

	return Result{
		Profile:     profile,
		Matched:     matched,
		Unmatched:   unmatched,
		TotalKeys:   total,
		MatchedKeys: matchedCount,
	}, nil
}

func categorize(key string, p Profile) string {
	upper := strings.ToUpper(key)
	for cat, prefixes := range p.Categories {
		for _, prefix := range prefixes {
			if strings.HasPrefix(upper, prefix) || upper == strings.TrimSuffix(prefix, "_") {
				return cat
			}
		}
	}
	return ""
}

func availableProfiles() []string {
	names := make([]string, 0, len(BuiltinProfiles))
	for k := range BuiltinProfiles {
		names = append(names, k)
	}
	return names
}
