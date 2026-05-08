package comparator

import (
	"fmt"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/parser"
)

// Result holds the outcome of comparing two env files.
type Result struct {
	SourceFile string
	TargetFile string
	Changes    []differ.Change
	SourceEnv  map[string]string
	TargetEnv  map[string]string
}

// Options controls comparator behaviour.
type Options struct {
	// StrictMode treats missing keys as errors rather than changes.
	StrictMode bool
}

// Compare parses both env files and returns a Result containing all
// detected changes between them.
func Compare(sourceFile, targetFile string, opts Options) (*Result, error) {
	src, err := parser.ParseFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("comparator: reading source %q: %w", sourceFile, err)
	}

	tgt, err := parser.ParseFile(targetFile)
	if err != nil {
		return nil, fmt.Errorf("comparator: reading target %q: %w", targetFile, err)
	}

	if opts.StrictMode {
		if err := strictCheck(src, tgt); err != nil {
			return nil, err
		}
	}

	changes := differ.Compare(src, tgt)

	return &Result{
		SourceFile: sourceFile,
		TargetFile: targetFile,
		Changes:    changes,
		SourceEnv:  src,
		TargetEnv:  tgt,
	}, nil
}

// strictCheck returns an error if any key present in source is absent in target.
func strictCheck(src, tgt map[string]string) error {
	for k := range src {
		if _, ok := tgt[k]; !ok {
			return fmt.Errorf("comparator: strict mode: key %q removed in target", k)
		}
	}
	return nil
}
