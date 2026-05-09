// Package templater provides functionality for validating environment variable
// maps against a declared template of required and optional keys.
//
// A Template defines which keys are expected in an environment configuration.
// Keys marked as "required" must be present; keys marked as "optional" may
// appear but are not enforced. Any key present in the env map that is not
// declared in the template is flagged as unexpected.
//
// Example usage:
//
//	 tmpl := &templater.Template{
//	     Required: []string{"DB_HOST", "DB_PORT"},
//	     Optional: []string{"LOG_LEVEL"},
//	 }
//	 result := templater.Validate(tmpl, envMap)
//	 if result.HasIssues() {
//	     fmt.Println(result.MissingRequired)
//	     fmt.Println(result.UnexpectedKeys)
//	 }
//
// Templates can also be constructed from a map using FromMap, where each
// value is either "required" or "optional".
package templater
