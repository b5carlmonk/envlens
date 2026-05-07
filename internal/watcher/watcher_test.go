package watcher_test

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/envlens/internal/watcher"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestNew_ValidPaths(t *testing.T) {
	w, err := watcher.New("source.env", "target.env", watcher.WatchOptions{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestNew_EmptySource_ReturnsError(t *testing.T) {
	_, err := watcher.New("", "target.env", watcher.WatchOptions{})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNew_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := watcher.New("source.env", "", watcher.WatchOptions{})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestWatcher_DetectsFileChange(t *testing.T) {
	dir := t.TempDir()
	src := writeTempFile(t, dir, "source.env", "KEY=old")
	tgt := writeTempFile(t, dir, "target.env", "KEY=old")

	var callCount atomic.Int32

	opts := watcher.WatchOptions{
		PollInterval: 50 * time.Millisecond,
		OnChange: func(source, target string) {
			callCount.Add(1)
		},
	}

	w, err := watcher.New(src, tgt, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go w.Start()
	defer w.Stop()

	// Wait for initial poll to settle
	time.Sleep(120 * time.Millisecond)

	// Modify the target file
	if err := os.WriteFile(tgt, []byte("KEY=new"), 0644); err != nil {
		t.Fatalf("failed to update file: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if callCount.Load() == 0 {
		t.Error("expected OnChange to be called after file modification")
	}
}

func TestWatcher_Stop_StopsPolling(t *testing.T) {
	dir := t.TempDir()
	src := writeTempFile(t, dir, "source.env", "A=1")
	tgt := writeTempFile(t, dir, "target.env", "A=1")

	var callCount atomic.Int32

	opts := watcher.WatchOptions{
		PollInterval: 30 * time.Millisecond,
		OnChange: func(_, _ string) {
			callCount.Add(1)
		},
	}

	w, _ := watcher.New(src, tgt, opts)
	go w.Start()
	time.Sleep(50 * time.Millisecond)
	w.Stop()

	// Modify after stop — should not trigger callback
	before := callCount.Load()
	os.WriteFile(tgt, []byte("A=2"), 0644)
	time.Sleep(100 * time.Millisecond)

	if callCount.Load() > before+1 {
		t.Error("OnChange was called after Stop")
	}
}
