// Package parser provides utilities for reading and parsing .env files
// into structured key-value maps for use by the envlens diffing and
// auditing engine.
//
// Supported .env formats:
//
//	# Comment lines are ignored
//	KEY=value
//	KEY="quoted value"
//	KEY='single quoted value'
//
// Blank lines are also ignored. Keys must be non-empty and values are
// automatically stripped of surrounding quotes.
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env["DATABASE_URL"])
package parser
