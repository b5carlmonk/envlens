package reporter

import (
	"encoding/json"
	"io"
)

type jsonChange struct {
	Type     string `json:"type"`
	Key      string `json:"key"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

type jsonReport struct {
	FromFile string       `json:"from_file"`
	ToFile   string       `json:"to_file"`
	Summary  string       `json:"summary"`
	Changes  []jsonChange `json:"changes"`
}

func renderJSON(w io.Writer, r *Report) error {
	changes := make([]jsonChange, 0, len(r.Diff))
	for _, c := range r.Diff {
		changes = append(changes, jsonChange{
			Type:     string(c.Type),
			Key:      c.Key,
			OldValue: c.OldValue,
			NewValue: c.NewValue,
		})
	}
	payload := jsonReport{
		FromFile: r.FromFile,
		ToFile:   r.ToFile,
		Summary:  r.Summary(),
		Changes:  changes,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
