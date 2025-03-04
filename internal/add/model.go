package add

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"github.com/dworthen/changelog/internal/models/common"
	"github.com/dworthen/changelog/internal/models/pipelinemodel"
)

type addModel struct {
	model tea.Model
	state addModelState
	err   error
}

func NewAddModel() (*addModel, error) {
	pipelineModel, err := pipelinemodel.NewPipelineModelBuilder().
		WithModels(formModelConstructors).
		WithStepCompletionMsg(pipelineStepCompleteMsg{}).
		WithOnComplete(onAddCompleteCmd).
		Build()

	if err != nil {
		return nil, err
	}

	return &addModel{
		model: pipelineModel,
		state: addModelStateRunning,
	}, nil
}

func (m addModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("addModel.Update: received msg", "msg", spew.Sdump(msg))

	switch msg := msg.(type) {
	case pipelineCompleteMsg:
		slog.Info("addModel.Update: PipelineCompleteMsg received. Exiting")
		return m, tea.Quit
	case common.ErrorMsg:
		slog.Error("addModel.Update: ErrorMsg received", "error", msg.Err)
		m.state = addModelStateError
		m.err = msg.Err
		return m, nil
	case tea.KeyMsg:
		slog.Debug("addModel.Update: KeyMsg received")
		switch msg.String() {
		case "ctrl+c":
			slog.Info("addModel.Update: ctrl+c received. Quitting")
			return m, tea.Quit
		case "enter":
			switch m.state {
			case addModelStateError, addModelStateComplete:
				slog.Info("addModel.Update: Enter key received at end of application. Quitting")
				return m, tea.Quit
			}
		}
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m addModel) View() string {
	if m.state == addModelStateError {
		return "Error: " + m.err.Error()
	}
	return m.model.View()
}
