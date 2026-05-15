package scorer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderText returns a human-readable text representation of a ScoreResult.
func RenderText(result ScoreResult) string {
	var sb strings.Builder

	sb.WriteString("=== Risk Score Report ===\n")
	sb.WriteString(fmt.Sprintf("Risk Level : %s\n", result.Level))
	sb.WriteString(fmt.Sprintf("Score      : %d\n", result.Score))
	sb.WriteString(fmt.Sprintf("Total Keys : %d\n", result.TotalKeys))

	if len(result.Reasons) == 0 {
		sb.WriteString("Reasons    : none\n")
		return sb.String()
	}

	sb.WriteString("Reasons:\n")
	for _, r := range result.Reasons {
		sb.WriteString(fmt.Sprintf("  - %s\n", r))
	}

	return sb.String()
}

// RenderJSON returns a JSON representation of a ScoreResult.
func RenderJSON(result ScoreResult) (string, error) {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("scorer: failed to render JSON: %w", err)
	}
	return string(b), nil
}
