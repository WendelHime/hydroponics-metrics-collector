// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/WendelHime/hydroponics-metrics-collector/internal/storage (interfaces: InfluxClient,MetricRepository)

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	models "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	gomock "go.uber.org/mock/gomock"
)

// MockInfluxClient is a mock of InfluxClient interface.
type MockInfluxClient struct {
	ctrl     *gomock.Controller
	recorder *MockInfluxClientMockRecorder
}

// MockInfluxClientMockRecorder is the mock recorder for MockInfluxClient.
type MockInfluxClientMockRecorder struct {
	mock *MockInfluxClient
}

// NewMockInfluxClient creates a new mock instance.
func NewMockInfluxClient(ctrl *gomock.Controller) *MockInfluxClient {
	mock := &MockInfluxClient{ctrl: ctrl}
	mock.recorder = &MockInfluxClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInfluxClient) EXPECT() *MockInfluxClientMockRecorder {
	return m.recorder
}

// WriteData mocks base method.
func (m *MockInfluxClient) WriteData(arg0 context.Context, arg1 string, arg2 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteData", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteData indicates an expected call of WriteData.
func (mr *MockInfluxClientMockRecorder) WriteData(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteData", reflect.TypeOf((*MockInfluxClient)(nil).WriteData), varargs...)
}

// MockMetricRepository is a mock of MetricRepository interface.
type MockMetricRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMetricRepositoryMockRecorder
}

// MockMetricRepositoryMockRecorder is the mock recorder for MockMetricRepository.
type MockMetricRepositoryMockRecorder struct {
	mock *MockMetricRepository
}

// NewMockMetricRepository creates a new mock instance.
func NewMockMetricRepository(ctrl *gomock.Controller) *MockMetricRepository {
	mock := &MockMetricRepository{ctrl: ctrl}
	mock.recorder = &MockMetricRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricRepository) EXPECT() *MockMetricRepositoryMockRecorder {
	return m.recorder
}

// WriteMeasurement mocks base method.
func (m *MockMetricRepository) WriteMeasurement(arg0 context.Context, arg1 ...models.SensorRequest) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteMeasurement", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMeasurement indicates an expected call of WriteMeasurement.
func (mr *MockMetricRepositoryMockRecorder) WriteMeasurement(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMeasurement", reflect.TypeOf((*MockMetricRepository)(nil).WriteMeasurement), varargs...)
}
