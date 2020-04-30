package melonade_client_go

import (
	"context"

	"github.com/devit-tel/goerror"
)

//go:generate mockery -name=Service
type Service interface {
	StartWorkflow(ctx context.Context, workflowName, revision, transactionId string, payload interface{}) (*StartWorkflowResponse, goerror.Error)
}
