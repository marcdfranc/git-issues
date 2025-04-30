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
	c := conf.New()

	if len(os.Args) < 2 {
		help.PrintHelp()
		return
	}
	command := os.Args[1]

	if command == "init" {
		c.Init()
		return
	}

	config, err := c.GetConfig()
	if err != nil {
		fmt.Println(err)
		fmt.Println("please run 'git_issues init' to configure.")
		return
	}

	switch command {
	case "list":
		issue.ListFeature(config)
	case "create":
		issue.Create(config)
	default:
		fmt.Printf("command not found: %s\n", command)
		help.PrintHelp()
	}
}
