package auditor

import (
	"time"

	"github.com/yourusername/envlens/internal/differ"
	"github.com/yourusername/envlens/internal/scorer"
)

// AuditResult holds the full audit output for a diff operation.
type AuditResult struct {
	Timestamp   time.Time
	Source      string
	Target      string
	Changes     []differ.Change
	RiskScore   scorer.Score
	Annotations []Annotation
}

// Annotation attaches a human-readable note to a specific change key.
type Annotation struct {
	Key     string
	Message string
	Severity string // "info", "warning", "critical"
}

// Run performs a full audit: diff, score, and annotate.
func Run(source, target string, changes []differ.Change) AuditResult {
	score := scorer.Evaluate(changes)
	annotations := annotate(changes)

	return AuditResult{
		Timestamp:   time.Now().UTC(),
		Source:      source,
		Target:      target,
		Changes:     changes,
		RiskScore:   score,
		Annotations: annotations,
	}
}

// annotate generates annotations for notable changes.
func annotate(changes []differ.Change) []Annotation {
	var annotations []Annotation
	for _, c := range changes {
		if ann, ok := buildAnnotation(c); ok {
			annotations = append(annotations, ann)
		}
	}
	return annotations
}

func buildAnnotation(c differ.Change) (Annotation, bool) {
	switch c.Type {
	case differ.Removed:
		return Annotation{
			Key:      c.Key,
			Message:  "Key was removed; ensure dependents are updated.",
			Severity: "warning",
		}, true
	case differ.Modified:
		if isSensitiveKey(c.Key) {
			return Annotation{
				Key:      c.Key,
				Message:  "Sensitive key was modified; verify intentional rotation.",
				Severity: "critical",
			}, true
		}
	}
	return Annotation{}, false
}

func isSensitiveKey(key string) bool {
	sensitivePatterns := []string{"SECRET", "PASSWORD", "TOKEN", "API_KEY", "PRIVATE"}
	for _, p := range sensitivePatterns {
		if containsIgnoreCase(key, p) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > 0 && containsStr(toUpper(s), toUpper(substr)))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func toUpper(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		b[i] = c
	}
	return string(b)
}
