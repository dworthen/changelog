package apply

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/davecgh/go-spew/spew"
	"github.com/dworthen/changelog/internal/models/formwrapmodel"
)

func newApproveChangelogForm() (tea.Model, error) {

	approved := false

	onComplete := func() tea.Cmd {
		return func() tea.Msg {
			returnMsg := ApplyModelConfirmationCompleteMsg{
				approved: approved,
			}
			slog.Debug("Command Start - End: newApproveChangelogForm.OnComplete returning msg", "msg", spew.Sdump(returnMsg))
			return returnMsg
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Approve Changelog").
				Description("Continue with the above changes?").
				Value(&approved).
				Affirmative("Yes").
				Negative("No"),
		),
	)

	return formwrapmodel.NewFormWrapModelBuilder().
		WithForm(form).
		WithOnComplete(onComplete).
		WithShowHelp(false).
		Build()

}
