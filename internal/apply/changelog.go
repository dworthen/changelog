package apply

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/gitmanage"
	"github.com/dworthen/changelog/internal/globals"
	"github.com/dworthen/changelog/internal/iowriters"
	"github.com/dworthen/changelog/internal/utils"
	"github.com/hashicorp/go-version"
	"github.com/mailgun/raymond/v2"
)

type changelog struct {
	WorkingDirectory   string
	ChangelogDirectory string
	ConfigFile         string
	ChangelogFile      string
	ChangelogFiles     []string
	Commands           []string
	MajorChanges       []changeDescription
	MinorChanges       []changeDescription
	PatchChanges       []changeDescription
	NumMajorChanges    int
	NumMinorChanges    int
	NumPatchChanges    int
	BumpingVersion     bool
	VersionBumpType    string
	OldVersion         string
	Version            string
	VersionFiles       []config.VersionFile
	VersionFileMatches map[string]versionFileMatchesContainer
	ChangelogEntry     string
	Commit             bool
	Tag                bool
	TagFormat          string
	TagParsed          string
}

func newChangelog() (*changelog, error) {
	c := &changelog{
		WorkingDirectory:   utils.GetCWD(),
		ChangelogDirectory: utils.GetChangelogDirPath(),
		ConfigFile:         "",
		ChangelogFiles:     []string{},
		Commands:           []string{},
		MajorChanges:       []changeDescription{},
		MinorChanges:       []changeDescription{},
		PatchChanges:       []changeDescription{},
		NumMajorChanges:    0,
		NumMinorChanges:    0,
		NumPatchChanges:    0,
		BumpingVersion:     false,
		VersionBumpType:    "None",
		OldVersion:         "",
		Version:            "",
		VersionFiles:       []config.VersionFile{},
		VersionFileMatches: map[string]versionFileMatchesContainer{},
		ChangelogEntry:     "",
		Commit:             false,
		Tag:                false,
		TagFormat:          "",
		TagParsed:          "",
	}

	return c.populateChangelog(withConfigData(), withVersionFileMatches(), withChangelogFiles(), withChangeDescriptions(), withBumpInfo(), withChangelogEntry())
}

type changelogLoaders = func(*changelog) error

