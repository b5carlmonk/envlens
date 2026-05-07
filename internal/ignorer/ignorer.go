package ignorer

import (
	"bufio"
	"os"
	"strings"
)

// Options configures the ignore behaviour.
type Options struct {
	// Keys is an explicit list of env var keys to ignore.
	Keys []string
	// Prefixes causes any key starting with one of these prefixes to be ignored.
	Prefixes []string
}

// FromFile reads an ignore file where each non-blank, non-comment line is
// treated as a key or prefix pattern (prefix patterns end with '*').
func FromFile(path string) (Options, error) {
	f, err := os.Open(path)
	if err != nil {
		return Options{}, err
	}
	defer f.Close()

	var opts Options
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, "*") {
			opts.Prefixes = append(opts.Prefixes, strings.TrimSuffix(line, "*"))
		} else {
			opts.Keys = append(opts.Keys, line)
		}
	}
	return opts, scanner.Err()
}

// Apply removes any entry whose key matches an ignored key or prefix.
// It returns a new map without the ignored keys.
func Apply(env map[string]string, opts Options) map[string]string {
	ignored := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		ignored[k] = true
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if ignored[k] {
			continue
		}
		skip := false
		for _, p := range opts.Prefixes {
			if strings.HasPrefix(k, p) {
				skip = true
				break
			}
		}
		if !skip {
			result[k] = v
		}
	}
	return result
}
