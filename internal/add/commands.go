package add

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "embed"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dworthen/changelog/internal/common"
	"github.com/dworthen/changelog/internal/config"
	"github.com/dworthen/changelog/internal/gitmanage"
	"github.com/mailgun/raymond/v2"
)

//go:embed fileTemplate.hbs
var fileTemplate string

func saveFile(change string, description string) tea.Cmd {
	return func() tea.Msg {

		filename := fmt.Sprintf("%d.md", time.Now().UTC().Unix())
		relFilePath := filepath.Join(".changelog", filename)
		filePath, err := filepath.Abs(relFilePath)
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}
		dir := filepath.Dir(filePath)

		stats, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				return common.NewErrorMsg(".changelog directory does not exist. Run changelog init.")
			} else {
				return common.ErrorMsg{
					Err: err,
				}
			}
		}
		if !stats.IsDir() {
			return common.NewErrorMsg(".changelog is not a directory")
		}

		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		contents, err := raymond.Render(fileTemplate, map[string]string{
			"change":      strings.ToLower(change),
			"description": description,
		})
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		os.WriteFile(filePath, []byte(contents), 0644)

		conf, err := config.GetConfig()
		if err != nil {
			return common.ErrorMsg{
				Err: err,
			}
		}

		if conf.OnAdd.ComitFiles {
			err = gitmanage.CommitFiles([]string{relFilePath}, description)
			if err != nil {
				return common.ErrorMsg{
					Err: err,
				}
			}
		}

		return common.CompletedMsg{}
	}
}
