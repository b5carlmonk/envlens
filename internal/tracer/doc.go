// Package tracer provides a lightweight event tracing mechanism for recording
// sequences of environment diff operations across deployments or sessions.
//
// A Trace accumulates ordered entries, each capturing the label, source,
// target, timestamp, and list of changes produced by a diff operation.
//
// # Usage
//
//	tr := tracer.New()
//	tr.Add("deploy-v2", "staging.env", "prod.env", changes)
//
//	fmt.Println(tracer.RenderText(tr))
//
// Entries can be filtered by label and rendered as plain text or JSON.
// This is useful for audit trails, CI pipelines, and deployment history.
package tracer
