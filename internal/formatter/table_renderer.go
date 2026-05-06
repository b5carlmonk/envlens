package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

const (
	colKeyWidth   = 32
	colTypeWidth  = 10
	colValueWidth = 40
)

// RenderTable writes a human-readable table of environment variable changes
// to the provided writer. Each row shows the change type, key, old value,
// and new value (truncated if necessary).
func RenderTable(w io.Writer, changes []differ.Change) error {
	header := fmt.Sprintf(
		"%-*s %-*s %-*s %s\n",
		colTypeWidth, "TYPE",
		colKeyWidth, "KEY",
		colValueWidth, "OLD VALUE",
		"NEW VALUE",
	)
	separator := strings.Repeat("-", colTypeWidth+colKeyWidth+colValueWidth*2+6) + "\n"

	if _, err := fmt.Fprint(w, header); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, separator); err != nil {
		return err
	}

	for _, c := range changes {
		row := fmt.Sprintf(
			"%-*s %-*s %-*s %s\n",
			colTypeWidth, strings.ToUpper(string(c.Type)),
			colKeyWidth, truncate(c.Key, colKeyWidth),
			colValueWidth, truncate(c.OldValue, colValueWidth),
			truncate(c.NewValue, colValueWidth),
		)
		if _, err := fmt.Fprint(w, row); err != nil {
			return err
		}
	}

	return nil
}

// truncate shortens s to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
