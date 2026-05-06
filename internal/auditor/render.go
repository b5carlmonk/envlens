package auditor

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// RenderText writes a human-readable audit summary to w.
func RenderText(w io.Writer, r AuditResult) {
	fmt.Fprintf(w, "Audit Report\n")
	fmt.Fprintf(w, "============\n")
	fmt.Fprintf(w, "Source : %s\n", r.Source)
	fmt.Fprintf(w, "Target : %s\n", r.Target)
	fmt.Fprintf(w, "Time   : %s\n", r.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Risk   : %s (score: %d)\n", r.RiskScore.Level, r.RiskScore.Points)
	fmt.Fprintf(w, "Changes: %d\n", len(r.Changes))

	if len(r.Annotations) > 0 {
		fmt.Fprintf(w, "\nAnnotations:\n")
		for _, ann := range r.Annotations {
			severityTag := strings.ToUpper(ann.Severity)
			fmt.Fprintf(w, "  [%s] %s — %s\n", severityTag, ann.Key, ann.Message)
		}
	}
}

// RenderJSON writes the audit result as JSON to w.
func RenderJSON(w io.Writer, r AuditResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
