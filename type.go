package melonade_client_go

type WorkflowState string

const (
	WORKFLOW_STATE_RUNNING   WorkflowState = "RUNNING"
	WORKFLOW_STATE_CANCELLED WorkflowState = "CANCELLED"
	WORKFLOW_STATE_FAILED    WorkflowState = "FAILED"
	WORKFLOW_STATE_PAUSED    WorkflowState = "PAUSED"
	WORKFLOW_STATE_TIMEOUT   WorkflowState = "TIMEOUT"
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
