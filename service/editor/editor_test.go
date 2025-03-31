package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"git-issues/domain"
)

var osCreateTemp = os.CreateTemp

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetIssueContentFromEditor(t *testing.T) {
	// ARRANGE
	mockEditor := filepath.Join(os.TempDir(), "mock-editor")
	if runtime.GOOS == "windows" {
		mockEditor += ".bat"
	} else {
		mockEditor += ".sh"
	}

	setupMockEditor := func(content string) {
		var script string
		if runtime.GOOS == "windows" {
			script = fmt.Sprintf(`(@echo %s) > %%1`, strings.ReplaceAll(content, "\n", "\r\necho."))
		} else {
			script = fmt.Sprintf(`#!/bin/sh echo '%s' > "$1"`, content)
		}
		os.WriteFile(mockEditor, []byte(script), 0700)
	}

	t.Cleanup(func() {
		os.Remove(mockEditor)
	})

	tests := []struct {
		name         string
		service      *Service
		initialTitle string
		initialBody  string
		mockContent  string
		wantTitle    string
		wantBody     string
		wantErr      bool
		setup        func()
	}{
		{
			name:         "successful edit with title and body",
			service:      New(&domain.Config{Editor: mockEditor}),
			initialTitle: "Test Title",
			initialBody:  "Test Body",
			mockContent:  "New Title\n\nNew Body",
			wantTitle:    "New Title",
			wantBody:     "New Body",
			wantErr:      false,
			setup:        func() { setupMockEditor("New Title\n\nNew Body") },
		},
		{
			name:         "successful edit with title only",
			service:      New(&domain.Config{Editor: mockEditor}),
			initialTitle: "Test Title",
			initialBody:  "",
			mockContent:  "New Title Only",
			wantTitle:    "New Title Only",
			wantBody:     "",
			wantErr:      false,
			setup:        func() { setupMockEditor("New Title Only") },
		},
		{
			name:         "error when editor fails",
			service:      New(&domain.Config{Editor: "/non/existent/editor"}),
			initialTitle: "Test",
			initialBody:  "Test",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setup != nil {
				originalCreateTemp := osCreateTemp
				if tt.setup != nil {
					tt.setup()
				}
				defer func() { osCreateTemp = originalCreateTemp }()
			}

			// ACT
			gotTitle, gotBody, err := tt.service.GetIssueContentFromEditor(tt.initialTitle, tt.initialBody)

			// ASSERT
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error = %v, got err = %v", tt.wantErr, err)
				return
			}

			if !tt.wantErr {
				if gotTitle != tt.wantTitle {
					t.Errorf("expected title = %q, got %q", tt.wantTitle, gotTitle)
				}
				if gotBody != tt.wantBody {
					t.Errorf("expected body = %q, got %q", tt.wantBody, gotBody)
				}
			}
		})
	}
}

func TestGetEditorSuccess(t *testing.T) {
	originalGOOS := goos
	t.Cleanup(func() { goos = originalGOOS })

	// ARRANGE
	tests := []struct {
		name    string
		service *Service
		goos    string
		want    string
	}{
		{
			name:    "default windows success",
			service: New(&domain.Config{}),
			goos:    "windows",
			want:    "notepad",
		},
		{
			name:    "default unix success",
			goos:    "linux",
			service: New(&domain.Config{}),
			want:    "vi",
		},
		{
			name: "editor set success",
			goos: "windows",
			service: New(&domain.Config{
				Editor: "notepad++",
			}),
			want: "notepad++",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goos = tt.goos
			// ACT
			got := tt.service.getEditor()

			// ASSERT
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}

		})
	}

}
