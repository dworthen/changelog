package initialize

import (
	_ "embed"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/models/common"
	"github.com/dworthen/changelog/internal/utils"
)

//go:embed changelogTemplate.hbs
var changeLogTemplate []byte

//go:embed gitignoreTemplate.hbs
var gitignoreTemplate []byte

func onInitializeCompleteCmd() tea.Cmd {
	return func() tea.Msg {
		configData, err := config.LoadConfig()
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}
		err = configData.Save()
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		changelogTemplatePath := utils.GetChangelogTemplatePath()

		_, err = os.Stat(changelogTemplatePath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.WriteFile(changelogTemplatePath, changeLogTemplate, 0644)
				if err != nil {
					return common.ErrorMsg{
						Err: fmt.Errorf("Error creating changelog template file: %w", err),
					}
				}
			} else {
				return common.ErrorMsg{
					Err: fmt.Errorf("Error checking for changelog template file: %w", err),
				}
			}
		}

		gitignoreTemplatePath := utils.GetGitIgnorePath()
		err = os.WriteFile(gitignoreTemplatePath, gitignoreTemplate, 0644)
		if err != nil {
			return common.ErrorMsg{
				Err: fmt.Errorf("Error creating .gitignore file: %w", err),
			}
		}

		returnMsg := pipelineCompleteMsg{}
		return returnMsg
	}
}
