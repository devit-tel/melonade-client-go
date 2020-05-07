// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	goerror "github.com/devit-tel/goerror"
	melonade_client_go "github.com/devit-tel/melonade-client-go"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// StartWorkflow provides a mock function with given fields: ctx, workflowName, revision, transactionId, payload
func (_m *Service) StartWorkflow(ctx context.Context, workflowName string, revision string, transactionId string, payload interface{}) (*melonade_client_go.StartWorkflowResponse, goerror.Error) {
	ret := _m.Called(ctx, workflowName, revision, transactionId, payload)

	var r0 *melonade_client_go.StartWorkflowResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, interface{}) *melonade_client_go.StartWorkflowResponse); ok {
		r0 = rf(ctx, workflowName, revision, transactionId, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*melonade_client_go.StartWorkflowResponse)
		}
	}

	var r1 goerror.Error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, interface{}) goerror.Error); ok {
		r1 = rf(ctx, workflowName, revision, transactionId, payload)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(goerror.Error)
		}
	}

	return r0, r1
}