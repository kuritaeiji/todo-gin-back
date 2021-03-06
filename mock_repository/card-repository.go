// Code generated by MockGen. DO NOT EDIT.
// Source: repository/card-repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/kuritaeiji/todo-gin-back/model"
)

// MockCardRepository is a mock of CardRepository interface.
type MockCardRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCardRepositoryMockRecorder
}

// MockCardRepositoryMockRecorder is the mock recorder for MockCardRepository.
type MockCardRepositoryMockRecorder struct {
	mock *MockCardRepository
}

// NewMockCardRepository creates a new mock instance.
func NewMockCardRepository(ctrl *gomock.Controller) *MockCardRepository {
	mock := &MockCardRepository{ctrl: ctrl}
	mock.recorder = &MockCardRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCardRepository) EXPECT() *MockCardRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCardRepository) Create(arg0 *model.Card, arg1 *model.List) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCardRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCardRepository)(nil).Create), arg0, arg1)
}

// Destroy mocks base method.
func (m *MockCardRepository) Destroy(card *model.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", card)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockCardRepositoryMockRecorder) Destroy(card interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockCardRepository)(nil).Destroy), card)
}

// Find mocks base method.
func (m *MockCardRepository) Find(id int) (model.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", id)
	ret0, _ := ret[0].(model.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockCardRepositoryMockRecorder) Find(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockCardRepository)(nil).Find), id)
}

// Move mocks base method.
func (m *MockCardRepository) Move(card *model.Card, toListID, toIndex int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Move", card, toListID, toIndex)
	ret0, _ := ret[0].(error)
	return ret0
}

// Move indicates an expected call of Move.
func (mr *MockCardRepositoryMockRecorder) Move(card, toListID, toIndex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockCardRepository)(nil).Move), card, toListID, toIndex)
}

// Update mocks base method.
func (m *MockCardRepository) Update(card, updatingCard *model.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", card, updatingCard)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCardRepositoryMockRecorder) Update(card, updatingCard interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCardRepository)(nil).Update), card, updatingCard)
}
