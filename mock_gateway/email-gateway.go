// Code generated by MockGen. DO NOT EDIT.
// Source: gateway/email-gateway.go

// Package mock_gateway is a generated GoMock package.
package mock_gateway

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailGateway is a mock of EmailGateway interface.
type MockEmailGateway struct {
	ctrl     *gomock.Controller
	recorder *MockEmailGatewayMockRecorder
}

// MockEmailGatewayMockRecorder is the mock recorder for MockEmailGateway.
type MockEmailGatewayMockRecorder struct {
	mock *MockEmailGateway
}

// NewMockEmailGateway creates a new mock instance.
func NewMockEmailGateway(ctrl *gomock.Controller) *MockEmailGateway {
	mock := &MockEmailGateway{ctrl: ctrl}
	mock.recorder = &MockEmailGatewayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailGateway) EXPECT() *MockEmailGatewayMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockEmailGateway) Send(to, subject, htmlString string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", to, subject, htmlString)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockEmailGatewayMockRecorder) Send(to, subject, htmlString interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockEmailGateway)(nil).Send), to, subject, htmlString)
}
