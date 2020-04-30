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
