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
		if _, err := w.Write(want); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	service := New(defaultConfig)

	data := domain.Issue{
		Number: 0,
		Title:  "",
		Body:   "",
		State:  "",
	}

	// Act
	got, err := service.MakeRequest(http.MethodPost, server.URL, &data)

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
		if _, err := w.Write([]byte(`{"message":"bad request"}`)); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	type args struct {
		method string
		url    string
		data   *domain.Issue
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "when malformed request received",
			args: args{
				method: http.MethodGet,
				url:    "http://invalid.url", // malformed url
				data:   nil,
			},
			want: domain.ErrRequest,
		},
		{
			name: "when github api is down",
			args: args{
				method: http.MethodGet,
				url:    server.URL, // malformed url
				data:   nil,
			},
			want: domain.ErrApi,
		},
	}

	service := New(defaultConfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, got := service.MakeRequest(tt.args.method, tt.args.url, tt.args.data)

			// Assert
			if !errors.Is(got, tt.want) {
				t.Errorf("unexpected error got %v, want: %v", got, tt.want)
			}
		})
	}

}

func TestMakeGitHubRequest_CreateRequestError(t *testing.T) {
	service := New(defaultConfig)

	// invalid method to trigger request creation error
	_, err := service.MakeRequest("\x00BAD", "http://example.com", nil)
	if !errors.Is(err, domain.ErrCreateRequest) {
		t.Errorf("expected ErrCreateRequest, got %v", err)
	}
}

func TestMakeGitHubRequest_ReadBodyError(t *testing.T) {
	service := New(defaultConfig)

	// Mock HTTP server that closes connection abruptly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		err := conn.Close() // close connection immediately
		if err != nil {
			t.Fatalf("failed to close connection: %v", err)
		}
	}))
	defer server.Close()

	_, err := service.MakeRequest(http.MethodGet, server.URL, nil)
	if err == nil {
		t.Errorf("expected error on reading body, got nil")
	}
}

var (
	defaultConfig = &domain.Config{Token: "mockToken"}
)
