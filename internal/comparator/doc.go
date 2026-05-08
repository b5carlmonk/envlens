// Package comparator provides a high-level API for comparing two .env files.
//
// It combines the parser and differ packages into a single Compare call,
// returning a Result that contains the parsed environments and the list of
// detected changes.
//
// Basic usage:
//
//	res, err := comparator.Compare(".env.staging", ".env.production", comparator.Options{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	comparator.RenderText(os.Stdout, res)
//
// StrictMode
//
// When Options.StrictMode is true, Compare returns an error if any key
// present in the source file is absent from the target file. This is useful
// for CI pipelines where key removal must be an explicit, reviewed action.
package comparator
