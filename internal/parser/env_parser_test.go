package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env["APP_ENV"])
	}
	if env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", env["PORT"])
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# This is a comment\n\nDB_HOST=localhost\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
	if env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", env["DB_HOST"])
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret value" {
		t.Errorf("expected unquoted value, got %q", env["SECRET"])
	}
	if env["TOKEN"] != "abc123" {
		t.Errorf("expected unquoted value, got %q", env["TOKEN"])
	}
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
