package initialize

type initializeModelState int

const (
	initializeModelStateRunning initializeModelState = iota
	initializeModelStateComplete
	initializeModelStateError
)
