// Code generated by MockGen. DO NOT EDIT.
// Source: storj.io/storj/pkg/eestream (interfaces: ErasureScheme)

// Package mock_eestream is a generated GoMock package.
package mock_eestream

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockErasureScheme is a mock of ErasureScheme interface
type MockErasureScheme struct {
	ctrl     *gomock.Controller
	recorder *MockErasureSchemeMockRecorder
}

// MockErasureSchemeMockRecorder is the mock recorder for MockErasureScheme
type MockErasureSchemeMockRecorder struct {
	mock *MockErasureScheme
}

// NewMockErasureScheme creates a new mock instance
func NewMockErasureScheme(ctrl *gomock.Controller) *MockErasureScheme {
	mock := &MockErasureScheme{ctrl: ctrl}
	mock.recorder = &MockErasureSchemeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockErasureScheme) EXPECT() *MockErasureSchemeMockRecorder {
	return m.recorder
}

// Decode mocks base method
func (m *MockErasureScheme) Decode(arg0 []byte, arg1 map[int][]byte) ([]byte, error) {
	ret := m.ctrl.Call(m, "Decode", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decode indicates an expected call of Decode
func (mr *MockErasureSchemeMockRecorder) Decode(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockErasureScheme)(nil).Decode), arg0, arg1)
}

// Encode mocks base method
func (m *MockErasureScheme) Encode(arg0 []byte, arg1 func(int, []byte)) error {
	ret := m.ctrl.Call(m, "Encode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Encode indicates an expected call of Encode
func (mr *MockErasureSchemeMockRecorder) Encode(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockErasureScheme)(nil).Encode), arg0, arg1)
}

// EncodeSingle mocks base method
func (m *MockErasureScheme) EncodeSingle(arg0, arg1 []byte, arg2 int) error {
	ret := m.ctrl.Call(m, "EncodeSingle", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// EncodeSingle indicates an expected call of EncodeSingle
func (mr *MockErasureSchemeMockRecorder) EncodeSingle(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodeSingle", reflect.TypeOf((*MockErasureScheme)(nil).EncodeSingle), arg0, arg1, arg2)
}

// ErasureShareSize mocks base method
func (m *MockErasureScheme) ErasureShareSize() int {
	ret := m.ctrl.Call(m, "ErasureShareSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// ErasureShareSize indicates an expected call of ErasureShareSize
func (mr *MockErasureSchemeMockRecorder) ErasureShareSize() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErasureShareSize", reflect.TypeOf((*MockErasureScheme)(nil).ErasureShareSize))
}

// RequiredCount mocks base method
func (m *MockErasureScheme) RequiredCount() int {
	ret := m.ctrl.Call(m, "RequiredCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// RequiredCount indicates an expected call of RequiredCount
func (mr *MockErasureSchemeMockRecorder) RequiredCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequiredCount", reflect.TypeOf((*MockErasureScheme)(nil).RequiredCount))
}

// StripeSize mocks base method
func (m *MockErasureScheme) StripeSize() int {
	ret := m.ctrl.Call(m, "StripeSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// StripeSize indicates an expected call of StripeSize
func (mr *MockErasureSchemeMockRecorder) StripeSize() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StripeSize", reflect.TypeOf((*MockErasureScheme)(nil).StripeSize))
}

// TotalCount mocks base method
func (m *MockErasureScheme) TotalCount() int {
	ret := m.ctrl.Call(m, "TotalCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// TotalCount indicates an expected call of TotalCount
func (mr *MockErasureSchemeMockRecorder) TotalCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalCount", reflect.TypeOf((*MockErasureScheme)(nil).TotalCount))
}