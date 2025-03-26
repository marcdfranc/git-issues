package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"git-issues/domain"
	"os"
	"strings"
)

const configFile = ".ghissuescli"

func InitConfig() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("GitHub Personal Access Token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	fmt.Print("Repository Owner (username/organization): ")
	owner, _ := reader.ReadString('\n')
	owner = strings.TrimSpace(owner)

	fmt.Print("Repository Name: ")
	repo, _ := reader.ReadString('\n')
	repo = strings.TrimSpace(repo)

	fmt.Print("Default text editor (empty for system default): ")
	editor, _ := reader.ReadString('\n')
	editor = strings.TrimSpace(editor)

	config := domain.Config{
		Token:      token,
		Owner:      owner,
		Repo:       repo,
		Editor:     editor,
		APIBaseURL: "https://api.github.com",
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("could not generat config: %v\n", err)
		return
	}

	err = os.WriteFile(configFile, configData, 0600)
	if err != nil {
		fmt.Printf("could not save config: %v\n", err)
		return
	}

	fmt.Println("config created with success!")
}

func LoadConfig() (*domain.Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config domain.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
