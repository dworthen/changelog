package initialize

import (
	"fmt"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/davecgh/go-spew/spew"
	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/models/formwrapmodel"
	"github.com/dworthen/changelog/internal/models/pipelinemodel"
	"github.com/hashicorp/go-version"
)

func newChangelogEntryForm() (tea.Model, error) {
	configData, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	changelogPath := configData.GetChangelogFile()
	currentVersion := configData.GetVersion()

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			slog.Debug("Command Start: forms.newChangelogEntryForm.OnComplete", "changelogPath", changelogPath)
			configData.SetChangelogFile(changelogPath)
			configData.SetVersion(currentVersion)
			returnMsg := pipelineStepCompleteMsg{}
			slog.Debug("Command End: forms.newChangelogEntryForm.OnComplete returning msg", "config", configData, "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Changelog File Path").
				Description("The relative path to the changelog file to manage, e.g., CHANGELOG.md").
				Validate(func(value string) error {
					if strings.TrimSpace(value) == "" {
						return fmt.Errorf("Required.")
					}
					return nil
				}).Value(&changelogPath),
			huh.NewInput().
				Title("Current Version").
				Description("The current semantic version of the project, e.g., 0.0.0").
				Validate(func(value string) error {
					if strings.TrimSpace(value) == "" {
						return fmt.Errorf("Required.")
					}
					_, err := version.NewSemver(value)
					if err != nil {
						return fmt.Errorf("Not a valid semantic version. Must adhere to https://semver.org/")
					}
					return nil
				}).Value(&currentVersion),
		),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		Build()

}

var changelogFilesDescription = strings.TrimSpace(`
Changelog can manage version files, bumping the version in each file when running 'changelog apply'.
Versions are updated based on a regular expression defined per file.
Each regular expression must have a single capture group that captures the version to update/replace.

Example:
	- path: package.json
	  pattern: "version":\\s\*"(\\d+\\.\\d+\\.\\d+)"
	- path: pyproject.toml
	  pattern: version\\s*=\\s*"(\\d+\\.\\d+\\.\\d+)"
`)

var providedVersionFiles = false

func newChangelogFilesForm() (tea.Model, error) {
	configData, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	filePath := ""
	filePattern := ""

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			slog.Debug("Command Start: forms.newChangelogFilesForm.OnComplete", "filePath", filePath, "filePattern", filePattern)
			if filePath != "" && filePattern != "" {
				configData.AddVersionFile(filePath, filePattern)
			}
			returnMsg := pipelineStepCompleteMsg{}
			slog.Debug("Command End: forms.newChangelogFilesForm.OnComplete returning msg", "config", configData, "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Version Files").Description(changelogFilesDescription),
			huh.NewInput().
				Title("File Path").
				Description("The relative path to the file to manage, e.g., package.json.").
				Value(&filePath),
			huh.NewInput().
				Title("Regular Expression Pattern").
				Description("The pattern to use to extract the version from the file.").
				Validate(func(value string) error {
					if filePath != "" && strings.TrimSpace(value) == "" {
						return fmt.Errorf("Required.")
					}

					if filePath != "" && value != "" {
						providedVersionFiles = true
					} else {
						providedVersionFiles = false
					}

					return nil
				}).
				Value(&filePattern),
		),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		Build()
}

func newAddNewFilesPromptForm() (tea.Model, error) {
	addAnotherFile := false

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			slog.Debug("Command Start: forms.newAddNewFilesPromptForm.OnComplete", "addAnotherFile", addAnotherFile)
			var returnMsg tea.Msg = pipelineStepCompleteMsg{}
			if addAnotherFile {
				returnMsg = pipelinemodel.PipelineModelChangeIndexMsg{
					Amount: -1,
				}
			}
			slog.Debug("Command End: forms.newAddNewFilesPromptForm.OnComplete returning msg", "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add Another File").
				Description("Would you like to add another version file to manage?").
				Affirmative("Yes").
				Negative("No").
				Value(&addAnotherFile),
		).WithHideFunc(func() bool {
			return !providedVersionFiles
		}),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		Build()
}

