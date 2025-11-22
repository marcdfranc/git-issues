package issue

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"git-issues/domain"
)

type errWriter struct {
	err error
}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, e.err
}

func TestPrintIssues(t *testing.T) {
	tests := []struct {
		name    string
		issues  []domain.Issue
		want    string
		wantErr error
	}{
		{
			name: "successful print multiple issues",
			issues: []domain.Issue{
				{Number: 1, Title: "t1", State: "open"},
				{Number: 2, Title: "t2", State: "closed"},
			},
			want: fmt.Sprintf("\nIssues:\n#%v - %s (%s)\n#%v - %s (%s)\n",
				1, "t1", "open",
				2, "t2", "closed",
			),
			wantErr: nil,
		},
		{
			name:    "empty list prints only header",
			issues:  []domain.Issue{},
			want:    "\nIssues:\n",
			wantErr: nil,
		},
		{
			name:    "writer returns error",
			issues:  []domain.Issue{{Number: 1, Title: "t", State: "open"}},
			want:    "",
			wantErr: errors.New("write fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			var w ioWriter = &buf
			if tt.wantErr != nil {
				// replace writer with one that fails
				w = &errWriter{err: tt.wantErr}
			}

			// Act
			err := PrintIssues(w, tt.issues)

			// Assert
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("unexpected error: got %v want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

// small interface to allow swapping between bytes.Buffer and errWriter
type ioWriter interface {
	Write([]byte) (int, error)
}

func TestPrintIssue(t *testing.T) {
	tests := []struct {
		name    string
		issue   *domain.Issue
		want    string
		wantErr error
	}{
		{
			name: "successful print detail",
			issue: &domain.Issue{
				Number: 10,
				Title:  "title",
				State:  "open",
				Body:   "body content",
			},
			want: fmt.Sprintf("\nIssue #%d\nTitle: %s\nState: %s\nBody:\n%s\n",
				10, "title", "open", "body content"),
			wantErr: nil,
		},
		{
			name:    "writer error",
			issue:   &domain.Issue{Number: 1, Title: "t", State: "s", Body: "b"},
			want:    "",
			wantErr: errors.New("write fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			var w ioWriter = &buf
			if tt.wantErr != nil {
				w = &errWriter{err: tt.wantErr}
			}

			// Act
			err := PrintIssue(w, tt.issue)

			// Assert
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("unexpected error: got %v want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got := buf.String(); got != tt.want {
				t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
