package initialize

import (
	_ "embed"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"github.com/dworthen/changelog/internal/models/common"
	"github.com/dworthen/changelog/internal/models/pipelinemodel"
)

type initializeModel struct {
	model tea.Model
	state initializeModelState
	err   error
}

func NewInitializeModel() (*initializeModel, error) {

	pipelineModel, err := pipelinemodel.NewPipelineModelBuilder().
		WithModels(formModelConstructors).
		WithStepCompletionMsg(pipelineStepCompleteMsg{}).
		WithOnComplete(onInitializeCompleteCmd).
		Build()

	if err != nil {
		return nil, err
	}

	return &initializeModel{
		model: pipelineModel,
		state: initializeModelStateRunning,
	}, nil
}

func (m initializeModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m initializeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("initializeModel.Update: received msg", "msg", spew.Sdump(msg))

	switch msg := msg.(type) {
	case pipelineCompleteMsg:
		slog.Info("initializeModel.Update: PipelineCompleteMsg received. Exiting")
		return m, tea.Quit
	case common.ErrorMsg:
		slog.Info("initializeModel.Update: ErrorMsg received")
		m.state = initializeModelStateError
		m.err = msg.Err
		return m, nil
	case tea.KeyMsg:
		slog.Debug("initializeModel.Update: KeyMsg received")
		switch msg.String() {
		case "ctrl+c":
			slog.Info("initializeModel.Update: ctrl+c received. Quitting")
			return m, tea.Quit
		case "enter":
			switch m.state {
			case initializeModelStateError, initializeModelStateComplete:
				slog.Info("initializeModel.Update: Enter key received at end of application. Quitting")
				return m, tea.Quit
			}
		}
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m initializeModel) View() string {
	if m.state == initializeModelStateError {
		return "Error: " + m.err.Error()
	}
	return m.model.View()
}
