package apply

import (
	"fmt"
	"os"
	"strings"

	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/utils"
)

type versionFileMatchesContainer struct {
	VersionFiles config.VersionFile
	Matches      []versionFileMatch
}

func (c *versionFileMatchesContainer) bump(newVersion string) error {
	cwd := utils.GetCWD()

	fullPath := utils.JoinPaths(cwd, c.VersionFiles.Path)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Error bumping version. Version file, %s, does not exist within working directory, %s. Error: %w", c.VersionFiles.Path, cwd, err)
		}
		return err
	}

	fileBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("Error bumping version. Failed to read %s. Error: %w", fullPath, err)
	}

	fileContents := string(fileBytes)
	fileContentSegments := []string{}
	startIndex := 0
	endIndex := 0

	for _, match := range c.Matches {
		endIndex = match.StartIndex
		fileContentSegments = append(fileContentSegments, fileContents[startIndex:endIndex])
		fileContentSegments = append(fileContentSegments, newVersion)
		startIndex = match.EndIndex
	}
	fileContentSegments = append(fileContentSegments, fileContents[startIndex:])

	updatedFileContents := strings.Join(fileContentSegments, "")
	err = os.WriteFile(fullPath, []byte(updatedFileContents), 0644)
	if err != nil {
		return fmt.Errorf("Error bumping version. Failed to write to %s. Error: %w", fullPath, err)
	}
	return nil
}
