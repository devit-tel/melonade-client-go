package melonade_client_go

//go:generate mockery -name=Service
type Service interface {
	StartWorkflow(workflowName, revision, transactionId string, payload interface{}) (*StartWorkflowResponse, error)
}
