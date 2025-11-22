# git-issues

A CLI client to manage GitHub issues.  
Exercise from "The Go Programming Language" (chapter 4, exercise 11). Provides commands to create, list, view, update and close issues.


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

or, if entrypoint is in cmd:

```bachh
go run ./cmd/ghissues
```

## Configuration

The `init` command interactively collects configuration (GitHub token, repository owner, repository name, preferred editor) and saves it to a configuration file. You can also create a `config.json` file manually.

```bash
./ghissues init
```

Example `config.json`:
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

Place `config.json` in the working directory or the path expected by the application.

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
./ghissues create
./ghissues list
./ghissues view 12
./ghissues close 12
```

## Testing

Run tests with:

```bash
go test ./...
```

## Repository Layout

```
git-issues/
├── application/              
│   └── config.go
├── domain/           
│   ├── config.go         # Configuration definitions
│   ├── errors.go         # errors definitions
│   └── issue.go          # Issue definitions
├── features/
│   ├──conf/
│   │  └── init.go  # Configuration service
│   │  └── init_test.go # Configuration service tests
│   ├──help/
│   │  └── view.go  # Help service
│   │  └── help_test.go # Help service tests
│   ├──issue/
│   │  └── create.go      # Create issue service
│   │  └── create_test.go # Create issue service tests
│   │  └── list.go        # List issues service
│   │  └── list_test.go   # List issues service tests
│   │  └── view.go        # View issue service
│   │  └── view_test.go   # View issue service tests
│   │  └── update.go      # Update issue service
│   │  └── update_test.go # Update issue service tests
│   │  └── close.go       # Close issue service
│   │  └── close_test.go  # Close issue service tests
├── service/
├── tests/              # Test data and utilities
├── go.mod              # Go module file
├── go.sum              # Go module checksums
└── README.md           # Project documentation
```

## Contributing

- Open issues for bugs and feature requests.
- Use branches for features and fixes, then open PRs.
- Add tests for new features and follow project style.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Maintainers

- Maintainer: Marcelo de Oliveira Francisco (`marcdfranc@gmail.com` ou `@marcdfranc` no GitHub)

## Troubleshooting

- Invalid token: verify `config.json` or environment variable contains a valid token.
- Permission errors: ensure token has correct scopes for the target repository.
- Editor not found: set editor in `config.json` to a command available in PATH (Windows: `notepad` or `code`).

