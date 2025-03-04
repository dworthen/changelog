package addcmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dworthen/changelog/internal/add"
	"github.com/dworthen/changelog/internal/utils"
	"github.com/dworthen/changelog/internal/versioninfo"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a changelog entry.",
	Long:  "Add a new changelog entry (.md file) in the .changelog directory to be applied at a later time. The entry describes the version bump to occur (patch|minor|major) and the changes associated with the bump.",
	Run: func(cmd *cobra.Command, args []string) {
		addModel, err := add.NewAddModel()
		utils.CheckError(err)
		_, err = tea.NewProgram(addModel, tea.WithAltScreen()).Run()
		utils.CheckError(err)
		err = versioninfo.PrintAvailableUpdate()
		utils.CheckError(err)
	},
}

func init() {
	// Here you will define your flags and configuration settings.
}
