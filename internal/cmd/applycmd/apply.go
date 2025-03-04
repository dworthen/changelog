package applycmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dworthen/changelog/internal/apply"
	"github.com/dworthen/changelog/internal/globals"
	"github.com/dworthen/changelog/internal/utils"
	"github.com/dworthen/changelog/internal/versioninfo"
	"github.com/spf13/cobra"
)

var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply changelog entries.",
	Long:  "Apply changelog entries (.md files) in the .changelog directory to the current working directory. The entries describe the version bump to occur (patch|minor|major) and the changes associated with the bump. The final version will be the max version bump present across all changelog entries.",
	Run: func(cmd *cobra.Command, args []string) {
		applyModel, err := apply.NewApplyModel()
		utils.CheckError(err)

		globals.Program = tea.NewProgram(applyModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
		_, err = globals.Program.Run()
		utils.CheckError(err)

		err = versioninfo.PrintAvailableUpdate()
		utils.CheckError(err)
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
