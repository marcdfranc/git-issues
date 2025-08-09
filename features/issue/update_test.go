package issue

import (
	"errors"
	"testing"

	"git-issues/domain"
	"git-issues/service/client"
	"git-issues/service/editor"
	"git-issues/testdata/stubs"
)

var defaultIssue = []byte("{\n  \"number\": 1,\n  \"title\": \"Example Issue\",\n  \"body\": \"This is a sample issue body.\",\n  \"html_url\": \"https://github.com/owner/repo/issues/1\"\n}")

func TestUpdate(t *testing.T) {
	f := UpdateFeature{config: &domain.Config{}}

	tests := []struct {
		name       string
		number     int
		clientStub client.GitHubClient
		editorStub editor.Editor
		wantErr    error
	}{
		{
			name:   "successful update",
			number: 1,
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(issue *domain.Issue) error {
					issue.Title = "Updated Title"
					issue.Body = "Updated Body"
					return nil
				},
			},
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return []byte(`{"number":1,"title":"Updated Title","body":"Updated Body"}`), nil
				},
			},
			wantErr: nil,
		},
		{
			name: "number required",
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(issue *domain.Issue) error {
					return nil
				},
			},
			clientStub: &stubs.ClientStub{},
			wantErr:    errNumberIsRequered,
		},
		{
			name:       "invalid issue",
			number:     3,
			editorStub: &stubs.EditorStub{},
			clientStub: &stubs.ClientStub{},
			wantErr:    errProcessing,
		},
		{
			name:   "error editor",
			number: 4,
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(issue *domain.Issue) error {
					return errProcessing
				},
			},
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return defaultIssue, nil
				},
			},
			wantErr: errUpdate,
		},
		{
			name:   "not found",
			number: 5,
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(issue *domain.Issue) error {
					issue.Title = "Title"
					issue.Body = "Body"
					return nil
				},
			},
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data *domain.Issue) ([]byte, error) {
					return nil, domain.ErrApi
				},
			},
			wantErr: errNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f.editor = tt.editorStub
			f.client = tt.clientStub
			err := f.Update(tt.number)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error got: %v, want: %v", err, tt.wantErr)
			}
		})
	}
}
