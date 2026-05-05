package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed environment file as a key-value map.
type EnvMap map[string]string

// ParseFile reads a .env file and returns an EnvMap.
// It skips blank lines and comments (lines starting with '#').
// It supports KEY=VALUE and KEY="VALUE" formats.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file %q: %w", path, err)
	}

	return env, nil
}

// parseLine parses a single KEY=VALUE line.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format %q: expected KEY=VALUE", line)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	value := strings.TrimSpace(parts[1])
	value = strings.Trim(value, "\"'")

	return key, value, nil
}
