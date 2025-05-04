package utils

const (
	StatusPending = "PENDING"
	StatusRunning = "RUNNING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)

const (
	StatusPendingCode = 0
	StatusRunningCode = 1
	StatusSuccessCode = 2
	StatusFailedCode  = 3
)

const (
	ExecuteImmediately = 0 // 立即执行 - Execute immediately

	ExecutePeriodically = 1 // 周期执行 - Execute on schedule

	ExecuteOnSchedule = 2 // 定时执行 - Execute periodically
)
