package add

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"github.com/dworthen/changelog/internal/models/common"
)

func onAddCompleteCmd() tea.Cmd {
	return func() tea.Msg {
		changelogEntry := GetChangeLogEntry()
		err := changelogEntry.Save()
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		returnMsg := pipelineCompleteMsg{}
		slog.Debug("Command End: add.NewAddModel.OnComplete - Saved changelog entry. Returning msg", "changelogEntry", changelogEntry, "msg", spew.Sdump(returnMsg))
		return returnMsg
	}
}
