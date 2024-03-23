package common

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type CompletedMsg struct{}

type ErrorMsg struct {
	Err error
}

func NewErrorMsg(err string) ErrorMsg {
	return ErrorMsg{
		Err: fmt.Errorf("%s", err),
	}
}

func (e ErrorMsg) Error() string {
	return e.Err.Error()
}

func NewErrorCmd(err string) tea.Cmd {
	return func() tea.Msg {
		return NewErrorMsg(err)
	}
}
