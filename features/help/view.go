package help

import "fmt"

func PrintHelp() {
	fmt.Println(`GitHub Issues CLI - Application to manage GitHub issues

Usage:
  ghissues <comand> [args]

Commands:
  init       conf the app
  create     Create a new issue
  list       List all issues
  view <n>   View the issue number n
  update <n> Update the issue number n
  close <n>  close the issue number n
  help       Display Help

Examples:
  ghissues init
  ghissues create
  ghissues list
  ghissues view 123
  ghissues update 123
  ghissues close 123`)
}
