package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeFile writes content to a file, creating parent directories if needed
func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// ============================================================================
// Basic Template
// ============================================================================

func (g *Generator) createBasicTemplate() error {
	mainGo := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello from %s!")
}
`, g.config.Name)

	if err := writeFile(filepath.Join(g.config.Name, "main.go"), mainGo); err != nil {
		return err
	}

	if g.config.IncludeTests {
		testGo := `package main

import "testing"

func TestMain(t *testing.T) {
	// Add your tests here
}
`
		if err := writeFile(filepath.Join(g.config.Name, "main_test.go"), testGo); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// CLI Template (Cobra)
// ============================================================================

func (g *Generator) createCLITemplate() error {
	// Main entry point
	mainGo := fmt.Sprintf(`package main

import (
	"os"

	"%s/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
`, g.config.ModulePath)

	if err := writeFile(filepath.Join(g.config.Name, "cmd", g.config.Name, "main.go"), mainGo); err != nil {
		return err
	}

	// Root command
	rootCmd := fmt.Sprintf(`package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "%s",
	Short: "A brief description of your application",
	Long: `+"`"+`%s is a CLI application.

Add a longer description here.`+"`"+`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to %s!")
		fmt.Println("Use --help to see available commands.")
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}
`, g.config.Name, g.config.Name, g.config.Name)

	if err := writeFile(filepath.Join(g.config.Name, "internal", "cmd", "root.go"), rootCmd); err != nil {
		return err
	}

	// Version command
	versionCmd := `package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
`
	if err := writeFile(filepath.Join(g.config.Name, "internal", "cmd", "version.go"), versionCmd); err != nil {
		return err
	}

	return nil
}

// ============================================================================
// API Template (Chi Router)
// ============================================================================

func (g *Generator) createAPITemplate() error {
	// Main entry point
	mainGo := fmt.Sprintf(`package main

import (
	"log"
	"net/http"

	"%s/internal/router"
)

func main() {
	r := router.New()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
`, g.config.ModulePath)

	if err := writeFile(filepath.Join(g.config.Name, "cmd", g.config.Name, "main.go"), mainGo); err != nil {
		return err
	}

	// Router
	routerGo := fmt.Sprintf(`package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"%s/internal/handler"
	mw "%s/internal/middleware"
)

// New creates a new router with all routes configured
func New() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mw.ContentType)

	// Routes
	r.Get("/", handler.Home)
	r.Get("/health", handler.Health)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/hello", handler.Hello)
	})

	return r
}
`, g.config.ModulePath, g.config.ModulePath)

	if err := writeFile(filepath.Join(g.config.Name, "internal", "router", "router.go"), routerGo); err != nil {
		return err
	}

	// Handlers
	handlerGo := `package handler

import (
	"encoding/json"
	"net/http"
)