var onAddDescription = strings.TrimSpace(`
The command 'changelog add' creates a new changelog description file within the .changelog directory.
Changelog add can be configured to add the new description file to staged git files and then commit
the staged files using the description as the git commit message. In short, changelog add can replace
git commit.

Example workflow
	1. Make changes
	2. git add files...
	3. changelog add
`)

func newOnAddForm() (tea.Model, error) {
	configData, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	commitFiles := configData.GetOnAddCommitFiles()

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			slog.Debug("Command Start: forms.newOnAddForm.OnComplete", "commitFiles", commitFiles)
			configData.SetOnAddCommitFiles(commitFiles)
			returnMsg := pipelineStepCompleteMsg{}
			slog.Debug("Command End: forms.newOnAddForm.OnComplete returning msg", "config", configData, "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Changelog Add Settings").Description(onAddDescription),
			huh.NewConfirm().
				Title("On 'changelog add', Commit Files?").
				Description("Would you like to commit staged files when adding a new changelog entry? The changelog description is used as the commit message.").
				Affirmative("Yes").
				Negative("No").
				Value(&commitFiles),
		),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		Build()
}

var onApplyDescription = strings.TrimSpace(`
Like 'changelog add', 'changelog apply' can be configured to stage and commit impacted files,
those that have received a version bump along with the changelog file (e.g., CHANGELOG.md).
The command can also be configured to tag the commit with the new version.

Changelog apply will also run any commands provided in the 'Commands' field after bumping the version
and before committing the files. Example commands include:

	- npm install
	- go mod tidy
	- uv sync
`)

func newOnApplyForm() (tea.Model, error) {
	configData, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	commitFiles := configData.GetOnApplyCommitFiles()
	tagCommit := configData.GetOnApplyTagCommit()
	tagFormat := configData.GetOnApplyTagFormat()
	commands := ""

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			slog.Debug("Command Start: forms.newOnApplyForm.OnComplete", "commitFiles", commitFiles, "tagCommit", tagCommit, "tagFormat", tagFormat, "commands", commands)
			if strings.TrimSpace(commands) != "" {
				commandsList := strings.Split(commands, "\n")
				configData.SetOnApplyCommands(commandsList)
			}
			configData.SetOnApplyCommitFiles(commitFiles)
			if !commitFiles {
				configData.SetOnApplyTagCommit(false)
				configData.SetOnApplyTagFormat("")
			} else {
				configData.SetOnApplyTagCommit(tagCommit)
				configData.SetOnApplyTagFormat(tagFormat)
			}
			returnMsg := pipelineStepCompleteMsg{}
			slog.Debug("Command End: forms.newOnApplyForm.OnComplete returning msg", "config", configData, "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	noteField := huh.NewNote().Title("Changelog Apply Settings").Description(onApplyDescription)

	form := huh.NewForm(
		huh.NewGroup(
			noteField,
			huh.NewText().
				Title("Commands").
				Description("The commands to run when applying the changelog changes. One command per line.").
				Value(&commands),
		),
		huh.NewGroup(
			noteField,
			huh.NewConfirm().
				Title("On 'changelog apply', Commit Files?").
				Description("Would you like to commit the changelog file(s) when applying the changelog changes?").
				Affirmative("Yes").
				Negative("No").
				Value(&commitFiles),
		),
		huh.NewGroup(
			noteField,
			huh.NewConfirm().
				Title("On 'changelog apply', Tag Commit").
				Description("Would you like to tag the commit when applying the changelog changes?").
				Affirmative("Yes").
				Negative("No").
				Value(&tagCommit),
		).WithHideFunc(func() bool {
			return !commitFiles
		}),
		huh.NewGroup(
			noteField,
			huh.NewInput().
				Title("'changelog apply' Tag Format").
				Description("The format to use when tagging the commit. The version will be available as {{version}}, e.g., v{{version}}").
				Validate(func(value string) error {
					if tagCommit && strings.TrimSpace(value) == "" {
						return fmt.Errorf("Required.")
					}
					return nil
				}).
				Value(&tagFormat),
		).WithHideFunc(func() bool {
			return !commitFiles || !tagCommit
		}),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		Build()
}

var formModelConstructors = []func() (tea.Model, error){
	newChangelogEntryForm,
	newChangelogFilesForm,
	newAddNewFilesPromptForm,
	newOnAddForm,
	newOnApplyForm,
}
