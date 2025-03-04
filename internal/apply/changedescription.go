package apply

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/dworthen/changelog/internal/gitmanage"
	"github.com/dworthen/changelog/internal/utils"
)

type changeDescription struct {
	Sha         string
	ShortSha    string
	Change      string
	Description string
	FilePath    string
}

type frontMatter struct {
	Change string `yaml:"change"`
}

func loadChangeDescriptions(files ...string) ([]changeDescription, error) {
	changeDescriptions := []changeDescription{}

	for _, file := range files {
		cd := changeDescription{}

		sha, err := gitmanage.GetCommitHashForFile(file)
		if err != nil {
			return nil, err
		}
		cd.Sha = sha
		cd.ShortSha = sha[:7]

		fullPath := utils.JoinPaths(utils.GetCWD(), file)
		cd.FilePath = fullPath

		fileContents, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("Error loading changelog entries. Failed to read file %s. Error: %w", fullPath, err)
		}

		var fm frontMatter
		desc, err := frontmatter.Parse(strings.NewReader(string(fileContents)), &fm)
		if err != nil {
			return nil, fmt.Errorf("Error loading changelog entries. Failed to parse frontmatter in file %s. Error: %w", fullPath, err)
		}
		description := strings.TrimSpace(string(desc))
		description = strings.ReplaceAll(description, "\r\n", "\n")
		descriptionLines := strings.Split(description, "\n")
		newDescriptionLines := []string{descriptionLines[0]}
		if len(descriptionLines) > 1 {
			for i := 1; i < len(descriptionLines); i++ {
				if strings.TrimSpace(descriptionLines[i]) == "" {
					continue
				}
				newDescriptionLines = append(newDescriptionLines, fmt.Sprintf("  %s", descriptionLines[i]))
			}
		}
		description = strings.Join(newDescriptionLines, "\n")

		cd.Change = strings.ToLower(fm.Change)
		cd.Description = description

		changeDescriptions = append(changeDescriptions, cd)
	}

	return changeDescriptions, nil
}
