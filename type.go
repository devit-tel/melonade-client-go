package melonade_client_go

type WorkflowState string

type TransactionStatus string
type WorkflowStatus string
type TaskStatus string
type TaskType string
type CommandType string
type EventType string

const (
	WORKFLOW_STATE_RUNNING   WorkflowState = "RUNNING"
	WORKFLOW_STATE_CANCELLED WorkflowState = "CANCELLED"
	WORKFLOW_STATE_FAILED    WorkflowState = "FAILED"
	WORKFLOW_STATE_PAUSED    WorkflowState = "PAUSED"
	WORKFLOW_STATE_TIMEOUT   WorkflowState = "TIMEOUT"
)

const (
	TransactionStatusRunning     TransactionStatus = "RUNNING"
	TransactionStatusPaused      TransactionStatus = "PAUSED"
	TransactionStatusCompleted   TransactionStatus = "COMPLETED"
	TransactionStatusFailed      TransactionStatus = "FAILED"
	TransactionStatusCancelled   TransactionStatus = "CANCELLED"
	TransactionStatusCompensated TransactionStatus = "COMPENSATED"
)

const (
	WorkflowStatusRunning   WorkflowStatus = "RUNNING"
	WorkflowStatusPaused    WorkflowStatus = "PAUSED"
	WorkflowStatusCompleted WorkflowStatus = "COMPLETED"
	WorkflowStatusFailed    WorkflowStatus = "FAILED"
	WorkflowStatusTimeout   WorkflowStatus = "TIMEOUT"
	WorkflowStatusCancelled WorkflowStatus = "CANCELLED"
)

const (
	TaskStatusScheduled  TaskStatus = "SCHEDULED"
	TaskStatusInProgress TaskStatus = "INPROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusFailed     TaskStatus = "FAILED"
	TaskStatusTimeout    TaskStatus = "TIMEOUT"
	TaskStatusAckTimeOut TaskStatus = "ACK_TIMEOUT"
)

const (
	TaskTypeTask       TaskType = "TASK"
	TaskTypeCompensate TaskType = "COMPENSATE"
)

const (
	CommandTypeStartTransaction CommandType = "START_TRANSACTION"
)

const (
	EventTypeTransaction EventType = "TRANSACTION"
	EventTypeWorkflow    EventType = "WORKFLOW"
	EventTypeTask        EventType = "TASK"
	EventTypeSystem      EventType = "SYSTEM"
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

type WorkflowDefinition struct {
	Name             string        `json:"name"`
	Rev              string        `json:"rev"`
	Description      string        `json:"description"`
	FailureStrategy  string        `json:"failureStrategy"`
	Tasks            []interface{} `json:"tasks"`
	OutputParameters interface{}   `json:"outputParameters"`
	Retry            struct {
		Limit int `json:"limit"`
	} `json:"retry"`
}

type TaskDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AckTimeout  int    `json:"ackTimeout"`
	Retry       struct {
		Delay int `json:"delay"`
		Limit int `json:"limit"`
	} `json:"retry"`
	SyncWorker bool `json:"syncWorker"`
	Timeout    int  `json:"timeout"`
}

type Transaction struct {
	TransactionID      string             `json:"transactionId"`
	Status             TransactionStatus  `json:"status"`
	Input              interface{}        `json:"input"`
	Output             interface{}        `json:"output"`
	CreateTime         int64              `json:"createTime"`
	EndTime            int64              `json:"endTime"`
	WorkflowDefinition WorkflowDefinition `json:"workflowDefinition"`
	Tags               []string           `json:"tags"`
}

