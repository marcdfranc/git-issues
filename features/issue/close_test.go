package issue

import (
	"errors"
	"testing"

	"git-issues/domain"
	"git-issues/testdata/stubs"
)

func TestCloseFeature(t *testing.T) {
	// common config used by tests
	cfg := &domain.Config{
		APIBaseURL: "https://api.example.com",
		Owner:      "owner",
		Repo:       "repo",
	}

	patchErr := errors.New("patch error")

	tests := []struct {
		name       string
		number     int
		clientStub *stubs.ClientStub
		want       error
	}{
		{
			name:   "successful close",
			number: 1,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					if method == "GET" {
						// return an open issue
						return []byte(`{"number":1,"state":"open","title":"t","body":"b"}`), nil
					}
					// PATCH returns closed issue
					return []byte(`{"number":1,"state":"closed","title":"t","body":"b"}`), nil
				},
			},
			want: nil,
		},
		{
			name:   "number is zero",
			number: 0,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return nil, nil
				},
			},
			want: errNumberIsRequered,
		},
		{
			name:   "get request error -> not found",
			number: 2,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					if method == "GET" {
						return nil, errors.New("network")
					}
					return nil, nil
				},
			},
			want: errNotFound,
		},
		{
			name:   "invalid json on GET -> processing error",
			number: 3,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					if method == "GET" {
						// invalid JSON
						return []byte(`{ invalid json `), nil
					}
					return nil, nil
				},
			},
			want: errProcessing,
		},
		{
			name:   "patch returns error -> forwarded",
			number: 4,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					if method == "GET" {
						return []byte(`{"number":4,"state":"open","title":"t","body":"b"}`), nil
					}
					return nil, patchErr
				},
			},
			want: patchErr,
		},
		{
			name:   "invalid json on PATCH response -> processing error",
			number: 5,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					if method == "GET" {
						return []byte(`{"number":5,"state":"open","title":"t","body":"b"}`), nil
					}
					// PATCH returns invalid JSON
					return []byte(`{ not json`), nil
				},
			},
			want: errProcessing,
		},
		{
			name:   "patch response not closed -> errClose",
			number: 6,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					if method == "GET" {
						return []byte(`{"number":6,"state":"open","title":"t","body":"b"}`), nil
					}
					// PATCH returns issue still open
					return []byte(`{"number":6,"state":"open","title":"t","body":"b"}`), nil
				},
			},
			want: errClose,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			f := CloseFeature{
				config: cfg,
				client: tt.clientStub,
			}

			// Act
			err := f.Close(tt.number)

			// Assert
			if !errors.Is(err, tt.want) {
				t.Fatalf("unexpected error: got %v want %v", err, tt.want)
			}
		})
	}
}
