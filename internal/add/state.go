package add

type addModelState int

const (
	addModelStateRunning addModelState = iota
	addModelStateComplete
	addModelStateError
)
