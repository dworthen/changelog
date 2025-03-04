package add

import (
	"fmt"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/davecgh/go-spew/spew"
	"github.com/dworthen/changelog/internal/models/formwrapmodel"
)

var changeType string

func newChangeTypeForm() (tea.Model, error) {
	changelogEntry := GetChangeLogEntry()

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			slog.Debug("Command Start: forms.newChangeTypeForm.OnComplete", "changeType", changeType)
			changelogEntry.SetChange(changeType)
			returnMsg := pipelineStepCompleteMsg{}
			slog.Debug("Command End: forms.newChangeTypeForm.OnComplete returning msg", "changelogEntry", changelogEntry, "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Change Type").
				Description("The type of change to make.").
				Options(huh.NewOptions[string]("patch", "minor", "major")...).
				Value(&changeType),
		),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		Build()
}

func newDescriptionForm() (tea.Model, error) {
	changelogEntry := GetChangeLogEntry()
	description := ""

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			slog.Debug("Command Start: forms.newDescriptionForm.OnComplete", "description", description)
			changelogEntry.SetDescription(description)
			returnMsg := pipelineStepCompleteMsg{}
			slog.Debug("Command End: forms.newDescriptionForm.OnComplete returning msg", "changelogEntry", changelogEntry, "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Description").
				DescriptionFunc(func() string {
					return fmt.Sprintf("Describe the %s changes.", changeType)
				}, &changeType).
				Validate(func(value string) error {
					if strings.TrimSpace(value) == "" {
						return fmt.Errorf("Required.")
					}
					return nil
				}).
				Value(&description),
		),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		Build()
}

var formModelConstructors = []func() (tea.Model, error){
	newChangeTypeForm,
	newDescriptionForm,
}
