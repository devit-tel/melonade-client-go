package withtracer

import (
	"context"

	"github.com/devit-tel/goerror"
	json "github.com/json-iterator/go"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	wrapService "github.com/devit-tel/melonade-client-go"
)

type ServiceWithTracer struct {
	service wrapService.Service
}

func (swt *ServiceWithTracer) CreateTaskDefinition(t wrapService.TaskDefinition) goerror.Error {
	panic("implement me")
}

func (swt *ServiceWithTracer) UpdateTaskDefinition(t wrapService.TaskDefinition) goerror.Error {
	panic("implement me")
}

func (swt *ServiceWithTracer) CreateWorkflowDefinition(t wrapService.WorkflowDefinition) goerror.Error {
	panic("implement me")
}

func (swt *ServiceWithTracer) UpdateWorkflowDefinition(t wrapService.WorkflowDefinition) goerror.Error {
	panic("implement me")
}

func (swt *ServiceWithTracer) DeleteWorkflowDefinition(name, rev string) goerror.Error {
	panic("implement me")
}

func (swt *ServiceWithTracer) GetTransactionData(tID string) (*wrapService.Transaction, goerror.Error) {
	panic("implement me")
}

func Wrap(service wrapService.Service) wrapService.Service {
	return &ServiceWithTracer{
		service: service,
	}
}

func (swt *ServiceWithTracer) StartWorkflow(ctx context.Context, workflowName, revision, transactionId string, payload interface{}) (*wrapService.StartWorkflowResponse, goerror.Error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "external.wfmgateway.SendOptimizeResult")
	defer sp.Finish()

	var logPayload interface{}
	if payloadBytes, err := json.Marshal(payload); err != nil {
		logPayload = string(payloadBytes)
	} else {
		logPayload = payload
	}

	sp.LogFields(
		log.String("workflowName", workflowName),
		log.String("revision", revision),
		log.String("transactionId", transactionId),
		log.Object("payload", logPayload),
	)

	resp, err := swt.service.StartWorkflow(ctx, workflowName, revision, transactionId, payload)
	if err != nil {
		sp.LogKV("error", err)
		sp.LogKV("errorCause", err.Cause())
	}

	return resp, err
}

func (swt *ServiceWithTracer) GetWorkflowDefinitions() ([]*wrapService.WorkflowDefinition, goerror.Error) {
	return swt.service.GetWorkflowDefinitions()
}

func (swt *ServiceWithTracer) GetTaskDefinitions() ([]*wrapService.TaskDefinition, goerror.Error) {
	return swt.service.GetTaskDefinitions()
}

func (swt *ServiceWithTracer) SetTaskDefinition(t wrapService.TaskDefinition) goerror.Error {
	return swt.service.SetTaskDefinition(t)
}

func (swt *ServiceWithTracer) SetWorkflowDefinition(t wrapService.WorkflowDefinition) goerror.Error {
	return swt.service.SetWorkflowDefinition(t)
}
