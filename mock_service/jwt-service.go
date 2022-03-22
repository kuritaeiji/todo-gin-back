// Code generated by MockGen. DO NOT EDIT.
// Source: service/jwt-service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockJWTService is a mock of JWTService interface.
type MockJWTService struct {
	ctrl     *gomock.Controller
	recorder *MockJWTServiceMockRecorder
}

// MockJWTServiceMockRecorder is the mock recorder for MockJWTService.
type MockJWTServiceMockRecorder struct {
	mock *MockJWTService
}

// NewMockJWTService creates a new mock instance.
func NewMockJWTService(ctrl *gomock.Controller) *MockJWTService {
	mock := &MockJWTService{ctrl: ctrl}
	mock.recorder = &MockJWTServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTService) EXPECT() *MockJWTServiceMockRecorder {
	return m.recorder
}

// CreateJWT mocks base method.
func (m *MockJWTService) CreateJWT(id, dayFromNow int) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJWT", id, dayFromNow)
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateJWT indicates an expected call of CreateJWT.
func (mr *MockJWTServiceMockRecorder) CreateJWT(id, dayFromNow interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJWT", reflect.TypeOf((*MockJWTService)(nil).CreateJWT), id, dayFromNow)
}