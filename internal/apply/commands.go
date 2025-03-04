package apply

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dworthen/changelog/internal/globals"
	"github.com/dworthen/changelog/internal/iowriters"
	"github.com/dworthen/changelog/internal/models/common"
)

func loadChangelogCmd() tea.Cmd {
	return func() tea.Msg {
		changelogData, err := newChangelog()
		if err != nil {
			return common.ErrorMsg{Err: err}
		}
		returnMsg := ApplyModelChangelogLoaddedMsg{
			changelog: changelogData,
		}
		return returnMsg
	}
}

func applyModelCompleteCmd(c *changelog) tea.Cmd {
	return func() tea.Msg {
		var stdioWriterFunc iowriters.StdioWriterFunc = func(output string) {
			globals.Program.Send(ApplyModelAppendCommandsOutputMsg{
				Output: output,
			})
		}
		stdioWriter := iowriters.NewStdioWriter(stdioWriterFunc)

		err := c.apply(stdioWriter)
		if err != nil {
			return common.ErrorMsg{Err: err}
		}
		returnMsg := ApplyModelSetStateMsg{
			State: ApplyModelStateComplete,
		}
		return returnMsg
	}
}
