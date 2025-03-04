package apply

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/dworthen/changelog/internal/models/common"
	"github.com/dworthen/changelog/internal/models/formwrapmodel"
	"github.com/dworthen/changelog/internal/models/helppanel"
	"github.com/dworthen/changelog/internal/models/scrollwindow"
)

var theme = huh.ThemeBase16()

type applyModel struct {
	State             ApplyModelState
	changelog         *changelog
	err               error
	scrollWindowModel tea.Model
	confirmModel      tea.Model
	helpPanelModel    tea.Model
	commandOutputs    string
	width             int
	height            int
}

func NewApplyModel() (*applyModel, error) {

	scrollWindowModel := scrollwindow.NewScrollWindowModelBuilder().WithTitle("Summary of Changes").Build()

	confirmModel, err := newApproveChangelogForm()
	if err != nil {
		return nil, err
	}

	keys := []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Next Panel")),
	}
	helpPanelModel := helppanel.NewHelpPanelModel(keys)

	return &applyModel{
		State:             ApplyModelStateLoading,
		scrollWindowModel: scrollWindowModel,
		confirmModel:      confirmModel,
		helpPanelModel:    helpPanelModel,
	}, nil
}

var (
	borderHeight = 4
	borderWidth  = 2
	modelStyle   = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

func (m applyModel) Init() tea.Cmd {
	return tea.Batch(m.scrollWindowModel.Init(), loadChangelogCmd(), m.confirmModel.Init())
}

func (m applyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("applyModel.Update: received msg", "msg", spew.Sdump(msg))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		slog.Info("applyModel.Update: WindowSizeMsg received", "msg", spew.Sdump(msg))
		m.width = msg.Width
		m.height = msg.Height
		return m, scrollwindow.ScrollWindowResizeCmd(msg.Width-borderWidth, m.getScrollWindowHeight())
	case ApplyModelSetStateMsg:
		slog.Info("applyModel.Update: ApplyModelSetStateMsg received", "msg", spew.Sdump(msg))
		m.State = msg.State
		cmds := []tea.Cmd{scrollwindow.ScrollWindowResizeCmd(m.width-borderWidth, m.getScrollWindowHeight())}
		if m.State == ApplyModelStateRunningCommands {
			cmds = append(cmds, scrollwindow.ScrollWindowSetTitleCmd("Running Commands"))
		}
		return m, tea.Batch(cmds...)
	case ApplyModelChangelogLoaddedMsg:
		slog.Info("applyModel.Update: ChangelogLoadedMsg received", "msg", spew.Sdump(msg))
		m.changelog = msg.changelog
		if m.changelog.BumpingVersion {
			m.State = ApplyModelStateReviewingScrollWindowActive
		} else {
			m.State = ApplyModelStateNoChanges
		}
		contents, err := m.changelog.getSummary()
		if err != nil {
			return m, common.ErrorMsgCmd(err)
		}
		cmds := []tea.Cmd{scrollwindow.ScrollWindowSetContentCmd(contents), scrollwindow.ScrollWindowResizeCmd(m.width-borderWidth, m.getScrollWindowHeight())}
		return m, tea.Batch(cmds...)
	case ApplyModelAppendCommandsOutputMsg:
		m.commandOutputs += msg.Output
		return m, scrollwindow.ScrollWindowSetContentCmd(m.commandOutputs)
	case common.ErrorMsg:
		slog.Error("applyModel.Update: ErrorMsg received", "msg", spew.Sdump(msg))
		m.err = msg.Err
		m.State = ApplyModelStateError
		return m, nil
	case tea.KeyMsg:
		if m.State == ApplyModelStateError ||
			m.State == ApplyModelStateComplete ||
			m.State == ApplyModelStateCancelled {
			return m, tea.Quit
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			if m.State == ApplyModelStateReviewingScrollWindowActive {
				m.State = ApplyModelStateReviewingConfirmActive
				return m, formwrapmodel.FormWrapWithHelpCmd(true)
			} else if m.State == ApplyModelStateReviewingConfirmActive {
				m.State = ApplyModelStateReviewingScrollWindowActive
				return m, formwrapmodel.FormWrapWithHelpCmd(false)
			}
		case "enter":
			if m.State == ApplyModelStateNoChanges {
				return m, tea.Quit
			}
		}
	case helppanel.HelpPanelWithKeysMsg:
		newHelpPanel, cmd := m.helpPanelModel.Update(msg)
		m.helpPanelModel = newHelpPanel
		return m, cmd
	case formwrapmodel.FormWrapWithHelpMsg:
		slog.Info("applyModel.Update: FormWrapWithHelpMsg received", "msg", spew.Sdump(msg))
		newForm, cmd := m.confirmModel.Update(msg)
		m.confirmModel = newForm
		return m, cmd
	case ApplyModelConfirmationCompleteMsg:
		slog.Info("applyModel.Update: ChangelogConfirmationCompleteMsg received", "msg", spew.Sdump(msg))
		if msg.approved {
			m.State = ApplyModelStateApplying
			return m, applyModelCompleteCmd(m.changelog)
		}
		m.State = ApplyModelStateCancelled
		return m, nil
	case tea.MouseMsg, scrollwindow.ScrollWindowResizeMsg, scrollwindow.ScrollWindowSetContentMsg, scrollwindow.ScrollWiindowSetTitleMsg:
		newScrollWindow, cmd := m.scrollWindowModel.Update(msg)
		m.scrollWindowModel = newScrollWindow
		return m, cmd
	}

	if m.State == ApplyModelStateReviewingConfirmActive {
		newConfirm, cmd := m.confirmModel.Update(msg)
		m.confirmModel = newConfirm
		return m, cmd
	} else {
		newScrollWindow, cmd := m.scrollWindowModel.Update(msg)
		m.scrollWindowModel = newScrollWindow
		return m, cmd
	}
}

