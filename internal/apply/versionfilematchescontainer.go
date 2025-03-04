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

	fullpath := utils.JoinPaths(cwd, c.VersionFiles.Path)
	_, err := os.Stat(fullpath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Version file, %s, does not exist within working directory, %s", c.VersionFiles.Path, cwd)
		}
		return err
	}

	fileBytes, err := os.ReadFile(fullpath)
	if err != nil {
		return err
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
	return os.WriteFile(fullpath, []byte(updatedFileContents), 0644)
}
