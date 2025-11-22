package issue

import (
	"errors"
	"testing"

	"git-issues/domain"
	"git-issues/testdata/stubs"
)

func TestListFeature(t *testing.T) {
	// common config used by tests
	cfg := &domain.Config{
		APIBaseURL: "https://api.example.com",
		Owner:      "owner",
		Repo:       "repo",
	}

	fetchErr := errors.New("network")

	tests := []struct {
		name       string
		clientStub *stubs.ClientStub
		wantErr    error
		wantLen    int
		wantIssue  *domain.Issue
	}{
		{
			name: "successful list",
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return []byte(`[{"number":1,"state":"open","title":"t","body":"b"}]`), nil
				},
			},
			wantErr: nil,
			wantLen: 1,
			wantIssue: &domain.Issue{
				Number: 1,
				State:  "open",
				Title:  "t",
				Body:   "b",
			},
		},
		{
			name: "request error forwarded",
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return nil, fetchErr
				},
			},
			wantErr: fetchErr,
		},
		{
			name: "invalid json -> processing error",
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return []byte(`{ not json`), nil
				},
			},
			wantErr: errProcessing,
		},
		{
			name: "empty list",
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return []byte(`[]`), nil
				},
			},
			wantErr: nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			f := NewList(cfg, tt.clientStub)

			// Act
			got, err := f.List()

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
			if len(got) != tt.wantLen {
				t.Fatalf("unexpected length: got %d want %d", len(got), tt.wantLen)
			}

			if tt.wantLen > 0 && tt.wantIssue != nil {
				if got[0].Number != tt.wantIssue.Number {
					t.Fatalf("Number mismatch: got %d want %d", got[0].Number, tt.wantIssue.Number)
				}
				if got[0].State != tt.wantIssue.State {
					t.Fatalf("State mismatch: got %q want %q", got[0].State, tt.wantIssue.State)
				}
				if got[0].Title != tt.wantIssue.Title {
					t.Fatalf("Title mismatch: got %q want %q", got[0].Title, tt.wantIssue.Title)
				}
				if got[0].Body != tt.wantIssue.Body {
					t.Fatalf("Body mismatch: got %q want %q", got[0].Body, tt.wantIssue.Body)
				}
			}
		})
	}
}
