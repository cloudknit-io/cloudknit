// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/compuzest/zlifecycle-il-operator/controller/lib/watcher (interfaces: API)

// Package watcher is a generated GoMock package.
package watcher

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAPI is a mock of API interface.
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
}

// MockAPIMockRecorder is the mock recorder for MockAPI.
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance.
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// Watch mocks base method.
func (m *MockAPI) Watch(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Watch indicates an expected call of Watch.
func (mr *MockAPIMockRecorder) Watch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockAPI)(nil).Watch), arg0)
}
