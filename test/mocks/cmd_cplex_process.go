// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lothar1998/v2x-optimizer/cmd/v2x-optimizer-performance/cmd (interfaces: CPLEXProcess)

// Package mocks is a generated GoMock package.
package mocks

import (
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCPLEXProcess is a mock of CPLEXProcess interface.
type MockCPLEXProcess struct {
	ctrl     *gomock.Controller
	recorder *MockCPLEXProcessMockRecorder
}

// MockCPLEXProcessMockRecorder is the mock recorder for MockCPLEXProcess.
type MockCPLEXProcessMockRecorder struct {
	mock *MockCPLEXProcess
}

// NewMockCPLEXProcess creates a new mock instance.
func NewMockCPLEXProcess(ctrl *gomock.Controller) *MockCPLEXProcess {
	mock := &MockCPLEXProcess{ctrl: ctrl}
	mock.recorder = &MockCPLEXProcessMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCPLEXProcess) EXPECT() *MockCPLEXProcessMockRecorder {
	return m.recorder
}

// Output mocks base method.
func (m *MockCPLEXProcess) Output() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Output")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Output indicates an expected call of Output.
func (mr *MockCPLEXProcessMockRecorder) Output() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockCPLEXProcess)(nil).Output))
}

// Signal mocks base method.
func (m *MockCPLEXProcess) Signal(arg0 os.Signal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Signal indicates an expected call of Signal.
func (mr *MockCPLEXProcessMockRecorder) Signal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signal", reflect.TypeOf((*MockCPLEXProcess)(nil).Signal), arg0)
}
