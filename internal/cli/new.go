package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/azrakarakaya1/goscaffold/internal/generator"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// ProjectConfig holds project configuration
type ProjectConfig struct {
	Name             string
	ModulePath       string
	Template         string
	GitHubUser       string
	IncludeMake      bool
	IncludeDocker    bool
	IncludeCI        bool
	IncludeLint      bool
	IncludePreCommit bool
	IncludeTests     bool
	InitGit          bool
}

var config ProjectConfig
var allDevOps bool
var allQuality bool
var noInteractive bool

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Long: `Create a new Go project with the specified template and options.

Templates:
  basic    - Minimal Go project (default)
  cli      - CLI application with Cobra
  api      - REST API with Chi router
  grpc     - gRPC service with proto files
  library  - Reusable Go library

Examples:
  goscaffold new myapp
  goscaffold new myapi -t api -g username --all-devops
  goscaffold new mycli -t cli -g username -D -Q`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Template flags
	newCmd.Flags().StringVarP(&config.Template, "template", "t", "basic", "Project template (basic|cli|api|grpc|library)")
	newCmd.Flags().StringVarP(&config.GitHubUser, "github", "g", "", "GitHub username for module path")
	newCmd.Flags().StringVarP(&config.ModulePath, "module", "m", "", "Custom module path (overrides github)")

	// DevOps flags
	newCmd.Flags().BoolVar(&config.IncludeMake, "makefile", false, "Include Makefile")
	newCmd.Flags().BoolVar(&config.IncludeDocker, "docker", false, "Include Dockerfile and docker-compose")
	newCmd.Flags().BoolVar(&config.IncludeCI, "ci", false, "Include GitHub Actions CI workflow")
	newCmd.Flags().BoolVarP(&allDevOps, "all-devops", "D", false, "Include all DevOps files")

	// Quality flags
	newCmd.Flags().BoolVar(&config.IncludeLint, "lint", false, "Include golangci-lint config")
	newCmd.Flags().BoolVar(&config.IncludePreCommit, "precommit", false, "Include pre-commit hooks config")
	newCmd.Flags().BoolVar(&config.IncludeTests, "tests", false, "Include test file scaffolding")
	newCmd.Flags().BoolVarP(&allQuality, "all-quality", "Q", false, "Include all quality tools")

	// Other flags
	newCmd.Flags().BoolVar(&config.InitGit, "git", false, "Initialize git repository")
	newCmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "Skip interactive prompts")
}

