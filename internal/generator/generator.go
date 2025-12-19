package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

// Config holds the project generation configuration
type Config struct {
	Name             string
	ModulePath       string
	Template         string
	IncludeMakefile  bool
	IncludeDocker    bool
	IncludeCI        bool
	IncludeLint      bool
	IncludePreCommit bool
	IncludeTests     bool
	InitGit          bool
}

// Generator handles project generation
type Generator struct {
	config Config
	info   func(a ...interface{}) string
}

// New creates a new Generator
func New(cfg Config) *Generator {
	return &Generator{
		config: cfg,
		info:   color.New(color.FgCyan).SprintFunc(),
	}
}

// Generate creates the project
func (g *Generator) Generate() error {
	// Create base directories
	if err := g.createDirectories(); err != nil {
		return err
	}

	// Create go.mod
	if err := g.createGoMod(); err != nil {
		return err
	}

	// Create template-specific files
	if err := g.createTemplateFiles(); err != nil {
		return err
	}

	// Create .gitignore
	if err := g.createGitignore(); err != nil {
		return err
	}

	// Create DevOps files if requested
	if g.config.IncludeMakefile {
		if err := g.createMakefile(); err != nil {
			return err
		}
	}

	if g.config.IncludeDocker {
		if err := g.createDockerFiles(); err != nil {
			return err
		}
	}

	if g.config.IncludeCI {
		if err := g.createCIWorkflow(); err != nil {
			return err
		}
	}

	// Create quality files if requested
	if g.config.IncludeLint {
		if err := g.createLintConfig(); err != nil {
			return err
		}
	}

	if g.config.IncludePreCommit {
		if err := g.createPreCommitConfig(); err != nil {
			return err
		}
	}

	// Create README
	if err := g.createReadme(); err != nil {
		return err
	}

	// Initialize git if requested
	if g.config.InitGit {
		if err := g.initGit(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) createDirectories() error {
	fmt.Printf("  %s Creating directories...\n", g.info("→"))

	dirs := []string{
		g.config.Name,
	}

	// Add template-specific directories
	switch g.config.Template {
	case "basic":
		// No additional directories needed
	case "cli":
		dirs = append(dirs,
			filepath.Join(g.config.Name, "cmd", g.config.Name),
			filepath.Join(g.config.Name, "internal"),
		)
	case "api":
		dirs = append(dirs,
			filepath.Join(g.config.Name, "cmd", g.config.Name),
			filepath.Join(g.config.Name, "internal", "handler"),
			filepath.Join(g.config.Name, "internal", "middleware"),
			filepath.Join(g.config.Name, "internal", "router"),
			filepath.Join(g.config.Name, "pkg"),
		)
	case "grpc":
		dirs = append(dirs,
			filepath.Join(g.config.Name, "cmd", g.config.Name),
			filepath.Join(g.config.Name, "internal", "server"),
			filepath.Join(g.config.Name, "proto"),
			filepath.Join(g.config.Name, "pkg"),
		)
	case "library":
		dirs = append(dirs,
			filepath.Join(g.config.Name, "pkg", g.config.Name),
			filepath.Join(g.config.Name, "examples"),
		)
	default:
		dirs = append(dirs,
			filepath.Join(g.config.Name, "cmd"),
			filepath.Join(g.config.Name, "internal"),
			filepath.Join(g.config.Name, "pkg"),
		)
	}

	// Add CI directory if needed
	if g.config.IncludeCI {
		dirs = append(dirs, filepath.Join(g.config.Name, ".github", "workflows"))
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (g *Generator) createGoMod() error {
	fmt.Printf("  %s Creating go.mod...\n", g.info("→"))

	content := fmt.Sprintf(`module %s

go 1.21
`, g.config.ModulePath)

	return writeFile(filepath.Join(g.config.Name, "go.mod"), content)
}

func (g *Generator) createTemplateFiles() error {
	fmt.Printf("  %s Creating template files...\n", g.info("→"))

	switch g.config.Template {
	case "basic":
		return g.createBasicTemplate()
	case "cli":
		return g.createCLITemplate()
	case "api":
		return g.createAPITemplate()
	case "grpc":
		return g.createGRPCTemplate()
	case "library":
		return g.createLibraryTemplate()
	default:
		return g.createBasicTemplate()
	}
}

func (g *Generator) createGitignore() error {
	content := `# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Dependency directories
vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local

# Build
dist/

# Logs
*.log
`
	return writeFile(filepath.Join(g.config.Name, ".gitignore"), content)
}

func (g *Generator) initGit() error {
	fmt.Printf("  %s Initializing git repository...\n", g.info("→"))

	cmd := exec.Command("git", "init")
	cmd.Dir = g.config.Name
	return cmd.Run()
}
