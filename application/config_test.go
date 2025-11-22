package application

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"git-issues/domain"
)

func TestLoadConfig_Success(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	file := filepath.Join(dir, "config.json")

	expected := &domain.Config{
		Token:      "token123",
		Owner:      "owner-name",
		Repo:       "repo-name",
		Editor:     "vim",
		APIBaseURL: "https://api.example.com",
	}

	data, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal expected config: %v", err)
	}

	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Act
	cfg, err := LoadConfig(file)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	// Assert
	if cfg.Token != expected.Token {
		t.Fatalf("Token mismatch: got %q want %q", cfg.Token, expected.Token)
	}
	if cfg.Owner != expected.Owner {
		t.Fatalf("Owner mismatch: got %q want %q", cfg.Owner, expected.Owner)
	}
	if cfg.Repo != expected.Repo {
		t.Fatalf("Repo mismatch: got %q want %q", cfg.Repo, expected.Repo)
	}
	if cfg.Editor != expected.Editor {
		t.Fatalf("Editor mismatch: got %q want %q", cfg.Editor, expected.Editor)
	}
	if cfg.APIBaseURL != expected.APIBaseURL {
		t.Fatalf("APIBaseURL mismatch: got %q want %q", cfg.APIBaseURL, expected.APIBaseURL)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	file := filepath.Join(dir, "config.json")

	// syntactically invalid JSON (truncated)
	invalidJSON := `{"token":"token123",`

	if err := os.WriteFile(file, []byte(invalidJSON), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Act
	_, err := LoadConfig(file)

	// Assert
	if err == nil {
		t.Fatal("expected unmarshal error, got nil")
	}
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError, got %T: %v", err, err)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	file := filepath.Join(dir, "does-not-exist.json")

	// Act
	_, err := LoadConfig(file)

	// Assert
	if err == nil {
		t.Fatal("expected read error, got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected file-not-found error, got %T: %v", err, err)
	}
}
