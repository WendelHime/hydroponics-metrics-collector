// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/WendelHime/hydroponics-metrics-collector/internal/logic (interfaces: UserLogic)

// Package logic is a generated GoMock package.
package logic

import (
	context "context"
	reflect "reflect"

	models "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	gomock "go.uber.org/mock/gomock"
)

// MockUserLogic is a mock of UserLogic interface.
type MockUserLogic struct {
	ctrl     *gomock.Controller
	recorder *MockUserLogicMockRecorder
}

// MockUserLogicMockRecorder is the mock recorder for MockUserLogic.
type MockUserLogicMockRecorder struct {
	mock *MockUserLogic
}

// NewMockUserLogic creates a new mock instance.
func NewMockUserLogic(ctrl *gomock.Controller) *MockUserLogic {
	mock := &MockUserLogic{ctrl: ctrl}
	mock.recorder = &MockUserLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserLogic) EXPECT() *MockUserLogicMockRecorder {
	return m.recorder
}

// AddDevice mocks base method.
func (m *MockUserLogic) AddDevice(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDevice", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDevice indicates an expected call of AddDevice.
func (mr *MockUserLogicMockRecorder) AddDevice(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDevice", reflect.TypeOf((*MockUserLogic)(nil).AddDevice), arg0, arg1, arg2)
}

// CreateAccount mocks base method.
func (m *MockUserLogic) CreateAccount(arg0 context.Context, arg1 models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockUserLogicMockRecorder) CreateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockUserLogic)(nil).CreateAccount), arg0, arg1)
}

// GetDevices mocks base method.
func (m *MockUserLogic) GetDevices(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevices", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevices indicates an expected call of GetDevices.
func (mr *MockUserLogicMockRecorder) GetDevices(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevices", reflect.TypeOf((*MockUserLogic)(nil).GetDevices), arg0, arg1)
}

// Login mocks base method.
func (m *MockUserLogic) Login(arg0 context.Context, arg1 models.Credentials) (models.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(models.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserLogicMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserLogic)(nil).Login), arg0, arg1)
}
