package known

import "github.com/RichardKnop/machinery/v2/tasks"

const (
	CREATED string = "CREATED"
	REVOKED string = "REVOKED"
	RUNNING string = tasks.StateStarted
	PENDING string = tasks.StatePending
	SUCCESS string = tasks.StateSuccess
	FAILURE string = tasks.StateFailure
)
