package patcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/envlens/internal/differ"
)

// Strategy defines how conflicts are resolved when patching.
type Strategy string

const (
	StrategyOverwrite Strategy = "overwrite" // always apply patch value
	StrategySkip      Strategy = "skip"      // keep existing value on conflict
	StrategyError     Strategy = "error"     // return error on conflict
)

// Options controls patch behaviour.
type Options struct {
	Strategy Strategy
	DryRun   bool
}

// Result holds the outcome of a patch operation.
type Result struct {
	Applied  []string
	Skipped  []string
	Conflicts []string
	DryRun   bool
}

// Apply patches the env map in target using the provided diff changes.
// It reads the current target file, applies changes according to opts,
// and writes the result back (unless DryRun is set).
func Apply(targetPath string, changes []differ.Change, opts Options) (Result, error) {
	raw, err := os.ReadFile(targetPath)
	if err != nil {
		return Result{}, fmt.Errorf("patcher: read target: %w", err)
	}

	lines := strings.Split(string(raw), "\n")
	env := parseLines(lines)

	var result Result
	result.DryRun = opts.DryRun

	for _, ch := range changes {
		switch ch.Type {
		case differ.Added, differ.Modified:
			_, exists := env[ch.Key]
			if exists && ch.Type == differ.Modified {
				switch opts.Strategy {
				case StrategySkip:
					result.Skipped = append(result.Skipped, ch.Key)
					continue
				case StrategyError:
					return result, fmt.Errorf("patcher: conflict on key %q", ch.Key)
				}
				result.Conflicts = append(result.Conflicts, ch.Key)
			}
			env[ch.Key] = ch.NewValue
			result.Applied = append(result.Applied, ch.Key)
		case differ.Removed:
			delete(env, ch.Key)
			result.Applied = append(result.Applied, ch.Key)
		}
	}

	if !opts.DryRun {
		if err := writeEnv(targetPath, env); err != nil {
			return result, fmt.Errorf("patcher: write target: %w", err)
		}
	}

	return result, nil
}

func parseLines(lines []string) map[string]string {
	env := make(map[string]string)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}

func writeEnv(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
