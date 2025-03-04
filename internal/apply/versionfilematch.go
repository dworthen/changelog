package apply

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/utils"
)

type versionFileMatch struct {
	Match      string
	StartIndex int
	EndIndex   int
	LineNumber int
	LineIndex  int
}

func newVersionFileMatches(vf config.VersionFile, cwd string) ([]versionFileMatch, error) {
	versionPattern, err := regexp.Compile(vf.Pattern)
	if err != nil {
		return nil, fmt.Errorf("Invalid version pattern: %s", vf.Pattern)
	}

	fullpath := utils.JoinPaths(cwd, vf.Path)

	_, err = os.Stat(fullpath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Version file, %s, does not exist within working directory, %s", vf.Path, cwd)
		}
		return nil, err
	}

	fileBytes, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	fileContents := string(fileBytes)
	matches := versionPattern.FindAllStringSubmatchIndex(fileContents, -1)
	if matches == nil || len(matches) == 0 {
		return nil, fmt.Errorf("Version pattern, %s, not found in file, %s", vf.Pattern, vf.Path)
	}

	versionFileMatches := []versionFileMatch{}

	for _, match := range matches {
		versionMatch := fileContents[match[2]:match[3]]
		lineNumber := strings.Count(fileContents[:match[0]], "\n") + 1
		lastNewlineIndex := strings.LastIndex(fileContents[:match[0]], "\n")
		lineIndex := match[2] - lastNewlineIndex
		versionFileMatches = append(versionFileMatches, versionFileMatch{
			Match:      versionMatch,
			StartIndex: match[2],
			EndIndex:   match[3],
			LineNumber: lineNumber,
			LineIndex:  lineIndex,
		})
	}

	return versionFileMatches, nil

}
