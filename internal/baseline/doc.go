// Package baseline provides functionality for saving and loading named
// reference points of environment variable sets.
//
// A baseline captures the state of an environment at a specific point in time,
// enabling drift detection by comparing the baseline against a current env map.
//
// # Usage
//
//	b := baseline.New("prod-v1", ".env.prod", envMap)
//	baseline.Save(b, "baselines/prod-v1.json")
//
//	loaded, _ := baseline.Load("baselines/prod-v1.json")
//	drifted := baseline.DriftKeys(loaded, currentEnv)
//
// Rendering helpers are available via RenderText and RenderJSON to display
// baseline metadata and drift results in human-readable or machine-readable form.
package baseline