func (m applyModel) ShowTitleView() bool {
	return m.State == ApplyModelStateLoading ||
		m.State == ApplyModelStateNoChanges ||
		m.State == ApplyModelStateApplying ||
		m.State == ApplyModelStateError ||
		m.State == ApplyModelStateComplete ||
		m.State == ApplyModelStateCancelled
}

func (m applyModel) TitleView() string {
	return theme.Focused.Title.Render(string(m.State))
}

func (m applyModel) GetTitleViewHeight() int {
	if !m.ShowTitleView() {
		return 0
	}
	return lipgloss.Height(m.TitleView())
}

func (m applyModel) ShowHelpPanel() bool {
	return m.State == ApplyModelStateReviewingScrollWindowActive ||
		m.State == ApplyModelStateReviewingConfirmActive
}

func (m applyModel) helpPanelView() string {
	return m.helpPanelModel.View()
}

func (m applyModel) getHelpPanelHeight() int {
	if !m.ShowHelpPanel() {
		return 0
	}
	return lipgloss.Height(m.helpPanelView())
}

func (m applyModel) ShowConfirmView() bool {
	return m.State == ApplyModelStateReviewingConfirmActive ||
		m.State == ApplyModelStateReviewingScrollWindowActive
}

func (m applyModel) confirmView() string {
	confirmHeight := m.getConfirmHeight()
	if m.State == ApplyModelStateReviewingConfirmActive {
		return focusedModelStyle.Width(m.width - borderWidth).Height(confirmHeight).Render(m.confirmModel.View())
	}
	return modelStyle.Width(m.width - borderWidth).Height(confirmHeight).Render(m.confirmModel.View())
}

func (m applyModel) getConfirmHeight() int {
	if !m.ShowConfirmView() {
		return 0
	}
	return lipgloss.Height(m.confirmModel.View())
}

func (m applyModel) ShowScrollWindow() bool {
	return m.State == ApplyModelStateNoChanges ||
		m.State == ApplyModelStateReviewingScrollWindowActive ||
		m.State == ApplyModelStateReviewingConfirmActive ||
		m.State == ApplyModelStateRunningCommands ||
		m.State == ApplyModelStateApplying ||
		m.State == ApplyModelStateComplete
}

func (m applyModel) ScrollWindowView() string {
	scrollHeight := m.getScrollWindowHeight()
	if m.State == ApplyModelStateReviewingScrollWindowActive {
		return focusedModelStyle.Width(m.width - borderWidth).Height(scrollHeight).Render(m.scrollWindowModel.View())
	}
	return modelStyle.Width(m.width - borderWidth).Height(scrollHeight).Render(m.scrollWindowModel.View())
}

func (m applyModel) getScrollWindowHeight() int {
	return m.height - m.GetTitleViewHeight() - m.getConfirmHeight() - m.getHelpPanelHeight() - borderHeight
}

func (m applyModel) View() string {
	views := []string{}

	if m.ShowTitleView() {
		views = append(views, m.TitleView())
	}
	if m.State == ApplyModelStateError {
		views = append(views, m.err.Error())
		return lipgloss.JoinVertical(lipgloss.Top, views...)
	}

	if m.ShowScrollWindow() {
		views = append(views, m.ScrollWindowView())
	}
	if m.ShowConfirmView() {
		views = append(views, m.confirmView())
	}
	if m.ShowHelpPanel() {
		views = append(views, m.helpPanelView())
	}
	return lipgloss.JoinVertical(lipgloss.Top, views...)
}
