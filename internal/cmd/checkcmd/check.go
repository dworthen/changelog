package checkcmd

import (
	"fmt"
	"os"

	"github.com/dworthen/changelog/internal/gitmanage"
	"github.com/dworthen/changelog/internal/utils"
	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check that the current branch or last commit contains a changelog entry. Useful for CI workflows to enforce the presence of changelog entries.",
	Long:  "Check that the current branch or last commit contains a changelog entry. Useful for CI workflows to enforce the presence of changelog entries.",
	Run: func(cmd *cobra.Command, args []string) {
		hasChangelogEntry, err := gitmanage.LastCommitContainsChangelogEntry()
		utils.CheckError(err)
		if !hasChangelogEntry {
			fmt.Fprintf(os.Stderr, "Last commit does not contain a changelog entry\n")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Success: Last commit contains a changelog entry\n")
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
