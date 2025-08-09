package issue

import (
	"errors"
	"testing"

	"git-issues/domain"
	"git-issues/service/client"
	"git-issues/service/editor"
	"git-issues/testdata/stubs"
)

func TestViewIssue(t *testing.T) {
	tests := []struct {
		name       string
		number     int
		clientStub client.GitHubClient
		editorStub editor.Editor
		wantErr    error
	}{
		{
			name:   "success",
			number: 1,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return []byte(`{"number":1,"title":"Test Issue","body":"Test Body"}`), nil
				},
			},
			editorStub: &stubs.EditorStub{},
			wantErr:    nil,
		},
		{
			name:   "not found",
			number: 2,
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return nil, domain.ErrApi
				},
			},
			editorStub: &stubs.EditorStub{},
			wantErr:    domain.ErrApi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ViewFeature{
				config: &domain.Config{},
				client: tt.clientStub,
			}
			_, err := f.View(tt.number)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}
