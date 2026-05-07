package watcher

import (
	"fmt"
	"os"
	"time"
)

// WatchOptions configures the file watcher behavior.
type WatchOptions struct {
	// PollInterval is how often to check for file changes.
	PollInterval time.Duration
	// OnChange is called when a change is detected between the two files.
	OnChange func(source, target string)
}

// Watcher monitors two env files for changes and triggers a callback.
type Watcher struct {
	source  string
	target  string
	opts    WatchOptions
	stopCh  chan struct{}
}

// New creates a new Watcher for the given source and target env files.
func New(source, target string, opts WatchOptions) (*Watcher, error) {
	if source == "" || target == "" {
		return nil, fmt.Errorf("watcher: source and target paths must not be empty")
	}
	if opts.PollInterval <= 0 {
		opts.PollInterval = 5 * time.Second
	}
	return &Watcher{
		source: source,
		target: target,
		opts:   opts,
		stopCh: make(chan struct{}),
	}, nil
}

// Start begins watching the files and blocks until Stop is called.
func (w *Watcher) Start() {
	var lastSrcMod, lastTgtMod time.Time

	ticker := time.NewTicker(w.opts.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			srcMod := modTime(w.source)
			tgtMod := modTime(w.target)

			if (!srcMod.IsZero() && srcMod != lastSrcMod) ||
				(!tgtMod.IsZero() && tgtMod != lastTgtMod) {
				lastSrcMod = srcMod
				lastTgtMod = tgtMod
				if w.opts.OnChange != nil {
					w.opts.OnChange(w.source, w.target)
				}
			}
		}
	}
}

// Stop signals the watcher to stop polling.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

// modTime returns the modification time of a file, or zero time on error.
func modTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
