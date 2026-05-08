package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved reference point of environment variables
// used for future comparisons and drift detection.
type Baseline struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Source    string            `json:"source"`
	Env       map[string]string `json:"env"`
}

// New creates a new Baseline from the given env map.
func New(name, source string, env map[string]string) *Baseline {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Baseline{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Source:    source,
		Env:       copy,
	}
}

// Save writes the baseline to a JSON file at the given path.
func Save(b *Baseline, path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("baseline: write failed: %w", err)
	}
	return nil
}

// Load reads a baseline from a JSON file at the given path.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read failed: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal failed: %w", err)
	}
	return &b, nil
}

// DriftKeys returns keys whose values differ between the baseline and a current env map.
func DriftKeys(b *Baseline, current map[string]string) []string {
	var drifted []string
	for k, baseVal := range b.Env {
		if curVal, ok := current[k]; ok {
			if curVal != baseVal {
				drifted = append(drifted, k)
			}
		}
	}
	return drifted
}
