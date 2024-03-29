package cmd

import (
	"os"

	"github.com/dworthen/changelog/internal/cmd/addcmd"
	"github.com/dworthen/changelog/internal/cmd/applycmd"
	"github.com/dworthen/changelog/internal/cmd/initcmd"
	"github.com/dworthen/changelog/internal/cmd/updatecmd"
	"github.com/dworthen/changelog/internal/cmd/versioncmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Description",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initcmd.InitCmd)
	rootCmd.AddCommand(addcmd.AddCmd)
	rootCmd.AddCommand(applycmd.ApplyCmd)
	rootCmd.AddCommand(updatecmd.UpdateCmd)
	rootCmd.AddCommand(versioncmd.VersionCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goda.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
