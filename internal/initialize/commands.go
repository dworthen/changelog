package initialize

import (
	_ "embed"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dworthen/changelog/internal/common"
	"github.com/dworthen/changelog/internal/config"
)

//go:embed changelogTemplate.hbs
var changeLogTemplate []byte

func initialize(config *config.Config) tea.Cmd {
	return func() tea.Msg {

		dir, err := filepath.Abs(".changelog")
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		err = os.MkdirAll(dir, 0744)
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		filename := filepath.Join(dir, "changelogTemplate.hbs")
		err = os.WriteFile(filename, changeLogTemplate, 0644)
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		err = config.Save()
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		return common.CompletedMsg{}
	}
}
