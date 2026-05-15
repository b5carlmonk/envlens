// Package inspector provides deep inspection of individual environment
// variable changes, enriching each change with metadata such as value
// length delta, type hints, and whether the value appears numeric,
// boolean, or URL-like.
package inspector

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// TypeHint describes the inferred type of an environment variable value.
type TypeHint string

const (
	TypeString  TypeHint = "string"
	TypeNumeric TypeHint = "numeric"
	TypeBoolean TypeHint = "boolean"
	TypeURL     TypeHint = "url"
)

// Inspection holds enriched metadata for a single Change.
type Inspection struct {
	Change      differ.Change
	OldLen      int
	NewLen      int
	LengthDelta int
	OldHint     TypeHint
	NewHint     TypeHint
	ValueSame   bool
}

// Inspect analyses each change and returns a slice of Inspection values.
func Inspect(changes []differ.Change) []Inspection {
	out := make([]Inspection, 0, len(changes))
	for _, c := range changes {
		out = append(out, inspect(c))
	}
	return out
}

func inspect(c differ.Change) Inspection {
	oldLen := len(c.OldValue)
	newLen := len(c.NewValue)
	return Inspection{
		Change:      c,
		OldLen:      oldLen,
		NewLen:      newLen,
		LengthDelta: newLen - oldLen,
		OldHint:     inferType(c.OldValue),
		NewHint:     inferType(c.NewValue),
		ValueSame:   c.OldValue == c.NewValue,
	}
}

func inferType(v string) TypeHint {
	if v == "" {
		return TypeString
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return TypeNumeric
	}
	lower := strings.ToLower(v)
	if lower == "true" || lower == "false" || lower == "yes" || lower == "no" {
		return TypeBoolean
	}
	if u, err := url.ParseRequestURI(v); err == nil && u.Scheme != "" {
		return TypeURL
	}
	return TypeString
}
