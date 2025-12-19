package cli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goscaffold",
	Short: "A CLI tool to scaffold Go projects",
	Long: `goscaffold is a powerful CLI tool that helps you create new Go projects
with best-practice directory structures, templates, and DevOps configurations.

Features:
  - Multiple project templates (basic, cli, api, grpc, library)
  - DevOps files (Makefile, Dockerfile, CI workflows)
  - Code quality tools (linter configs, pre-commit hooks)
  - Interactive mode with sensible defaults

Example:
  goscaffold new myproject -t api -g yourusername --all-devops`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand provided, show help
		cmd.Help()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add color to the CLI
	cobra.AddTemplateFunc("cyan", color.CyanString)
	cobra.AddTemplateFunc("green", color.GreenString)
	cobra.AddTemplateFunc("yellow", color.YellowString)
}
