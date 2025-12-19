package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	versionStr string
	commitStr  string
	dateStr    string
)

// SetVersionInfo sets the version information from main
func SetVersionInfo(version, commit, date string) {
	versionStr = version
	commitStr = commit
	dateStr = date
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		cyan := color.New(color.FgCyan).SprintFunc()
		white := color.New(color.FgWhite).SprintFunc()

		fmt.Println()
		fmt.Printf("  %s %s\n", cyan("goscaffold"), white(versionStr))
		fmt.Printf("  %s   %s\n", cyan("commit:"), white(commitStr))
		fmt.Printf("  %s    %s\n", cyan("built:"), white(dateStr))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
