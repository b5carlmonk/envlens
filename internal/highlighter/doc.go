// Package highlighter applies ANSI terminal colour and prefix symbols to
// environment variable diffs, making added, removed, and modified changes
// visually distinct in terminal output.
//
// Usage:
//
//	changes := differ.Compare(source, target)
//	opts := highlighter.DefaultOptions()
//	lines := highlighter.Apply(changes, opts)
//	fmt.Println(highlighter.Render(lines))
//
// Set Options.NoColor = true to strip ANSI codes when writing to files or
// non-interactive terminals. Set Options.ShowOldValue = false to omit the
// previous value annotation on modified keys.
package highlighter
