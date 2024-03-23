package initialize

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/dworthen/changelog/internal/config"
	"github.com/hashicorp/go-version"
)

var versionForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("Current Version").
			Validate(func(value string) error {
				if strings.TrimSpace(value) == "" {
					return fmt.Errorf("Required.")
				}

				_, err := version.NewSemver(value)
				if err != nil {
					return fmt.Errorf("Not a valid semantic version. Must adhere to https://semver.org/")
				}

				conf, err := config.GetConfig()
				if err != nil {
					return err
				}
				conf.Version = value
				return nil
			}),
	),
)

var onAddForm = huh.NewForm(
	huh.NewGroup(
		huh.NewConfirm().
			Title("Commit files when adding new changelog entries? The changelog description will be used as the commit message.").
			Validate(func(value bool) error {
				conf, err := config.GetConfig()
				if err != nil {
					return err
				}
				conf.OnAdd.ComitFiles = value
				return nil
			}),
	),
)

var onApplyCommitForm = huh.NewForm(
	huh.NewGroup(
		huh.NewConfirm().
			Title("Commit files when applying changelog changes?").
			Validate(func(value bool) error {
				conf, err := config.GetConfig()
				if err != nil {
					return err
				}
				conf.OnApply.CommitFiles = value
				return nil
			}),
	),
)

var onApplyTagCommitForm = huh.NewForm(
	huh.NewGroup(
		huh.NewConfirm().
			Title("Tag Commit?").
			Validate(func(value bool) error {
				conf, err := config.GetConfig()
				if err != nil {
					return err
				}
				conf.OnApply.TagCommit = value
				return nil
			}),
	),
)

var onApplyTagFormatForm = huh.NewForm(
	huh.NewGroup(
		huh.NewInput().
			Title("Tag Format. (default v{{version}})").
			Value(&defaultTagFormat).
			Validate(func(value string) error {
				if strings.TrimSpace(value) == "" {
					return fmt.Errorf("Required.")
				}
				conf, err := config.GetConfig()
				if err != nil {
					return err
				}
				conf.OnApply.TagFormat = value
				return nil
			}),
	),
)

var bumpFilesForm = huh.NewForm(
	huh.NewGroup(
		huh.NewText().
			Title("Bump Files. Specify files to update version info when apply changelog changes.").
			Description("Enter a bump file per line in the format FILE_PATH=JSON_PATH where the file path is relative to cwd and json path describes the location using dot notation of semantic version value to bump. Example: package.json=version").
			Validate(func(value string) error {
				trimmedValue := strings.TrimSpace(value)
				if trimmedValue == "" {
					return nil
				}
				normalizedInput := strings.ReplaceAll(trimmedValue, "\r\n", "\n")
				var lines []string

				for _, line := range strings.Split(normalizedInput, "\n") {
					trimmedLine := strings.TrimSpace(line)
					if trimmedLine != "" {
						lines = append(lines, trimmedLine)
					}
				}

				bumpFiles := []config.BumpInfo{}

				for _, line := range lines {
					parts := strings.Split(line, "=")
					if len(parts) != 2 {
						return fmt.Errorf("Invalid format. Each row should be a single entry in form FILE_PATH=JSON_PATH_TO_FIELD_TO_BUMP")
					}
					bumpFiles = append(bumpFiles, config.BumpInfo{
						FileName: parts[0],
						JsonPath: parts[1],
					})
				}

				conf, err := config.GetConfig()
				if err != nil {
					return err
				}
				conf.BumpFiles = bumpFiles

				return nil
			}),
	),
)