func (c *changelog) populateChangelog(loaders ...changelogLoaders) (*changelog, error) {
	for _, loader := range loaders {
		if err := loader(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func withConfigData() changelogLoaders {
	return func(c *changelog) error {
		configData, err := config.LoadConfig()
		if err != nil {
			return err
		}

		cwd := utils.GetCWD()
		configFile := utils.GetConfigFilePath()
		relativeConfigFile, err := filepath.Rel(cwd, configFile)
		if err != nil {
			return err
		}
		c.ConfigFile = relativeConfigFile
		c.ChangelogFile = configData.ChangelogFile
		c.OldVersion = configData.Version
		c.VersionFiles = configData.Files
		c.Commands = configData.OnApply.Commands
		c.Commit = configData.OnApply.CommitFiles
		c.Tag = configData.OnApply.TagCommit
		c.TagFormat = configData.OnApply.TagFormat
		return nil
	}
}

func withVersionFileMatches() changelogLoaders {
	return func(c *changelog) error {
		cwd := utils.GetCWD()
		for _, vf := range c.VersionFiles {
			matches, err := newVersionFileMatches(vf, cwd)
			if err != nil {
				return err
			}
			if entry, ok := c.VersionFileMatches[vf.Path]; ok {
				entry.Matches = append(c.VersionFileMatches[vf.Path].Matches, matches...)
				c.VersionFileMatches[vf.Path] = entry
			} else {
				c.VersionFileMatches[vf.Path] = versionFileMatchesContainer{vf, matches}
			}
		}
		return nil
	}
}

func withChangelogFiles() changelogLoaders {
	return func(c *changelog) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		changelogDir := utils.GetChangelogDirPath()
		relDir, err := filepath.Rel(cwd, changelogDir)
		if err != nil {
			return err
		}
		changelogEntryPattern := filepath.Join(relDir, "*.md")
		changelogFiles, err := filepath.Glob(changelogEntryPattern)
		if err != nil {
			return err
		}

		resolvedChangeLogFiles := []string{}
		for _, file := range changelogFiles {
			fullPath := filepath.Join(cwd, file)
			relPath, err := filepath.Rel(utils.GetCWD(), fullPath)
			if err != nil {
				return err
			}
			resolvedChangeLogFiles = append(resolvedChangeLogFiles, relPath)
		}

		c.ChangelogFiles = resolvedChangeLogFiles

		return nil
	}
}

func withChangeDescriptions() changelogLoaders {
	return func(c *changelog) error {
		changeDescriptions, err := loadChangeDescriptions(c.ChangelogFiles...)
		if err != nil {
			return err
		}

		for _, cd := range changeDescriptions {
			c.BumpingVersion = true
			switch cd.Change {
			case "major":
				c.MajorChanges = append(c.MajorChanges, cd)
			case "minor":
				c.MinorChanges = append(c.MinorChanges, cd)
			case "patch":
				c.PatchChanges = append(c.PatchChanges, cd)
			default:
				return fmt.Errorf("invalid change type: %s. Supported types: 'patch', 'minor', 'major'", cd.Change)
			}
		}

		c.NumMajorChanges = len(c.MajorChanges)
		c.NumMinorChanges = len(c.MinorChanges)
		c.NumPatchChanges = len(c.PatchChanges)

		return nil
	}
}

func withBumpInfo() changelogLoaders {
	return func(c *changelog) error {
		if c.NumMajorChanges > 0 {
			c.VersionBumpType = "major"
		} else if c.NumMinorChanges > 0 {
			c.VersionBumpType = "minor"
		} else if c.NumPatchChanges > 0 {
			c.VersionBumpType = "patch"
		}

		if c.VersionBumpType == "None" {
			return nil
		}

		ver, err := version.NewSemver(c.OldVersion)
		if err != nil {
			return err
		}
		verSegments := ver.Segments()
		if len(verSegments) != 3 {
			return fmt.Errorf("Unexpected semantic version. Expected X.X.X but got %s", c.OldVersion)
		}

		majorVer := verSegments[0]
		minorVer := verSegments[1]
		patchVer := verSegments[2]

		switch c.VersionBumpType {
		case "major":
			majorVer++
			minorVer = 0
			patchVer = 0
		case "minor":
			minorVer++
			patchVer = 0
		case "patch":
			patchVer++
		}

		c.Version = fmt.Sprintf("%d.%d.%d", majorVer, minorVer, patchVer)
		tag, err := raymond.Render(c.TagFormat, map[string]string{"version": c.Version})
		if err != nil {
			return err
		}
		c.TagParsed = tag

		return nil
	}
}

func withChangelogEntry() changelogLoaders {
	return func(c *changelog) error {
		changelogTemplatePath := utils.GetChangelogTemplatePath()
		changelogTemplateContents, err := os.ReadFile(changelogTemplatePath)
		if err != nil {
			return err
		}
		changelogEntry, err := raymond.Render(string(changelogTemplateContents), c)
		if err != nil {
			return err
		}
		c.ChangelogEntry = changelogEntry
		return nil
	}
}

//go:embed summaryTemplate.hbs
var summaryTemplate string

func (c *changelog) getSummary() (string, error) {
	return raymond.Render(summaryTemplate, c)
}

func (c *changelog) updateChangelog() error {
	changelogPath := utils.JoinPaths(c.WorkingDirectory, c.ChangelogFile)
	changelogContents := ""

	_, err := os.Stat(changelogPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		changelogBytes, err := os.ReadFile(changelogPath)
		if err != nil {
			return err
		}
		changelogContents = string(changelogBytes)
	}

	changelogContents = c.ChangelogEntry + changelogContents
	return os.WriteFile(changelogPath, []byte(changelogContents), 0644)
}

func (c *changelog) apply(writer *iowriters.StdioWriter) error {
	configData, err := config.LoadConfig()
	if err != nil {
		return err
	}

	filesToCommit := []string{"."}

	err = c.updateChangelog()
	if err != nil {
		return err
	}
	// filesToCommit = append(filesToCommit, c.ChangelogFile)

	for _, vfm := range c.VersionFileMatches {
		if err := vfm.bump(c.Version); err != nil {
			return err
		}
		// filesToCommit = append(filesToCommit, vfm.VersionFiles.Path)
	}

	for _, changelogEntryFile := range c.ChangelogFiles {
		err := os.Remove(utils.JoinPaths(c.WorkingDirectory, changelogEntryFile))
		if err != nil {
			return err
		}

		// filesToCommit = append(filesToCommit, changelogEntryFile)
	}

	configData.SetVersion(c.Version)
	err = configData.Save()
	if err != nil {
		return err
	}
	// configPath := utils.GetConfigFilePath()
	cwd := utils.GetCWD()
	// relativeConfigPath, err := filepath.Rel(cwd, configPath)
	// if err != nil {
	// 	return err
	// }
	// filesToCommit = append(filesToCommit, relativeConfigPath)

	if len(c.Commands) > 0 {
		globals.Program.Send(ApplyModelSetStateMsg{
			State: ApplyModelStateRunningCommands,
		})
	}

	for _, command := range c.Commands {
		cmdToRun := theme.Focused.Title.Render(fmt.Sprintf("%s\n", command))
		writer.Write([]byte(cmdToRun))
		commandList := strings.Split(command, " ")
		cmd := exec.Command(commandList[0], commandList[1:]...)
		cmd.Dir = cwd
		cmd.Stdout = writer
		cmd.Stderr = writer
		err = cmd.Run()
		if err != nil {
			errMSg := theme.Focused.ErrorIndicator.Render(err.Error() + "\n")
			writer.Write([]byte(errMSg))
			break
		}
	}

	if c.Commit {
		err := gitmanage.CommitFiles(filesToCommit, fmt.Sprintf("bump version to %s", c.Version))
		if err != nil {
			return err
		}
	}

	if c.Tag {
		err := gitmanage.Tag(c.TagParsed)
		if err != nil {
			return err
		}
	}

	return nil
}
