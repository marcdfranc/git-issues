package main

import (
	"fmt"

	"git-issues/application"
	"git-issues/features/conf"
	"git-issues/features/help"
	"git-issues/features/issue"

	"os"
)

func main() {
	if len(os.Args) < 2 {
		help.PrintHelp()
		return
	}
	command := os.Args[1]

	if command == "init" {
		conf.InitConfig()
		return
	}

	config, err := application.LoadConfig()
	if err != nil {
		fmt.Printf("could not load conf: %v\n", err)
		fmt.Println("please run 'git_issues init' to configure.")
		return
	}

	switch command {
	case "create":
		issue.Create(config)
	case "list":
		issue.List(config)
	default:
		fmt.Printf("command not found: %s\n", command)
		help.PrintHelp()
	}
}
