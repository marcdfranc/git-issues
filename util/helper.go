package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"git-issues/domain"
)

var goos = runtime.GOOS

func GetIssueContentFromEditor(config *domain.Config, initialTitle, initialBody string) (string, string, error) {
	tempFile, err := os.CreateTemp("", "ghissue-*.md")
	if err != nil {
		return "", "", fmt.Errorf("could not create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := fmt.Sprintf("%s\n\n%s", initialTitle, initialBody)
	_, err = tempFile.WriteString(content)
	if err != nil {
		return "", "", fmt.Errorf("could not write temp file: %v", err)
	}
	tempFile.Close()

	editor := GetEditor(config)
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("error on exec text editor: %v", err)
	}

	editedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", "", fmt.Errorf("error on read file: %v", err)
	}

	parts := strings.SplitN(string(editedContent), "\n", 2)
	title := strings.TrimSpace(parts[0])
	body := ""
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}

	return title, body, nil
}

func GetEditor(config *domain.Config) string {
	if config.Editor != "" {
		return config.Editor
	}

	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	// Default editors per OS
	if goos == "windows" {
		return "notepad"
	}
	return "vi"
}
