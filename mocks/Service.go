// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

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

// GetTaskDefinitions provides a mock function with given fields:
func (_m *Service) GetTaskDefinitions() ([]*melonade_client_go.TaskDefinition, goerror.Error) {
	ret := _m.Called()

	var r0 []*melonade_client_go.TaskDefinition
	if rf, ok := ret.Get(0).(func() []*melonade_client_go.TaskDefinition); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*melonade_client_go.TaskDefinition)
		}
	}

	var r1 goerror.Error
	if rf, ok := ret.Get(1).(func() goerror.Error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(goerror.Error)
		}
	}

	return r0, r1
}

// GetWorkflowDefinitions provides a mock function with given fields:
func (_m *Service) GetWorkflowDefinitions() ([]*melonade_client_go.WorkflowDefinition, goerror.Error) {
	ret := _m.Called()

	var r0 []*melonade_client_go.WorkflowDefinition
	if rf, ok := ret.Get(0).(func() []*melonade_client_go.WorkflowDefinition); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*melonade_client_go.WorkflowDefinition)
		}
	}

	var r1 goerror.Error
	if rf, ok := ret.Get(1).(func() goerror.Error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(goerror.Error)
		}
	}

	return r0, r1
}

// SetTaskDefinition provides a mock function with given fields: t
func (_m *Service) SetTaskDefinition(t melonade_client_go.TaskDefinition) goerror.Error {
	ret := _m.Called(t)

	var r0 goerror.Error
	if rf, ok := ret.Get(0).(func(melonade_client_go.TaskDefinition) goerror.Error); ok {
		r0 = rf(t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goerror.Error)
		}
	}

	return r0
}

// SetWorkflowDefinition provides a mock function with given fields: t
func (_m *Service) SetWorkflowDefinition(t melonade_client_go.WorkflowDefinition) goerror.Error {
	ret := _m.Called(t)

	var r0 goerror.Error
	if rf, ok := ret.Get(0).(func(melonade_client_go.WorkflowDefinition) goerror.Error); ok {
		r0 = rf(t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goerror.Error)
		}
	}

	return r0
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
