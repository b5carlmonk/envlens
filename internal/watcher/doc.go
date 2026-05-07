// Package watcher provides file-based polling to detect changes between
// two environment files (source and target).
//
// It periodically checks the modification times of both files and invokes
// a user-supplied callback when a change is detected. This is useful for
// long-running processes or CLI watch modes that need to react to env file
// updates in real time.
//
// Example usage:
//
//	w, err := watcher.New(".env.staging", ".env.prod", watcher.WatchOptions{
//		PollInterval: 2 * time.Second,
//		OnChange: func(src, tgt string) {
//			fmt.Printf("Change detected between %s and %s\n", src, tgt)
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer w.Stop()
//	w.Start()
package watcher
