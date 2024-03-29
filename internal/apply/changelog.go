package apply

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/gitmanage"
	"github.com/hashicorp/go-version"
	"github.com/mailgun/raymond/v2"
)

type ChangeDescription struct {
	Sha         string
	ShortSha    string
	Change      string
	Description string
	File        string
}

type Changelog struct {
	OldVersion   string
	Version      string
	Change       string
	MajorChanges []ChangeDescription
	MinorChanges []ChangeDescription
	PatchChanges []ChangeDescription
}

type FrontMatter struct {
	Change string `json:"change"`
}

var changeSeverity map[string]int = map[string]int{
	"patch": 0,
	"minor": 1,
	"major": 2,
}

var severityChange map[int]string = map[int]string{
	0: "patch",
	1: "minor",
	2: "major",
}

func NewChangelog() (*Changelog, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	files, err := doublestar.FilepathGlob(".changelog/*.md", doublestar.WithFilesOnly())
	if err != nil {
		return nil, err
	}

	changelog := Changelog{
		MajorChanges: []ChangeDescription{},
		MinorChanges: []ChangeDescription{},
		PatchChanges: []ChangeDescription{},
	}

	currentChange := 0

	for _, file := range files {
		changeDescription, err := loadFile(file)
		if err != nil {
			return nil, err
		}
		newChange := changeSeverity[changeDescription.Change]
		if newChange > currentChange {
			currentChange = newChange
		}
		switch changeDescription.Change {
		case "patch":
			changelog.PatchChanges = append(changelog.PatchChanges, changeDescription)
		case "minor":
			changelog.MinorChanges = append(changelog.MinorChanges, changeDescription)
		case "major":
			changelog.MajorChanges = append(changelog.MajorChanges, changeDescription)
		default:
			return nil, fmt.Errorf("Unrecognized changelog change %s in file %s", changeDescription.Change, file)
		}
	}

	changelog.Change = severityChange[currentChange]
	changelog.OldVersion = conf.Version
	newVersion, err := bumpVersion(changelog.OldVersion, changelog.Change)
	if err != nil {
		return nil, err
	}
	changelog.Version = newVersion

	return &changelog, nil
}

func loadFile(filePath string) (ChangeDescription, error) {
	changeDescription := ChangeDescription{
		File: filePath,
	}

	hash, err := gitmanage.GetFileCommit(filePath)
	if err != nil {
		return changeDescription, err
	}
	changeDescription.Sha = hash
	changeDescription.ShortSha = hash[0:7]

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return changeDescription, err
	}

	var matter FrontMatter
	desc, err := frontmatter.Parse(strings.NewReader(string(fileContents)), &matter)
	description := strings.TrimSpace(string(desc))
	description = strings.ReplaceAll(description, "\r\n", "\n")
	descriptionLines := strings.Split(description, "\n")
	if len(descriptionLines) > 1 {
		for i := 1; i < len(descriptionLines); i++ {
			descriptionLines[i] = fmt.Sprintf("  %s", descriptionLines[i])
		}
	}
	description = strings.Join(descriptionLines, "\n")

	if err != nil {
		return changeDescription, err
	}

	changeDescription.Description = description
	changeDescription.Change = matter.Change

	return changeDescription, nil

}

func bumpVersion(oldVersion string, change string) (string, error) {
	ver, err := version.NewSemver(oldVersion)
	if err != nil {
		return "", err
	}
	verSegments := ver.Segments()
	if len(verSegments) != 3 {
		return "", fmt.Errorf("Unexpected semantic version. Expected X.X.X but got %s", oldVersion)
	}
	majorVer := verSegments[0]
	minorVer := verSegments[1]
	patchVer := verSegments[2]
	switch change {
	case "patch":
		patchVer += 1
	case "minor":
		minorVer += 1
		patchVer = 0
	case "major":
		majorVer += 1
		minorVer = 0
		patchVer = 0
	default:
		return "", fmt.Errorf("Unrecognized change type. Expected patch|minor|majore but got %s", change)
	}
	return fmt.Sprintf("%d.%d.%d", majorVer, minorVer, patchVer), nil
}

func getChangelogTemplate() (string, error) {
	filePath, err := filepath.Abs(".changelog/changelogTemplate.hbs")
	if err != nil {
		return "", err
	}

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(fileContents)), nil
}

func getChangelogContents() (string, error) {
	filePath, err := filepath.Abs("CHANGELOG.md")
	if err != nil {
		return "", err
	}

	_, err = os.Stat(filePath)
	if err != nil {
		return "", nil
	}

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(fileContents)), nil
}

func (changelog *Changelog) Save() error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	changelogTemplate, err := getChangelogTemplate()
	if err != nil {
		return err
	}
	changelogContents, err := getChangelogContents()
	if err != nil {
		return err
	}

	newContents, err := raymond.Render(changelogTemplate, changelog)
	if err != nil {
		return err
	}

	newChangelog := fmt.Sprintf("%s\n\n%s", strings.TrimSpace(newContents), changelogContents)

	err = os.WriteFile("CHANGELOG.md", []byte(newChangelog), 0644)
	if err != nil {
		return err
	}

	err = changelog.DeleteFiles()
	if err != nil {
		return err
	}

	conf.Version = changelog.Version
	return conf.Save()
}

func (changeDescription *ChangeDescription) DeleteFile() error {
	return os.Remove(changeDescription.File)
}

func deleteFiles(changes []ChangeDescription) error {
	for _, cd := range changes {
		err := cd.DeleteFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func (changelog *Changelog) DeleteFiles() error {
	err := deleteFiles(changelog.MajorChanges)
	if err != nil {
		return err
	}
	err = deleteFiles(changelog.MinorChanges)
	if err != nil {
		return err
	}
	err = deleteFiles(changelog.PatchChanges)
	if err != nil {
		return err
	}
	return nil
}

func (changelog *Changelog) BumpFiles() error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	for _, bumpInfo := range conf.BumpFiles {
		err = bumpInfo.Bump(changelog.Version)
		if err != nil {
			return err
		}
	}
	return nil
}
