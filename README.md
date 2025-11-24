# git-issues

A CLI client to manage GitHub issues.  
Exercise from "The Go Programming Language" (chapter 4, exercise 11). Provides commands to create, list, view, update and close issues.

![Go Version](https://img.shields.io/badge/Go-1.16%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## Table of Contents

- [Description](#description)
- [Prerequisites](#prerequisites)
- [Install & Build](#install--build)
- [Configuration](#configuration)
- [Usage](#usage)
- [Testing](#testing)
- [Repository Layout](#repository-layout)
- [Contributing](#contributing)
- [License](#license)
- [Maintainers](#maintainers)
- [Troubleshooting](#troubleshooting)


## Description

`git-issues` is a small command-line application written in Go to perform CRUD operations on GitHub issues. It is intended as a learning exercise and a lightweight tool to manage issues from the terminal.

## Prerequisites

- Go 1.16 or newer (recommended to use the latest stable Go)
- Git
- A GitHub personal access token with repository permissions
- On Windows: use PowerShell or CMD; ensure your editor command (e.g. `code`, `notepad`) is available in PATH

Verify Go is installed:

```bash
go version
```

## Install & Build

Clone and prepare the project:

```bash
git clone https://github.com/marcdfranc/git-issues.git
```

Build the binary:

```bash
go build -o ghissues
```

Run without building:

```bash
go run main.go
```

## Configuration

The `init` command interactively collects configuration (GitHub token, repository owner, repository name, preferred editor) and saves it to a configuration file. You can also create a `.ghissuescli` file manually.

```bash
./ghissues init
```

Example `.ghissuescli`:
```json
{
  "token": "YOUR_GITHUB_TOKEN",
  "owner": "repo-owner",
  "repo": "repo-name",
  "editor": "code"
}
```
- ***token:*** GitHub personal access token (or set via environment variable if supported).
- ***owner:*** repository owner or organization (only the username, don't use the complete email).
- ***repo:*** repository name.
- ***editor:*** command used to open the editor for issue title/body (e.g. code, notepad, vim).

Place `.ghissuescli` in the working directory or the path expected by the application.

## Usage

Commands:

- `init`: Configure the application
- `create`: Creates a new issue (opens the editor to write title and body) 
- `list`: Lists all issues
- `view <number>`: Shows the details of a specific issue
- `update <number>`: Updates an existing issue
- `close <number>`: Closes an issue

example:

```bash
./ghissues init
```

```bash
./ghissues create
```

```bash
./ghissues list
```

```bash
./ghissues view 12
```

```bash
./ghissues close 12
```
sample of `list` output:

```text
#12 Fix login bug (open)
#11 Add dark mode (closed)
#10 Improve README (open)
```


## Testing

Run tests with:

```bash
go test ./...
```

## Repository Layout

```
git-issues
│   .gitignore
│   Dockerfile
│   ghissues
│   go.mod
│   LICENSE
│   main.go
│   README.md
│
├───application
│       config.go
│       config_test.go
│       
├───domain
│       config.go
│       errors.go
│       issues.go
│       
├───features
│   ├───conf
│   │       init.go
│   │       init_test.go
│   │       
│   ├───help
│   │       view.go
│   │       
│   └───issue
│           close.go
│           close_test.go
│           common.go
│           create.go
│           create_test.go
│           list.go
│           list_test.go
│           print.go
│           print_test.go
│           update.go
│           update_test.go
│           view.go
│           view_test.go
│           
├───service
│   ├───client
│   │       github.go
│   │       github_test.go
│   │       
│   └───editor
│           editor.go
│           editor_test.go
│           
└───testdata
    ├───data
    │       data.go
    │       
    └───stubs
            serviceclient.go
            serviceeditor.go
```

## Contributing

- Open issues for bugs and feature requests.
- Use branches for features and fixes, then open PRs.
- Add tests for new features and follow project style.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Maintainers

- Maintainer: Marcelo de Oliveira Francisco (`marcdfranc@gmail.com` / `@marcdfranc` on GitHub)

## Troubleshooting

- Invalid token: Check if `.ghissues` was created and contains a valid token.
- Permission errors: Ensure the token is correctly scoped to the target repository.
- Editor not found: configure the editor in `.ghissues` to a command available in PATH (Windows: `notepad` or `code`), prefer to use the application's init command instead of directly editing the file.
