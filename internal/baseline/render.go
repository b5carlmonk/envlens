package baseline

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// RenderText writes a human-readable summary of the baseline to w.
func RenderText(b *Baseline, drifted []string, w io.Writer) {
	fmt.Fprintf(w, "Baseline : %s\n", b.Name)
	fmt.Fprintf(w, "Source   : %s\n", b.Source)
	fmt.Fprintf(w, "Created  : %s\n", b.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Keys     : %d\n", len(b.Env))

	if len(drifted) == 0 {
		fmt.Fprintln(w, "Drift    : none")
		return
	}

	sort.Strings(drifted)
	fmt.Fprintf(w, "Drift    : %d key(s) changed\n", len(drifted))
	for _, k := range drifted {
		fmt.Fprintf(w, "  ~ %s\n", k)
	}
}

// RenderJSON writes a JSON representation of the baseline and drift info to w.
func RenderJSON(b *Baseline, drifted []string, w io.Writer) error {
	type output struct {
		Baseline *Baseline `json:"baseline"`
		Drift    []string  `json:"drift"`
		HasDrift bool      `json:"has_drift"`
	}

	if drifted == nil {
		drifted = []string{}
	}

	out := output{
		Baseline: b,
		Drift:    drifted,
		HasDrift: len(drifted) > 0,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline render: %w", err)
	}

	_, err = fmt.Fprintln(w, strings.TrimSpace(string(data)))
	return err
}
