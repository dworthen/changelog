package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dworthen/changelog/internal/cliflags"
)

func ToPath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load user home dir.")
		os.Exit(1)
	}
	return filepath.Clean(strings.ReplaceAll(path, "~", homeDir))
}

func JoinPaths(paths ...string) string {
	resolvedPaths := []string{}
	for _, path := range paths {
		resolvedPaths = append(resolvedPaths, ToPath(path))
	}
	return filepath.Clean(filepath.Join(resolvedPaths...))
}

func ToFullPath(path string) string {
	fullPath, err := filepath.Abs(ToPath(path))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get absolute path.")
		os.Exit(1)
	}
	return fullPath
}

func GetCWD() string {
	if cliflags.CWD != "" {
		return ToFullPath(cliflags.CWD)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	return ToFullPath(cwd)
}

func GetChangelogDirPath() string {
	return JoinPaths(GetCWD(), ".changelog")
}

func GetLogsFilePath() string {
	return JoinPaths(GetChangelogDirPath(), "logs.json")
}

func GetConfigFilePath() string {
	return JoinPaths(GetChangelogDirPath(), "config.yaml")
}

func GetChangelogTemplatePath() string {
	return JoinPaths(GetChangelogDirPath(), "changelogTemplate.hbs")
}

func GetGitIgnorePath() string { return JoinPaths(GetChangelogDirPath(), ".gitignore") }
