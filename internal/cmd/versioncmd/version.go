package versioncmd

import (
	"fmt"

	"github.com/dworthen/changelog/internal/utils"
	"github.com/dworthen/changelog/internal/versioninfo"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version of changelog",
	Long:  "Print the current version of changelog",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := versioninfo.GetVersion()
		utils.CheckError(err)
		fmt.Printf("Changelog Version: %s\n", version)
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
