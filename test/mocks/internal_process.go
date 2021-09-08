// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lothar1998/v2x-optimizer/internal/performance/executor (interfaces: Process)

// Package mocks is a generated GoMock package.
package mocks

import (
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProcess is a mock of Process interface.
type MockProcess struct {
	ctrl     *gomock.Controller
	recorder *MockProcessMockRecorder
}

// MockProcessMockRecorder is the mock recorder for MockProcess.
type MockProcessMockRecorder struct {
	mock *MockProcess
}

// NewMockProcess creates a new mock instance.
func NewMockProcess(ctrl *gomock.Controller) *MockProcess {
	mock := &MockProcess{ctrl: ctrl}
	mock.recorder = &MockProcessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcess) EXPECT() *MockProcessMockRecorder {
	return m.recorder
}

// Output mocks base method.
func (m *MockProcess) Output() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Output")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Output indicates an expected call of Output.
func (mr *MockProcessMockRecorder) Output() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockProcess)(nil).Output))
}

// Signal mocks base method.
func (m *MockProcess) Signal(arg0 os.Signal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Signal indicates an expected call of Signal.
func (mr *MockProcessMockRecorder) Signal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signal", reflect.TypeOf((*MockProcess)(nil).Signal), arg0)
}
