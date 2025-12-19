# goscaffold

[![CI](https://github.com/azrakarakaya1/goscaffold/actions/workflows/ci.yml/badge.svg)](https://github.com/azrakarakaya1/goscaffold/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/azrakarakaya1/goscaffold)](https://goreportcard.com/report/github.com/azrakarakaya1/goscaffold)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)

A powerful CLI tool to scaffold Go projects with best-practice directory structures, templates, and DevOps configurations.

## Features

- **Multiple Project Templates**
  - `basic` - Minimal Go project
  - `cli` - CLI application with Cobra
  - `api` - REST API with Chi router
  - `grpc` - gRPC service with proto files
  - `library` - Reusable Go library

- **DevOps Integration**
  - Makefile with common targets
  - Dockerfile (multi-stage build)
  - docker-compose.yml
  - GitHub Actions CI workflow

- **Code Quality Tools**
  - golangci-lint configuration
  - pre-commit hooks
  - Test scaffolding

- **Interactive Mode** - Guided setup with sensible defaults
- **Non-Interactive Mode** - Perfect for automation and CI/CD

## Installation

### Using Go

```bash
go install github.com/azrakarakaya1/goscaffold@latest
```

### From Releases

Download the binary for your platform from the [releases page](https://github.com/azrakarakaya1/goscaffold/releases).

### From Source

```bash
git clone https://github.com/azrakarakaya1/goscaffold.git
cd goscaffold
go build -o goscaffold ./cmd/goscaffold
```

## Usage

### Interactive Mode (Recommended for beginners)

```bash
goscaffold new
```

You'll be prompted for:
- Project name
- GitHub username
- Template selection
- DevOps and quality tool options

### Command Line Mode

```bash
# Basic project
goscaffold new myapp -g yourusername

# REST API with all DevOps files
goscaffold new myapi -t api -g yourusername --all-devops

# CLI app with all quality tools
goscaffold new mycli -t cli -g yourusername --all-quality

# Full-featured project
goscaffold new myproject -t api -g yourusername -D -Q --git
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--template` | `-t` | Project template (basic\|cli\|api\|grpc\|library) |
| `--github` | `-g` | GitHub username for module path |
| `--module` | `-m` | Custom module path (overrides --github) |
| `--makefile` | | Include Makefile |
| `--docker` | | Include Dockerfile and docker-compose |
| `--ci` | | Include GitHub Actions CI workflow |
| `--all-devops` | `-D` | Include all DevOps files |
| `--lint` | | Include golangci-lint config |
| `--precommit` | | Include pre-commit hooks config |
| `--tests` | | Include test file scaffolding |
| `--all-quality` | `-Q` | Include all quality tools |
| `--git` | | Initialize git repository |
| `--no-interactive` | | Skip interactive prompts |

## Generated Project Structure

### API Template Example

```
myapi/
├── cmd/
│   └── myapi/
│       └── main.go
├── internal/
│   ├── handler/
│   │   └── handler.go
│   ├── middleware/
│   │   └── middleware.go
│   └── router/
│       └── router.go
├── pkg/
├── .github/
│   └── workflows/
│       └── ci.yml
├── go.mod
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── .golangci.yml
├── .pre-commit-config.yaml
├── .gitignore
└── README.md
```

## Examples

### Create a REST API

```bash
goscaffold new userservice -t api -g myusername -D -Q

cd userservice
go mod tidy
go run ./cmd/userservice

# In another terminal
curl http://localhost:8080/health
```

### Create a CLI Tool

```bash
goscaffold new mytool -t cli -g myusername -D -Q

cd mytool
go mod tidy
go run ./cmd/mytool --help
```

### Create a Library

```bash
goscaffold new mylib -t library -g myusername -Q

cd mylib
go test ./...
```

## Development

### Prerequisites

- Go 1.21 or later

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Chi](https://github.com/go-chi/chi) - HTTP router
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts
- [color](https://github.com/fatih/color) - Colored terminal output
