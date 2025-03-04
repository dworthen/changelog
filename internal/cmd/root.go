package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dworthen/changelog/internal/cliflags"
	"github.com/dworthen/changelog/internal/cmd/addcmd"
	"github.com/dworthen/changelog/internal/cmd/applycmd"
	"github.com/dworthen/changelog/internal/cmd/checkcmd"
	"github.com/dworthen/changelog/internal/cmd/initcmd"
	"github.com/dworthen/changelog/internal/cmd/updatecmd"
	"github.com/dworthen/changelog/internal/cmd/versioncmd"
	"github.com/dworthen/changelog/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Git-based changelog manager for JavaScript, Python, and Go projects.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	utils.CheckError(err)
}

func init() {
	cobra.OnInitialize(initialize)

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&cliflags.CWD, "cwd", "", ".", "project directory for changelog.")
	rootCmd.PersistentFlags().BoolVarP(&cliflags.Verbose, "verbose", "", false, "enable verbose logging.")

	rootCmd.AddCommand(initcmd.InitCmd)
	rootCmd.AddCommand(addcmd.AddCmd)
	rootCmd.AddCommand(applycmd.ApplyCmd)
	rootCmd.AddCommand(checkcmd.CheckCmd)
	rootCmd.AddCommand(updatecmd.UpdateCmd)
	rootCmd.AddCommand(versioncmd.VersionCmd)
}

func initialize() {
	cwd := utils.GetCWD()
	configDir := utils.GetChangelogDirPath()
	configFilePath := utils.GetConfigFilePath()
	logsFilePath := utils.GetLogsFilePath()
	directory := filepath.Dir(logsFilePath)
	_, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(directory, 0755)
			utils.CheckError(err)
		} else {
			utils.CheckError(err)
		}
	}

	logFileWriter, err := os.OpenFile(logsFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	utils.CheckError(err)

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	if cliflags.Verbose {
		opts.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(logFileWriter, opts))
	slog.SetDefault(logger)

	slog.Info("Starting changelog", "cwd", cwd, "configDir", configDir, "configFilePath", configFilePath, "logsFilePath", logsFilePath)
}
