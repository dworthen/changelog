package updatecmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dworthen/changelog/internal/utils"
	"github.com/dworthen/changelog/internal/versioninfo"

	"github.com/dworthen/updater"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to the latest version of changelog",
	Long:  "Update to the latest version of changelog",
	Run: func(cmd *cobra.Command, args []string) {
		currentVersion, err := versioninfo.GetVersion()
		utils.CheckError(err)
		isUpdate, newVersion, err := versioninfo.CheckForUpdate()
		utils.CheckError(err)
		if !isUpdate {
			fmt.Printf("Current version, %s, is the latest. Nothing to update.\n", currentVersion)
		} else {
			fmt.Printf("Updating from version %s to version %s\n", currentVersion, newVersion)
			err = versioninfo.Update()
			if err != nil {
				var notSupportedError *updater.NotSupportedError
				if errors.As(err, &notSupportedError) {
					fmt.Fprintf(os.Stderr, "Self updating is not supported for %s. Please reinstall.", notSupportedError.Platform)
				} else {
					fmt.Fprintf(os.Stderr, "Error updating: %v\n", err)
				}
				os.Exit(1)
			}
			fmt.Printf("Updated to version %s\n", newVersion)
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
