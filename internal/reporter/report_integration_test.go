package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/parser"
	"github.com/user/envlens/internal/reporter"
	"os"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envlens-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestIntegration_ParseDiffReport(t *testing.T) {
	oldContent := "APP_ENV=production\nDB_HOST=db.prod\nSECRET=abc123\n"
	newContent := "APP_ENV=staging\nDB_HOST=db.prod\nNEW_VAR=enabled\n"

	oldFile := writeTempEnvFile(t, oldContent)
	newFile := writeTempEnvFile(t, newContent)

	oldEnv, err := parser.ParseFile(oldFile)
	if err != nil {
		t.Fatalf("ParseFile old: %v", err)
	}
	newEnv, err := parser.ParseFile(newFile)
	if err != nil {
		t.Fatalf("ParseFile new: %v", err)
	}

	changes := differ.Compare(oldEnv, newEnv)
	r := reporter.NewReport(oldFile, newFile, changes)

	if !r.HasChanges() {
		t.Fatal("expected changes")
	}

	var buf bytes.Buffer
	if err := r.Render(&buf, reporter.FormatText); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected SECRET (removed) in output")
	}
	if !strings.Contains(out, "NEW_VAR") {
		t.Error("expected NEW_VAR (added) in output")
	}
}