func runNew(cmd *cobra.Command, args []string) error {
	success := color.New(color.FgGreen).SprintFunc()
	info := color.New(color.FgCyan).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()

	fmt.Println()
	fmt.Printf("  %s\n\n", info("goscaffold - Go Project Generator"))

	// Get project name
	if len(args) > 0 {
		config.Name = args[0]
	} else if !noInteractive {
		name, err := promptForInput("Project name", "myproject")
		if err != nil {
			return err
		}
		config.Name = name
	} else {
		return fmt.Errorf("project name is required")
	}

	// Validate project name
	if err := validateProjectName(config.Name); err != nil {
		return err
	}

	// Check if directory exists
	if _, err := os.Stat(config.Name); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", config.Name)
	}

	// Get GitHub username if not provided
	if config.ModulePath == "" && config.GitHubUser == "" && !noInteractive {
		username, err := promptForInput("GitHub username", "")
		if err != nil {
			return err
		}
		config.GitHubUser = username
	}

	// Build module path
	if config.ModulePath == "" {
		if config.GitHubUser != "" {
			config.ModulePath = fmt.Sprintf("github.com/%s/%s", config.GitHubUser, config.Name)
		} else {
			config.ModulePath = config.Name
		}
	}

	// Select template if not provided and interactive
	if !cmd.Flags().Changed("template") && !noInteractive {
		template, err := promptForTemplate()
		if err != nil {
			return err
		}
		config.Template = template
	}

	// Apply aggregate flags
	if allDevOps {
		config.IncludeMake = true
		config.IncludeDocker = true
		config.IncludeCI = true
	}
	if allQuality {
		config.IncludeLint = true
		config.IncludePreCommit = true
		config.IncludeTests = true
	}

	// Ask about DevOps if not specified and interactive
	if !noInteractive && !allDevOps && !cmd.Flags().Changed("makefile") && !cmd.Flags().Changed("docker") && !cmd.Flags().Changed("ci") {
		includeDevOps, err := promptForConfirm("Include DevOps files (Makefile, Docker, CI)?")
		if err != nil {
			return err
		}
		if includeDevOps {
			config.IncludeMake = true
			config.IncludeDocker = true
			config.IncludeCI = true
		}
	}

	// Ask about quality tools if not specified and interactive
	if !noInteractive && !allQuality && !cmd.Flags().Changed("lint") && !cmd.Flags().Changed("precommit") && !cmd.Flags().Changed("tests") {
		includeQuality, err := promptForConfirm("Include code quality tools (linter, pre-commit, tests)?")
		if err != nil {
			return err
		}
		if includeQuality {
			config.IncludeLint = true
			config.IncludePreCommit = true
			config.IncludeTests = true
		}
	}

	// Display configuration
	fmt.Printf("  %s %s\n", info("Project:"), config.Name)
	fmt.Printf("  %s %s\n", info("Module:"), config.ModulePath)
	fmt.Printf("  %s %s\n", info("Template:"), config.Template)
	fmt.Println()

	// Generate the project
	gen := generator.New(generator.Config{
		Name:             config.Name,
		ModulePath:       config.ModulePath,
		Template:         config.Template,
		IncludeMakefile:  config.IncludeMake,
		IncludeDocker:    config.IncludeDocker,
		IncludeCI:        config.IncludeCI,
		IncludeLint:      config.IncludeLint,
		IncludePreCommit: config.IncludePreCommit,
		IncludeTests:     config.IncludeTests,
		InitGit:          config.InitGit,
	})

	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Success message
	fmt.Printf("  %s Project '%s' created successfully!\n\n", success("✓"), config.Name)
	fmt.Printf("  %s\n", warn("Next steps:"))
	fmt.Printf("    cd %s\n", config.Name)
	fmt.Printf("    go mod tidy\n")
	if config.Template == "api" || config.Template == "grpc" {
		fmt.Printf("    go run ./cmd/%s\n", config.Name)
	} else if config.Template == "cli" {
		fmt.Printf("    go run ./cmd/%s\n", config.Name)
	} else {
		fmt.Printf("    go run .\n")
	}
	fmt.Println()

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for valid characters (alphanumeric, hyphen, underscore)
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, hyphens, or underscores")
	}

	// Check for reserved names
	reserved := []string{"internal", "pkg", "cmd", "vendor", "test", "main"}
	for _, r := range reserved {
		if strings.ToLower(name) == r {
			return fmt.Errorf("'%s' is a reserved name", name)
		}
	}

	return nil
}

func promptForInput(label, defaultVal string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultVal,
	}
	return prompt.Run()
}

func promptForConfirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}
	return strings.ToLower(result) == "y" || result == "", nil
}

func promptForTemplate() (string, error) {
	templates := []struct {
		Name        string
		Description string
	}{
		{"basic", "Minimal Go project"},
		{"cli", "CLI application with Cobra"},
		{"api", "REST API with Chi router"},
		{"grpc", "gRPC service"},
		{"library", "Reusable Go library"},
	}

	templateItems := make([]string, len(templates))
	for i, t := range templates {
		templateItems[i] = fmt.Sprintf("%s - %s", t.Name, t.Description)
	}

	prompt := promptui.Select{
		Label: "Select project template",
		Items: templateItems,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return templates[idx].Name, nil
}
