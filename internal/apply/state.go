package apply

type ApplyModelState string

const (
	ApplyModelStateLoading                     ApplyModelState = "Loading..."
	ApplyModelStateNoChanges                   ApplyModelState = "No changes. Press enter to exit."
	ApplyModelStateReviewingScrollWindowActive ApplyModelState = "Reviewing changes."
	ApplyModelStateReviewingConfirmActive      ApplyModelState = "Approve changes."
	ApplyModelStateApplying                    ApplyModelState = "Applying changes..."
	ApplyModelStateRunningCommands             ApplyModelState = "Running commands..."
	ApplyModelStateError                       ApplyModelState = "Error. Press any key to exit."
	ApplyModelStateComplete                    ApplyModelState = "Changes applied. Press any key to exit."
	ApplyModelStateCancelled                   ApplyModelState = "Cancelled. Press any key to exit."
)
