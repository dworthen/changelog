package initialize

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/dworthen/changelog/internal/common"
	"github.com/dworthen/changelog/internal/config"
)

type InitState int

const (
	StateInitial InitState = iota
	StateOnAdd
	StateOnApplyCommit
	StateOnApplyTagCommit
	StateOnApplyTagFormat
	StateBumpFiles
	StateCompleted
	StateError
)

type Model struct {
	State                InitState
	versionForm          *huh.Form
	onAddForm            *huh.Form
	onApplyCommitForm    *huh.Form
	onApplyTagCommitForm *huh.Form
	onApplyTagFormatForm *huh.Form
	bumpFilesForm        *huh.Form
	width                int
	height               int
	err                  error
}

var defaultTagFormat string = "v{{version}}"

func NewModel() Model {
	m := Model{
		State: StateInitial,
	}

	m.versionForm = versionForm
	m.onAddForm = onAddForm
	m.onApplyCommitForm = onApplyCommitForm
	m.onApplyTagCommitForm = onApplyTagCommitForm
	m.onApplyTagFormatForm = onApplyTagFormatForm
	m.bumpFilesForm = bumpFilesForm

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.versionForm.Init(),
		m.onAddForm.Init(),
		m.onApplyCommitForm.Init(),
		m.onApplyTagCommitForm.Init(),
		m.onApplyTagFormatForm.Init(),
		m.bumpFilesForm.Init(),
	)
}

func (m *Model) processMsg(msg tea.Msg) tea.Cmd {
	conf, err := config.GetConfig()
	if err != nil {
		return common.NewErrorCmd(err.Error())
	}

	var cmds []tea.Cmd
	switch m.State {
	case StateInitial:
		form, cmd := m.versionForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.versionForm = f
			cmds = append(cmds, cmd)
			if f.State == huh.StateCompleted {
				m.State = StateOnAdd
			}
		}
	case StateOnAdd:
		form, cmd := m.onAddForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.onAddForm = f
			cmds = append(cmds, cmd)
			if f.State == huh.StateCompleted {
				m.State = StateOnApplyCommit
			}
		}
	case StateOnApplyCommit:
		form, cmd := m.onApplyCommitForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.onApplyCommitForm = f
			cmds = append(cmds, cmd)
			if f.State == huh.StateCompleted {
				if conf.OnApply.CommitFiles {
					m.State = StateOnApplyTagCommit
				} else {
					conf.OnApply.TagCommit = false
					m.State = StateBumpFiles
				}
			}
		}
	case StateOnApplyTagCommit:
		form, cmd := m.onApplyTagCommitForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.onApplyTagCommitForm = f
			cmds = append(cmds, cmd)
			if f.State == huh.StateCompleted {
				if conf.OnApply.TagCommit {
					m.State = StateOnApplyTagFormat
				} else {
					m.State = StateBumpFiles
				}
			}
		}
	case StateOnApplyTagFormat:
		form, cmd := m.onApplyTagFormatForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.onApplyTagFormatForm = f
			cmds = append(cmds, cmd)
			if f.State == huh.StateCompleted {
				m.State = StateBumpFiles
			}
		}
	case StateBumpFiles:
		form, cmd := m.bumpFilesForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.bumpFilesForm = f
			cmds = append(cmds, cmd)
			if f.State == huh.StateCompleted {
				cmds = append(cmds, initialize(conf))
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case common.CompletedMsg:
		m.State = StateCompleted
		// return m, tea.Quit
	case common.ErrorMsg:
		m.State = StateError
		m.err = msg.Err
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.versionForm.WithWidth(msg.Width)
		m.onAddForm.WithWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.State {
			case StateCompleted, StateError:
				return m, tea.Quit
			}
		}
	}

	cmd := m.processMsg(msg)
	return m, cmd
}

func (m Model) View() string {
	conf, err := config.GetConfig()
	if err != nil {
		return err.Error()
	}
	switch m.State {
	case StateCompleted:
		return fmt.Sprintf("%#v\n", conf)
	case StateError:
		return m.err.Error()
	case StateOnAdd:
		return m.onAddForm.View()
	case StateOnApplyCommit:
		return m.onApplyCommitForm.View()
	case StateOnApplyTagCommit:
		return m.onApplyTagCommitForm.View()
	case StateOnApplyTagFormat:
		return m.onApplyTagFormatForm.View()
	case StateBumpFiles:
		return m.bumpFilesForm.View()
	default:
		return m.versionForm.View()
	}
}
