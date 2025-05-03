package main

import (
	"fmt"

	"git-issues/application"
	"git-issues/features/conf"
	"git-issues/features/help"
	"git-issues/features/issue"
	"git-issues/service/client"
	"git-issues/service/editor"

	"os"
)

func main() {
	if len(os.Args) < 2 {
		help.PrintHelp()
		return
	}

	var err error
	featureConfig := conf.New()

	command := os.Args[1]

	if command == "init" {
		err = featureConfig.InitConfig()
		if err != nil {
			fmt.Printf("error on start the application: %v\n", err)
		}
		return
	}

	config, err := application.LoadConfig()
	if err != nil {
		fmt.Printf("could not load conf: %v\n", err)
		fmt.Println("please run 'git_issues init' to configure.")
		return
	}

	var response string

	textEditor := editor.New(config)
	serviceClient := client.New(config)
	create := issue.NewCreate(config, textEditor, serviceClient)

	switch command {
	case "create":
		response, err = create.Create()
		if err != nil {
			fmt.Printf("error on create issue: %v\n", err)
		}
		fmt.Println(response)
	case "list":
		issue.List(config)
	default:
		fmt.Printf("command not found: %s\n", command)
		help.PrintHelp()
	}
}
