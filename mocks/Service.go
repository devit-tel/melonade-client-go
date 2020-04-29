// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	melonade_client_go "gitlab.com/sendit-th/melonade-client-go"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// StartWorkflow provides a mock function with given fields: workflowName, revision, transactionId, payload
func (_m *Service) StartWorkflow(workflowName string, revision string, transactionId string, payload interface{}) (*melonade_client_go.StartWorkflowResponse, error) {
	ret := _m.Called(workflowName, revision, transactionId, payload)

	var r0 *melonade_client_go.StartWorkflowResponse
	if rf, ok := ret.Get(0).(func(string, string, string, interface{}) *melonade_client_go.StartWorkflowResponse); ok {
		r0 = rf(workflowName, revision, transactionId, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*melonade_client_go.StartWorkflowResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, interface{}) error); ok {
		r1 = rf(workflowName, revision, transactionId, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
