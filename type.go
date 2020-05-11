package melonade_client_go

type WorkflowState string

type TaskStatus string
type TaskType string
type CommandType string

const (
	WORKFLOW_STATE_RUNNING   WorkflowState = "RUNNING"
	WORKFLOW_STATE_CANCELLED WorkflowState = "CANCELLED"
	WORKFLOW_STATE_FAILED    WorkflowState = "FAILED"
	WORKFLOW_STATE_PAUSED    WorkflowState = "PAUSED"
	WORKFLOW_STATE_TIMEOUT   WorkflowState = "TIMEOUT"
)

const (
	TaskStatusScheduled  TaskStatus = "SCHEDULED"
	TaskStatusInProgress TaskStatus = "INPROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusFailed     TaskStatus = "FAILED"
)

const (
	TaskTypeTask       TaskType = "TASK"
	TaskTypeCompensate TaskType = "COMPENSATE"
)

const (
	CommandTypeStartTransaction CommandType = "START_TRANSACTION"
)

type ResponseData struct {
	TransactionId string        `json:"transactionId"`
	Status        WorkflowState `json:"status"`
	Input         interface{}   `json:"input"`
	Output        interface{}   `json:"output"`
	CreateTime    int64         `json:"createTime"`
	EndTime       int64         `json:"endTime"`
}

type ResponseDebug struct {
	Method  string      `json:"method"`
	Url     string      `json:"url"`
	Headers interface{} `json:"headers"`
	Body    interface{} `json:"body"`
	Query   interface{} `json:"query"`
	Params  interface{} `json:"params"`
}

type ResponseError struct {
	Message string         `json:"message"`
	Stack   string         `json:"stack"`
	Debug   *ResponseDebug `json:"debug"`
}

type StartWorkflowResponse struct {
	Success bool           `json:"success"`
	Data    *ResponseData  `json:"data"`
	Error   *ResponseError `json:"error"`
}

type Task struct {
	TaskID            string        `json:"taskId"`
	TaskName          string        `json:"taskName"`
	TaskReferenceName string        `json:"taskReferenceName"`
	WorkflowID        string        `json:"workflowId"`
	TransactionID     string        `json:"transactionId"`
	Type              TaskType      `json:"type"`
	Status            TaskStatus    `json:"status"`
	IsRetried         bool          `json:"isRetried"`
	Input             interface{}   `json:"input"`
	CreateTime        int64         `json:"createTime"`
	StartTime         int64         `json:"startTime"`
	EndTime           int64         `json:"endTime"`
	Retries           int           `json:"retries"`
	RetryDelay        int           `json:"retryDelay"`
	AckTimeout        int           `json:"ackTimeout"`
	Timeout           int           `json:"timeout"`
	TaskPath          []interface{} `json:"taskPath"` // slice of int + string e.g. [0, "parallelTasks", 0, 1]
	Output            interface{}   `json:"output"`
	Logs              []interface{} `json:"logs"`
}

type TaskResult struct {
	TransactionID string        `json:"transactionId"`
	TaskID        string        `json:"taskId"`
	Status        TaskStatus    `json:"status"`
	Output        interface{}   `json:"output"`
	Logs          []interface{} `json:"logs"`       // logs to append
	DoNotRetry    bool          `json:"doNotRetry"` // If task failed do not retry
}

type workflowRef struct {
	Name string `json:"name"`
	Rev  string `json:"rev"`
}

type commandStartTransaction struct {
	TransactionID string       `json:"transactionId"`
	Type          CommandType  `json:"type"` // CommandTypeStartTransaction
	WorkflowRef   *workflowRef `json:"workflowRef,omitempty"`
	Input         interface{}  `json:"input"`
	Tags          []string     `json:"tags"`
}
