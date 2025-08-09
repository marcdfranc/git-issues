# git-issues

This repository implements a Git issues client.

This is an exercise proposed in "Go, The Programing Language" cap 4, exercise 11.

The idea is creating a CRUD for GitHub issues in Go Lang

## Prerequisites

1. Go installed (version 1.16 or higher recommended)
2. GitHub personal access token with permissions to repositories
3. Text editor configured (or will use the system default)

## Compile

```bash
go build -o ghissues
```

## Configure application
This will ask for your GitHub access token, repository owner name, repository name, and preferred text editor.
This feature is completed and can be used as follows.

```bash
./ghissues init
```

## Commands

- `init`: Configure the application
- `create`: Creates a new issue (opens the editor to write title and body) 
- `list`: Lists all issues
- `view <number>`: Shows the details of a specific issue
- `update <number>`: Updates an existing issue
- `close <number>`: Closes an issue