type Workflow struct {
	TransactionID      string             `json:"transactionId"`
	Type               string             `json:"type"`
	WorkflowID         string             `json:"workflowId"`
	Status             WorkflowStatus     `json:"status"`
	Retries            int                `json:"retries"`
	Input              interface{}        `json:"input"`
	Output             interface{}        `json:"output"`
	CreateTime         int64              `json:"createTime"`
	StartTime          int64              `json:"startTime"`
	EndTime            int64              `json:"endTime"`
	WorkflowDefinition WorkflowDefinition `json:"workflowDefinition"`
	TransactionDepth   int                `json:"transactionDepth"`
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

func (t *Task) ToTaskResult() *TaskResult {
	return &TaskResult{
		TransactionID: t.TransactionID,
		TaskID:        t.TaskID,
		Status:        t.Status,
		IsSystem:      false,
	}
}

type TransactionResult struct {
	TransactionID string      `json:"transactionId"`
	Status        TaskStatus  `json:"status"`
	Output        interface{} `json:"output"`
}

type WorkflowResult struct {
	TransactionID string         `json:"transactionId"`
	WorkflowID    string         `json:"workflowId"`
	Status        WorkflowStatus `json:"status"`
	Output        interface{}    `json:"output"`
}

type TaskResult struct {
	TransactionID string        `json:"transactionId"`
	TaskID        string        `json:"taskId"`
	Status        TaskStatus    `json:"status"`
	Output        interface{}   `json:"output"`
	Logs          []interface{} `json:"logs"`       // logs to append
	DoNotRetry    bool          `json:"doNotRetry"` // If task failed do not retry
	IsSystem      bool          `json:"isSystem"`   // Internal usage
}

type WorkflowRef struct {
	Name string `json:"name"`
	Rev  string `json:"rev"`
}

type commandStartTransaction struct {
	TransactionID string       `json:"transactionId"`
	Type          CommandType  `json:"type"` // CommandTypeStartTransaction
	WorkflowRef   *WorkflowRef `json:"workflowRef,omitempty"`
	Input         interface{}  `json:"input"`
	Tags          []string     `json:"tags"`
}

type BaseEvent struct {
	TransactionID string    `json:"transactionId"`
	Type          EventType `json:"type"`    // TRANSACTION
	IsError       bool      `json:"isError"` // false
}

type EventTransaction struct {
	TransactionID string      `json:"transactionId"`
	Type          EventType   `json:"type"`    // TRANSACTION
	IsError       bool        `json:"isError"` // false
	Timestamp     int64       `json:"timestamp"`
	Details       Transaction `json:"details"`
}

type EventTransactionError struct {
	TransactionID string            `json:"transactionId"`
	Type          EventType         `json:"type"`    // TRANSACTION
	IsError       bool              `json:"isError"` // true
	Timestamp     int64             `json:"timestamp"`
	Details       TransactionResult `json:"details"`
}

type EventWorkflow struct {
	TransactionID string    `json:"transactionId"`
	Type          EventType `json:"type"`    // WORKFLOW
	IsError       bool      `json:"isError"` // false
	Timestamp     int64     `json:"timestamp"`
	Details       Workflow  `json:"details"`
}

type EventWorkflowError struct {
	TransactionID string         `json:"transactionId"`
	Type          EventType      `json:"type"`    // WORKFLOW
	IsError       bool           `json:"isError"` // true
	Timestamp     int64          `json:"timestamp"`
	Details       WorkflowResult `json:"details"`
}

type EventTask struct {
	TransactionID string    `json:"transactionId"`
	Type          EventType `json:"type"`    // TASK
	IsError       bool      `json:"isError"` // false
	Timestamp     int64     `json:"timestamp"`
	Details       Task      `json:"details"`
}

type EventTaskError struct {
	TransactionID string     `json:"transactionId"`
	Type          EventType  `json:"type"`    // TASK
	IsError       bool       `json:"isError"` // true
	Timestamp     int64      `json:"timestamp"`
	Details       TaskResult `json:"details"`
}

type EventSystem struct {
	TransactionID string      `json:"transactionId"`
	Type          EventType   `json:"type"`    // SYSTEM
	IsError       bool        `json:"isError"` // false
	Details       interface{} `json:"details"`
	Timestamp     int64       `json:"timestamp"`
}

type EventSystemError struct {
	TransactionID string      `json:"transactionId"`
	Type          EventType   `json:"type"`    // SYSTEM
	IsError       bool        `json:"isError"` // true
	Details       interface{} `json:"details"`
	Error         string      `json:"error"`
	Timestamp     int64       `json:"timestamp"`
}
