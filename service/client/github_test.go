package client

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git-issues/domain"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestMakeGitHubRequest_Success(t *testing.T) {
	// Arrange
	want := []byte(`{"message":"success"}`)

	// Mock HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(want)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	service := New(defaultConfig)

	data := map[string]string{"example": "value"}

	// Act
	got, err := service.MakeGitHubRequest(http.MethodPost, server.URL, data)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("unexpected result got %q, want: %q", got, want)
	}
}

func TestMakeGitHubRequest_Errors(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad request"}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	type args struct {
		method string
		url    string
		data   interface{}
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "when malformed data received",
			args: args{
				method: http.MethodPost,
				url:    "http://example.com",
				data:   make(chan int),
			},
			want: errEncoding,
		},
		{
			name: "when malformed request received",
			args: args{
				method: http.MethodGet,
				url:    "http://invalid.url", // malformed url
				data:   nil,
			},
			want: errRequest,
		},
		{
			name: "when github api is down",
			args: args{
				method: http.MethodGet,
				url:    server.URL, // malformed url
				data:   nil,
			},
			want: errApi,
		},
	}

	service := New(defaultConfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, got := service.MakeGitHubRequest(tt.args.method, tt.args.url, tt.args.data)

			// Assert
			if !errors.Is(got, tt.want) {
				t.Errorf("unexpected error got %v, want: %v", got, tt.want)
			}
		})
	}

}

var (
	defaultConfig = &domain.Config{Token: "mockToken"}
)
