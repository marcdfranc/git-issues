package issue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git-issues/domain"
	"git-issues/util"
)

// MockGitHubAPI cria um servidor de teste para simular a API do GitHub
func MockGitHubAPI(statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}))
}

type MockEditor struct {
	Title string
	Body  string
	Error error
}

func (m *MockEditor) GetIssueContentFromEditor(config *domain.Config, defaultTitle, defaultBody string) (string, string, error) {
	return m.Title, m.Body, m.Error
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name           string
		mockEditor     *MockEditor
		apiStatusCode  int
		apiResponse    interface{}
		wantErr        bool
		expectedOutput string
	}{
		{
			name: "successful creation",
			mockEditor: &MockEditor{
				Title: "Test Issue",
				Body:  "Test Body",
			},
			apiStatusCode: http.StatusCreated,
			apiResponse: map[string]interface{}{
				"number":   1,
				"html_url": "http://example.com/issue/1",
				"title":    "Test Issue",
				"body":     "Test Body",
			},
			wantErr:        false,
			expectedOutput: "Issue created with success!\nNumber: 1\nURL: http://example.com/issue/1\n",
		},
		{
			name: "empty title",
			mockEditor: &MockEditor{
				Title: "",
				Body:  "Test Body",
			},
			wantErr:        true,
			expectedOutput: "title is required.\n",
		},
		{
			name: "editor error",
			mockEditor: &MockEditor{
				Error: fmt.Errorf("editor error"),
			},
			wantErr:        true,
			expectedOutput: "could not edit issue: editor error\n",
		},
		{
			name: "api error",
			mockEditor: &MockEditor{
				Title: "Test Issue",
				Body:  "Test Body",
			},
			apiStatusCode: http.StatusInternalServerError,
			apiResponse: map[string]interface{}{
				"message": "Internal server error",
			},
			wantErr:        true,
			expectedOutput: "Could not create issue:",
		},
		{
			name: "invalid api response",
			mockEditor: &MockEditor{
				Title: "Test Issue",
				Body:  "Test Body",
			},
			apiStatusCode:  http.StatusCreated,
			apiResponse:    "invalid json",
			wantErr:        true,
			expectedOutput: "error on process response:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Substitui a função original pelo mock
			oldEditor := util.GetIssueContentFromEditor
			defer func() { util.GetIssueContentFromEditor = oldEditor }()
			util.GetIssueContentFromEditor = tt.mockEditor.GetIssueContentFromEditor

			// Configura o servidor mock se necessário
			var server *httptest.Server
			if tt.apiStatusCode != 0 {
				server = MockGitHubAPI(tt.apiStatusCode, tt.apiResponse)
				defer server.Close()
			}

			// Configuração para o teste
			config := &domain.Config{
				Owner:      "test",
				Repo:       "repo",
				APIBaseURL: "http://example.com",
			}

			// Se tivermos um servidor mock, atualizamos a URL
			if server != nil {
				config.APIBaseURL = server.URL
			}

			// Captura a saída padrão
			output := captureOutput(func() {
				issue.Create(config)
			})

			// Verifica os resultados
			if tt.wantErr {
				if !contains(output, tt.expectedOutput) {
					t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
				}
			} else {
				if output != tt.expectedOutput {
					t.Errorf("Expected output %q, got %q", tt.expectedOutput, output)
				}
			}
		})
	}
}

// captureOutput captura a saída da função para teste
func captureOutput(f func()) string {
	rescue := fmt.Print
	defer func() { fmt.Print = rescue }()

	var output string
	fmt.Print = func(a ...interface{}) (n int, err error) {
		output = fmt.Sprint(a...)
		return len(output), nil
	}

	f()
	return output
}

// contains verifica se uma string contém outra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
