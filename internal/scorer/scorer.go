package scorer

import "github.com/user/envlens/internal/differ"

// RiskLevel represents the severity of a set of environment changes.
type RiskLevel string

const (
	RiskNone   RiskLevel = "none"
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// Score holds the computed risk assessment for a diff result.
type Score struct {
	Level   RiskLevel
	Points  int
	Reasons []string
}

// sensitivePatterns are substrings that indicate a key may hold sensitive data.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL",
}

// Evaluate computes a risk Score for the provided list of changes.
func Evaluate(changes []differ.Change) Score {
	points := 0
	var reasons []string

	for _, c := range changes {
		switch c.Type {
		case differ.Added:
			points += 1
		case differ.Removed:
			points += 2
			reasons = append(reasons, "removed key: "+c.Key)
		case differ.Modified:
			points += 3
			if isSensitive(c.Key) {
				points += 5
				reasons = append(reasons, "sensitive key modified: "+c.Key)
			}
		}
	}

	return Score{
		Level:   levelFromPoints(points),
		Points:  points,
		Reasons: reasons,
	}
}

func isSensitive(key string) bool {
	for _, pattern := range sensitivePatterns {
		if containsUpper(key, pattern) {
			return true
		}
	}
	return false
}

func containsUpper(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func levelFromPoints(points int) RiskLevel {
	switch {
	case points == 0:
		return RiskNone
	case points <= 3:
		return RiskLow
	case points <= 8:
		return RiskMedium
	default:
		return RiskHigh
	}
}
