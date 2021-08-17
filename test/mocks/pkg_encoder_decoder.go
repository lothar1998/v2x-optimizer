// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lothar1998/v2x-optimizer/pkg/data (interfaces: EncoderDecoder)

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	data "github.com/lothar1998/v2x-optimizer/pkg/data"
)

// MockEncoderDecoder is a mock of EncoderDecoder interface.
type MockEncoderDecoder struct {
	ctrl     *gomock.Controller
	recorder *MockEncoderDecoderMockRecorder
}

// MockEncoderDecoderMockRecorder is the mock recorder for MockEncoderDecoder.
type MockEncoderDecoderMockRecorder struct {
	mock *MockEncoderDecoder
}

// NewMockEncoderDecoder creates a new mock instance.
func NewMockEncoderDecoder(ctrl *gomock.Controller) *MockEncoderDecoder {
	mock := &MockEncoderDecoder{ctrl: ctrl}
	mock.recorder = &MockEncoderDecoderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEncoderDecoder) EXPECT() *MockEncoderDecoderMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockEncoderDecoder) Decode(arg0 io.Reader) (*data.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", arg0)
	ret0, _ := ret[0].(*data.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decode indicates an expected call of Decode.
func (mr *MockEncoderDecoderMockRecorder) Decode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockEncoderDecoder)(nil).Decode), arg0)
}

// Encode mocks base method.
func (m *MockEncoderDecoder) Encode(arg0 *data.Data, arg1 io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Encode indicates an expected call of Encode.
func (mr *MockEncoderDecoderMockRecorder) Encode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockEncoderDecoder)(nil).Encode), arg0, arg1)
}
