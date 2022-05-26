// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/compuzest/zlifecycle-il-operator/controller/common/github (interfaces: API)

// Package github is a generated GoMock package.
package github

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v42/github"
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

// CreateHook mocks base method.
func (m *MockAPI) CreateHook(arg0, arg1 string, arg2 *github.Hook) (*github.Hook, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHook", arg0, arg1, arg2)
	ret0, _ := ret[0].(*github.Hook)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateHook indicates an expected call of CreateHook.
func (mr *MockAPIMockRecorder) CreateHook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHook", reflect.TypeOf((*MockAPI)(nil).CreateHook), arg0, arg1, arg2)
}

// CreateInstallationToken mocks base method.
func (m *MockAPI) CreateInstallationToken(arg0 int64) (*github.InstallationToken, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInstallationToken", arg0)
	ret0, _ := ret[0].(*github.InstallationToken)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateInstallationToken indicates an expected call of CreateInstallationToken.
func (mr *MockAPIMockRecorder) CreateInstallationToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInstallationToken", reflect.TypeOf((*MockAPI)(nil).CreateInstallationToken), arg0)
}

// CreateRepository mocks base method.
func (m *MockAPI) CreateRepository(arg0, arg1 string) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRepository", arg0, arg1)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRepository indicates an expected call of CreateRepository.
func (mr *MockAPIMockRecorder) CreateRepository(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRepository", reflect.TypeOf((*MockAPI)(nil).CreateRepository), arg0, arg1)
}

// DownloadContents mocks base method.
func (m *MockAPI) DownloadContents(arg0, arg1, arg2, arg3 string) (io.ReadCloser, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadContents", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DownloadContents indicates an expected call of DownloadContents.
func (mr *MockAPIMockRecorder) DownloadContents(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadContents", reflect.TypeOf((*MockAPI)(nil).DownloadContents), arg0, arg1, arg2, arg3)
}

// FindOrganizationInstallation mocks base method.
func (m *MockAPI) FindOrganizationInstallation(arg0 string) (*github.Installation, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrganizationInstallation", arg0)
	ret0, _ := ret[0].(*github.Installation)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindOrganizationInstallation indicates an expected call of FindOrganizationInstallation.
func (mr *MockAPIMockRecorder) FindOrganizationInstallation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrganizationInstallation", reflect.TypeOf((*MockAPI)(nil).FindOrganizationInstallation), arg0)
}

// FindRepositoryInstallation mocks base method.
func (m *MockAPI) FindRepositoryInstallation(arg0, arg1 string) (*github.Installation, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRepositoryInstallation", arg0, arg1)
	ret0, _ := ret[0].(*github.Installation)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindRepositoryInstallation indicates an expected call of FindRepositoryInstallation.
func (mr *MockAPIMockRecorder) FindRepositoryInstallation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRepositoryInstallation", reflect.TypeOf((*MockAPI)(nil).FindRepositoryInstallation), arg0, arg1)
}

// GetRepository mocks base method.
func (m *MockAPI) GetRepository(arg0, arg1 string) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepository", arg0, arg1)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRepository indicates an expected call of GetRepository.
func (mr *MockAPIMockRecorder) GetRepository(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepository", reflect.TypeOf((*MockAPI)(nil).GetRepository), arg0, arg1)
}

// ListHooks mocks base method.
func (m *MockAPI) ListHooks(arg0, arg1 string, arg2 *github.ListOptions) ([]*github.Hook, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListHooks", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*github.Hook)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListHooks indicates an expected call of ListHooks.
func (mr *MockAPIMockRecorder) ListHooks(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListHooks", reflect.TypeOf((*MockAPI)(nil).ListHooks), arg0, arg1, arg2)
}
