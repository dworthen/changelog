package initcmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dworthen/changelog/internal/initialize"
	"github.com/dworthen/changelog/internal/utils"
	"github.com/dworthen/changelog/internal/versioninfo"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize project to use changelog.",
	Long:  "Initialize project to use changelog. Creates a .changelog directory with a config.yaml file that configures how the changelog CLI applies and bumps changelog entries.",
	Run: func(cmd *cobra.Command, args []string) {
		initModel, err := initialize.NewInitializeModel()
		utils.CheckError(err)
		_, err = tea.NewProgram(initModel, tea.WithAltScreen()).Run()
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
