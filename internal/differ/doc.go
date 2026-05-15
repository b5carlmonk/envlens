// Package differ provides functionality for comparing two environment
// variable maps and producing a structured list of changes.
//
// Changes are categorized as Added, Removed, or Modified. The Compare
// function accepts two maps representing the source and target environments
// and returns a slice of Change values suitable for further processing
// by reporter, filter, scorer, and other envlens packages.
package differ
