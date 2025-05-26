package main

import (
	"fmt"
	"strconv"

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
		err = featureConfig.Init()
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
	update := issue.NewUpdate(config, textEditor, serviceClient)
	list := issue.NewList(config, serviceClient)
	w := os.Stdout

	switch command {
	case "create":
		response, err = create.Create()
		if err != nil {
			fmt.Printf("error on create issue: %v\n", err)
			return
		}
		fmt.Println(response)
	case "list":
		issues, err := list.List()
		if err != nil {
			fmt.Printf("error on list issues: %v\n", err)
			return
		}
		err = issue.PrintIssues(w, issues)
		if err != nil {
			fmt.Printf("error on print issues: %v\n", err)
			return
		}

	case "update":
		if len(os.Args) < 3 {
			fmt.Println("please provide an issue number")
			return
		}

		number, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("please provide a valid issue number")
			return
		}

		err = update.Update(number)

		if err != nil {
			fmt.Printf("error on update issue: %v\n", err)
			return
		}
		fmt.Println("issue updated")

	default:
		fmt.Printf("command not found: %s\n", command)
		help.PrintHelp()
	}
}
