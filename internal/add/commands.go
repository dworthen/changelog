package add

import (
	tea "github.com/charmbracelet/bubbletea"
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
		return returnMsg
	}
}
