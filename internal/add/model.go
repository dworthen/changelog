package add

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/dworthen/changelog/internal/common"
)

type ChangeFormState int

const (
	StateInitial ChangeFormState = iota
	StateDescription
	StateCompleted
	StateError
)

type Model struct {
	State             ChangeFormState
	form              *huh.Form
	ChangeType        string
	ChangeDescription string
	width             int
	height            int
	err               error
}

func NewModel() Model {
	m := Model{
		State:             StateInitial,
		ChangeType:        "",
		ChangeDescription: "",
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("changeType").
				Title("Change Type").
				Options(huh.NewOptions[string]("Patch", "Minor", "Major")...),
		),
	)

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.form.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case common.CompletedMsg:
		m.State = StateCompleted
		return m, tea.Quit
	case common.ErrorMsg:
		m.State = StateError
		m.err = msg.Err
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form.WithWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.State {
			case StateCompleted, StateError:
				return m, tea.Quit
			}
		}
	}

	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		switch m.State {
		case StateDescription:
			m.ChangeDescription = m.form.GetString("description")
			cmds = append(cmds, saveFile(m.ChangeType, m.ChangeDescription))
			// cmds = append(cmds, tea.Quit)
		case StateInitial:
			m.ChangeType = m.form.GetString("changeType")
			m.State = StateDescription
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewText().
						Key("description").
						Title("Change Description").Validate(func(value string) error {
						if strings.TrimSpace(value) == "" {
							return fmt.Errorf("Required.")
						}
						return nil
					}),
				),
			)
			cmds = append(cmds, m.form.Init())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.State {
	case StateError:
		return m.err.Error()
	case StateCompleted:
		return fmt.Sprintf("%s %s", m.ChangeType, m.ChangeDescription)
	default:
		return m.form.View()
	}
}
