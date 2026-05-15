package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Result holds the digest output for an environment map.
type Result struct {
	Digest     string            `json:"digest"`
	KeyCount   int               `json:"key_count"`
	KeyDigests map[string]string `json:"key_digests"`
}

// Compute generates a deterministic SHA-256 digest for the given env map.
// Individual per-key digests are also computed for granular comparison.
func Compute(env map[string]string) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	keyDigests := make(map[string]string, len(keys))
	var sb strings.Builder

	for _, k := range keys {
		v := env[k]
		entry := fmt.Sprintf("%s=%s", k, v)
		keyDigests[k] = hashString(entry)
		sb.WriteString(entry)
		sb.WriteByte('\n')
	}

	return Result{
		Digest:     hashString(sb.String()),
		KeyCount:   len(keys),
		KeyDigests: keyDigests,
	}
}

// Equal returns true if two Results share the same top-level digest.
func Equal(a, b Result) bool {
	return a.Digest == b.Digest
}

// DiffDigests returns keys whose per-key digests differ between two Results.
func DiffDigests(a, b Result) []string {
	seen := make(map[string]struct{})
	var changed []string

	for k, da := range a.KeyDigests {
		seen[k] = struct{}{}
		if db, ok := b.KeyDigests[k]; !ok || da != db {
			changed = append(changed, k)
		}
	}
	for k := range b.KeyDigests {
		if _, ok := seen[k]; !ok {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return changed
}

func hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