// Response is a generic API response
type Response struct {
	Message string ` + "`json:\"message\"`" + `
	Status  int    ` + "`json:\"status\"`" + `
}

// Home handles the root endpoint
func Home(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, Response{
		Message: "Welcome to the API",
		Status:  http.StatusOK,
	})
}

// Health handles health check endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, Response{
		Message: "OK",
		Status:  http.StatusOK,
	})
}

// Hello handles the hello endpoint
func Hello(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, Response{
		Message: "Hello, World!",
		Status:  http.StatusOK,
	})
}

func respond(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
`
	if err := writeFile(filepath.Join(g.config.Name, "internal", "handler", "handler.go"), handlerGo); err != nil {
		return err
	}

	// Middleware
	middlewareGo := `package middleware

import "net/http"

// ContentType sets the Content-Type header to application/json
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
`
	if err := writeFile(filepath.Join(g.config.Name, "internal", "middleware", "middleware.go"), middlewareGo); err != nil {
		return err
	}

	// Tests
	if g.config.IncludeTests {
		handlerTestGo := `package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
`
		if err := writeFile(filepath.Join(g.config.Name, "internal", "handler", "handler_test.go"), handlerTestGo); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// gRPC Template
// ============================================================================

func (g *Generator) createGRPCTemplate() error {
	// Main entry point
	mainGo := fmt.Sprintf(`package main

import (
	"log"
	"net"

	"%s/internal/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %%v", err)
	}

	s := grpc.NewServer()
	server.Register(s)

	log.Println("gRPC server starting on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %%v", err)
	}
}
`, g.config.ModulePath)

	if err := writeFile(filepath.Join(g.config.Name, "cmd", g.config.Name, "main.go"), mainGo); err != nil {
		return err
	}

	// Server
	serverGo := `package server

import (
	"context"

	"google.golang.org/grpc"
)

// GreeterServer implements the Greeter service
type GreeterServer struct{}

// Register registers the server with gRPC
func Register(s *grpc.Server) {
	// Register your gRPC services here
	// pb.RegisterGreeterServer(s, &GreeterServer{})
}

// SayHello implements the SayHello RPC
func (s *GreeterServer) SayHello(ctx context.Context, name string) (string, error) {
	return "Hello, " + name + "!", nil
}
`
	if err := writeFile(filepath.Join(g.config.Name, "internal", "server", "server.go"), serverGo); err != nil {
		return err
	}

	// Proto file
	protoFile := fmt.Sprintf(`syntax = "proto3";

package %s;

option go_package = "%s/pkg/pb";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
`, g.config.Name, g.config.ModulePath)

	if err := writeFile(filepath.Join(g.config.Name, "proto", g.config.Name+".proto"), protoFile); err != nil {
		return err
	}

	return nil
}

// ============================================================================
// Library Template
// ============================================================================

func (g *Generator) createLibraryTemplate() error {
	// Main library file
	libGo := fmt.Sprintf(`// Package %s provides functionality for...
package %s

// Version is the current version of the library
const Version = "0.1.0"

// Example is an example function
func Example() string {
	return "Hello from %s library!"
}
`, g.config.Name, g.config.Name, g.config.Name)

	if err := writeFile(filepath.Join(g.config.Name, "pkg", g.config.Name, g.config.Name+".go"), libGo); err != nil {
		return err
	}

	// Example usage
	exampleGo := fmt.Sprintf(`//go:build ignore

package main

import (
	"fmt"

	"%s/pkg/%s"
)

func main() {
	fmt.Println(%s.Example())
}
`, g.config.ModulePath, g.config.Name, g.config.Name)

	if err := writeFile(filepath.Join(g.config.Name, "examples", "basic", "main.go"), exampleGo); err != nil {
		return err
	}

	// Tests
	if g.config.IncludeTests {
		testGo := fmt.Sprintf(`package %s

import "testing"

func TestExample(t *testing.T) {
	result := Example()
	expected := "Hello from %s library!"

	if result != expected {
		t.Errorf("expected %%q, got %%q", expected, result)
	}
}
`, g.config.Name, g.config.Name)

		if err := writeFile(filepath.Join(g.config.Name, "pkg", g.config.Name, g.config.Name+"_test.go"), testGo); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// DevOps Files
// ============================================================================

func (g *Generator) createMakefile() error {
	fmt.Printf("  %s Creating Makefile...\n", g.info("→"))

	var runTarget string
	switch g.config.Template {
	case "basic":
		runTarget = "go run ."
	case "library":
		runTarget = "go run ./examples/basic"
	default:
		runTarget = fmt.Sprintf("go run ./cmd/%s", g.config.Name)
	}

	content := fmt.Sprintf(`# Project variables
BINARY_NAME=%s
PKG=%s

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test lint run tidy help

all: lint test build

## build: Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

## clean: Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

## test: Run tests
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

## lint: Run linter
lint:
	$(GOLINT) run ./...

## run: Run the application
run:
	%s

## tidy: Tidy dependencies
tidy:
	$(GOMOD) tidy

## help: Show this help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
`, g.config.Name, g.config.ModulePath, runTarget)

	return writeFile(filepath.Join(g.config.Name, "Makefile"), content)
}

func (g *Generator) createDockerFiles() error {
	fmt.Printf("  %s Creating Docker files...\n", g.info("→"))

	// Dockerfile
	dockerfile := fmt.Sprintf(`# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /%s ./cmd/%s

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /%s .

EXPOSE 8080

CMD ["./%s"]
`, g.config.Name, g.config.Name, g.config.Name, g.config.Name)

	if err := writeFile(filepath.Join(g.config.Name, "Dockerfile"), dockerfile); err != nil {
		return err
	}

	// docker-compose.yml
	compose := fmt.Sprintf(`version: '3.8'

services:
  %s:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=development
    restart: unless-stopped
`, g.config.Name)

	return writeFile(filepath.Join(g.config.Name, "docker-compose.yml"), compose)
}

func (g *Generator) createCIWorkflow() error {
	fmt.Printf("  %s Creating CI workflow...\n", g.info("→"))

	workflow := `name: CI

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'

    - name: Install dependencies
      run: go mod download

    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v4
      with:
        version: latest

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Build
      run: go build -v ./...
`
	return writeFile(filepath.Join(g.config.Name, ".github", "workflows", "ci.yml"), workflow)
}

// ============================================================================
// Code Quality Files
// ============================================================================

func (g *Generator) createLintConfig() error {
	fmt.Printf("  %s Creating linter config...\n", g.info("→"))

	content := `run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert

linters-settings:
  gofmt:
    simplify: true
  goimports:
    local-prefixes: github.com

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
`
	return writeFile(filepath.Join(g.config.Name, ".golangci.yml"), content)
}

func (g *Generator) createPreCommitConfig() error {
	fmt.Printf("  %s Creating pre-commit config...\n", g.info("→"))

	content := `repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: https://github.com/golangci/golangci-lint
    rev: v1.55.2
    hooks:
      - id: golangci-lint

  - repo: local
    hooks:
      - id: go-mod-tidy
        name: go mod tidy
        entry: go mod tidy
        language: system
        pass_filenames: false
`
	return writeFile(filepath.Join(g.config.Name, ".pre-commit-config.yaml"), content)
}

// ============================================================================
// README
// ============================================================================

func (g *Generator) createReadme() error {
	fmt.Printf("  %s Creating README...\n", g.info("→"))

	var description, usage string

	switch g.config.Template {
	case "cli":
		description = "A command-line application built with Go and Cobra."
		usage = fmt.Sprintf("```bash\ngo run ./cmd/%s\n```", g.config.Name)
	case "api":
		description = "A REST API built with Go and Chi router."
		usage = fmt.Sprintf("```bash\ngo run ./cmd/%s\n# Server starts on :8080\ncurl http://localhost:8080/health\n```", g.config.Name)
	case "grpc":
		description = "A gRPC service built with Go."
		usage = fmt.Sprintf("```bash\ngo run ./cmd/%s\n# Server starts on :50051\n```", g.config.Name)
	case "library":
		description = "A reusable Go library."
		usage = fmt.Sprintf("```go\nimport \"%s/pkg/%s\"\n\nfunc main() {\n    result := %s.Example()\n}\n```", g.config.ModulePath, g.config.Name, g.config.Name)
	default:
		description = "A Go project."
		usage = "```bash\ngo run .\n```"
	}

	content := fmt.Sprintf(`# %s

%s

## Installation

`+"```bash"+`
go get %s
`+"```"+`

## Usage

%s

## Development

### Prerequisites

- Go 1.21 or later
`, g.config.Name, description, g.config.ModulePath, usage)

	if g.config.IncludeMakefile {
		content += `
### Available Commands

` + "```bash" + `
make help    # Show available commands
make build   # Build the binary
make test    # Run tests
make lint    # Run linter
make run     # Run the application
` + "```" + `
`
	}

	content += `
## License

MIT License
`

	return writeFile(filepath.Join(g.config.Name, "README.md"), content)
}
