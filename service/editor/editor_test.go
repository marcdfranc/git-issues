package editor

import (
	"errors"
	"fmt"
	"os"

	"runtime"
	"strings"
	"testing"

	"git-issues/domain"
)

func TestNew(t *testing.T) {
	config := &domain.Config{Editor: "vim"}
	service := New(config)

	if service.config != config {
		t.Errorf("Expected config %v, got %v", config, service.config)
	}
}

func TestService_GetEditor(t *testing.T) {
	tests := []struct {
		name           string
		configEditor   string
		envEditor      string
		expectedEditor string
		setupEnv       func()
		cleanupEnv     func()
	}{
		{
			name:           "should use config editor when set",
			configEditor:   "nano",
			expectedEditor: "nano",
		},
		{
			name:           "should use EDITOR environment variable when config is empty",
			configEditor:   "",
			envEditor:      "emacs",
			expectedEditor: "emacs",
			setupEnv: func() {
				err := os.Setenv("EDITOR", "emacs")
				if err != nil {
					t.Fatalf("Failed to set EDITOR env: %v", err)
				}
			},
			cleanupEnv: func() {
				err := os.Unsetenv("EDITOR")
				if err != nil {
					t.Fatalf("Failed to unset EDITOR env: %v", err)
				}
			},
		},
		{
			name:           "should use default editor when no config or env is set",
			configEditor:   "",
			envEditor:      "",
			expectedEditor: getDefaultEditor(),
			setupEnv: func() {
				err := os.Unsetenv("EDITOR")
				if err != nil {
					t.Fatalf("Failed to unset EDITOR env: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				tt.setupEnv()
			}
			if tt.cleanupEnv != nil {
				defer tt.cleanupEnv()
			}

			service := &Service{
				config: &domain.Config{Editor: tt.configEditor},
			}

			editor := service.getEditor()
			if editor != tt.expectedEditor {
				t.Errorf("Expected editor %s, got %s", tt.expectedEditor, editor)
			}
		})
	}
}

func getDefaultEditor() string {
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "vi"
}

func TestService_GetIssueContentFromEditor(t *testing.T) {
	// Create platform-appropriate mock editors
	mockEditorPreserveTitle := createMockEditorScript(t, "New Title", `New issue body`)
	mockEditorOnlyTitle := createMockEditorScript(t, "Unique Title", "")

	tests := []struct {
		name          string
		issue         *domain.Issue
		editor        string
		mockEditor    string
		expectedTitle string
		expectedBody  string
		expectError   bool
		errorContains string
	}{
		{
			name: "should preserve existing title and update body",
			issue: &domain.Issue{
				Title: "Existing Title",
				Body:  "Existing body",
			},
			mockEditor:    mockEditorPreserveTitle,
			expectedTitle: "Existing Title", // Title should be preserved when already set
			expectedBody:  "New issue body",
			expectError:   false,
		},
		{
			name: "should update title when empty",
			issue: &domain.Issue{
				Title: "",
				Body:  "Corpo antigo",
			},
			mockEditor:    mockEditorPreserveTitle,
			expectedTitle: "New Title",
			expectedBody:  "New issue body",
		},
		{
			name: "should set body empty when only one line",
			issue: &domain.Issue{
				Title: "",
				Body:  "Corpo antigo",
			},
			mockEditor:    mockEditorOnlyTitle,
			expectedTitle: "Unique Title",
			expectedBody:  "",
		},
		{
			name: "should handle editor command failure",
			issue: &domain.Issue{
				Title: "Test",
				Body:  "Test",
			},
			editor:        "/nonexistent/editor",
			expectError:   true,
			errorContains: "could not exec editor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up editor
			editorToUse := tt.editor
			if tt.mockEditor != "" {
				editorToUse = tt.mockEditor
			}

			service := &Service{
				config: &domain.Config{Editor: editorToUse},
			}

			err := service.GetIssueContentFromEditor(tt.issue)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, but got none")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%v'", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.issue.Title != tt.expectedTitle {
					t.Errorf("Expected title '%s', got '%s'", tt.expectedTitle, tt.issue.Title)
				}
				if tt.issue.Body != tt.expectedBody {
					t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, tt.issue.Body)
				}

				if strings.Contains(tt.issue.Title, "\r") {
					t.Errorf("Title contains CR")
				}
				if strings.Contains(tt.issue.Body, "\r") {
					t.Errorf("Body contains CR")
				}

			}
		})
	}
}

// --------------------
// Testes de erros de I/O
// --------------------
func TestService_GetIssueContentFromEditor_IOErrors(t *testing.T) {
	issue := &domain.Issue{Title: "Titulo", Body: "Corpo"}

	// Forçar falha na criação de arquivo temporário
	createTempFile = func(dir, pattern string) (*os.File, error) {
		return nil, errors.New("forced temp file error")
	}
	defer func() { createTempFile = os.CreateTemp }() // restaura

	// Editor não será chamado porque falha antes
	service := &Service{config: &domain.Config{Editor: "dummy"}}
	err := service.GetIssueContentFromEditor(issue)

	if err == nil || !strings.Contains(err.Error(), "could not create temp file") {
		t.Errorf("Expected errCreateTmpFile, got %v", err)
	}
}

// Test error variables
func TestErrorVariables(t *testing.T) {
	// Test that error variables are properly defined
	expectedErrors := map[string]error{
		"errCreateTmpFile": errCreateTmpFile,
		"errWriteTempFile": errWriteTempFile,
		"errExecEditor":    errExecEditor,
		"errReadEditor":    errReadEditor,
	}

	for name, err := range expectedErrors {
		if err == nil {
			t.Errorf("Expected %s to be defined", name)
		}
	}
}

// Helper function to create mock editor scripts
func createMockEditorScript(t *testing.T, title, body string) string {
	t.Helper()

	var content string
	if runtime.GOOS == "windows" {
		switch {
		case body == "":
			// Script batch para Windows que só define o título
			content = fmt.Sprintf(`@echo off
echo %s > "%%1"
echo. >> "%%1"`, title)
		case title == "":
			// Script batch para Windows que só define o corpo
			content = fmt.Sprintf(`@echo off
echo. > "%%1"
echo %s >> "%%1"`, body)
		default:
			// Script batch para Windows que define título e corpo
			content = fmt.Sprintf(`@echo off
echo %s > "%%1"
echo. >> "%%1" 
echo %s >> "%%1"`, title, body)
		}
	} else {
		// Script shell para Unix/Linux
		content = fmt.Sprintf(`#!/bin/sh
echo "%s" > "$1"
echo "%s" >> "$1"`, title, body)
	}

	// Criar arquivo temporário
	var ext string
	if runtime.GOOS == "windows" {
		ext = ".bat"
	} else {
		ext = ".sh"
	}

	tempFile, err := os.CreateTemp("", "mock-editor-*"+ext)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	t.Cleanup(func() {
		err := os.Remove(tempFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	})

	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	err = tempFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Tornar executável em sistemas Unix
	if runtime.GOOS != "windows" {
		err = os.Chmod(tempFile.Name(), 0755)
		if err != nil {
			t.Fatalf("Failed to make script executable: %v", err)
		}
	}

	return tempFile.Name()
}
