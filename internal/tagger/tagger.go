package tagger

import (
	"strings"

	"github.com/yourusername/envlens/internal/differ"
)

// Tag represents a semantic label assigned to an environment variable change.
type Tag string

const (
	TagDatabase    Tag = "database"
	TagAuth        Tag = "auth"
	TagNetwork     Tag = "network"
	TagFeatureFlag Tag = "feature-flag"
	TagSecret      Tag = "secret"
	TagLogging     Tag = "logging"
	TagUnknown     Tag = "unknown"
)

// TaggedChange wraps a differ.Change with one or more semantic tags.
type TaggedChange struct {
	Change differ.Change
	Tags   []Tag
}

// Result holds all tagged changes produced by Apply.
type Result struct {
	Tagged []TaggedChange
}

var tagRules = []struct {
	tag      Tag
	keywords []string
}{
	{TagDatabase, []string{"DB_", "DATABASE_", "POSTGRES", "MYSQL", "MONGO", "REDIS", "DSN"}},
	{TagAuth, []string{"AUTH_", "JWT_", "TOKEN", "SECRET", "PASSWORD", "OAUTH", "API_KEY"}},
	{TagNetwork, []string{"HOST", "PORT", "URL", "ADDR", "ENDPOINT", "PROXY", "TLS"}},
	{TagFeatureFlag, []string{"FEATURE_", "FLAG_", "ENABLE_", "DISABLE_"}},
	{TagSecret, []string{"SECRET", "PRIVATE_KEY", "CERT", "PASSPHRASE", "CREDENTIAL"}},
	{TagLogging, []string{"LOG_", "LOGGING_", "LOG_LEVEL", "DEBUG", "VERBOSE"}},
}

// Apply assigns semantic tags to each change based on key naming conventions.
func Apply(changes []differ.Change) Result {
	tagged := make([]TaggedChange, 0, len(changes))
	for _, c := range changes {
		tags := classify(c.Key)
		tagged = append(tagged, TaggedChange{Change: c, Tags: tags})
	}
	return Result{Tagged: tagged}
}

func classify(key string) []Tag {
	upper := strings.ToUpper(key)
	seen := map[Tag]bool{}
	var tags []Tag
	for _, rule := range tagRules {
		for _, kw := range rule.keywords {
			if strings.Contains(upper, kw) && !seen[rule.tag] {
				tags = append(tags, rule.tag)
				seen[rule.tag] = true
				break
			}
		}
	}
	if len(tags) == 0 {
		tags = []Tag{TagUnknown}
	}
	return tags
}
