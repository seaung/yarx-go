package api

type CreateTaskRequest struct {
	TaskType    int    `json:"task_type"`
	TaskName    string `json:"task_name"`
	Target      string `json:"target"`
	ConCurrency int    `json:"concurrency"`
	Timeout     int    `json:"timeout"`
	Interval    int    `json:"interval"`
	MaxRetry    int    `json:"max_retry"`
}
