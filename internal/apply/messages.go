package apply

type ApplyModelChangelogLoaddedMsg struct {
	changelog *changelog
}

type ApplyModelConfirmationCompleteMsg struct {
	approved bool
}

type ApplyModelAppendCommandsOutputMsg struct {
	Output string
}

type ApplyModelSetStateMsg struct {
	State ApplyModelState
}
