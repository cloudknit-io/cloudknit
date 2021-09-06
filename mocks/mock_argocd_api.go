// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/compuzest/zlifecycle-il-operator/controllers/argocd (interfaces: Api)

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	v1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argocd "github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	gomock "github.com/golang/mock/gomock"
)

// MockApi is a mock of Api interface.
type MockApi struct {
	ctrl     *gomock.Controller
	recorder *MockApiMockRecorder
}

// MockApiMockRecorder is the mock recorder for MockApi.
type MockApiMockRecorder struct {
	mock *MockApi
}

// NewMockApi creates a new mock instance.
func NewMockApi(ctrl *gomock.Controller) *MockApi {
	mock := &MockApi{ctrl: ctrl}
	mock.recorder = &MockApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApi) EXPECT() *MockApiMockRecorder {
	return m.recorder
}

// CreateApplication mocks base method.
func (m *MockApi) CreateApplication(arg0 *v1alpha1.Application, arg1 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateApplication", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateApplication indicates an expected call of CreateApplication.
func (mr *MockApiMockRecorder) CreateApplication(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApplication", reflect.TypeOf((*MockApi)(nil).CreateApplication), arg0, arg1)
}

// CreateProject mocks base method.
func (m *MockApi) CreateProject(arg0 argocd.CreateProjectBody, arg1 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockApiMockRecorder) CreateProject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockApi)(nil).CreateProject), arg0, arg1)
}

// CreateRepository mocks base method.
func (m *MockApi) CreateRepository(arg0 argocd.CreateRepoBody, arg1 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRepository", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRepository indicates an expected call of CreateRepository.
func (mr *MockApiMockRecorder) CreateRepository(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRepository", reflect.TypeOf((*MockApi)(nil).CreateRepository), arg0, arg1)
}

// DeleteApplication mocks base method.
func (m *MockApi) DeleteApplication(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteApplication", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteApplication indicates an expected call of DeleteApplication.
func (mr *MockApiMockRecorder) DeleteApplication(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApplication", reflect.TypeOf((*MockApi)(nil).DeleteApplication), arg0, arg1)
}

// DoesApplicationExist mocks base method.
func (m *MockApi) DoesApplicationExist(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoesApplicationExist", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoesApplicationExist indicates an expected call of DoesApplicationExist.
func (mr *MockApiMockRecorder) DoesApplicationExist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoesApplicationExist", reflect.TypeOf((*MockApi)(nil).DoesApplicationExist), arg0, arg1)
}

// DoesProjectExist mocks base method.
func (m *MockApi) DoesProjectExist(arg0, arg1 string) (bool, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoesProjectExist", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DoesProjectExist indicates an expected call of DoesProjectExist.
func (mr *MockApiMockRecorder) DoesProjectExist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoesProjectExist", reflect.TypeOf((*MockApi)(nil).DoesProjectExist), arg0, arg1)
}

// GetAuthToken mocks base method.
func (m *MockApi) GetAuthToken() (*argocd.GetTokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthToken")
	ret0, _ := ret[0].(*argocd.GetTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthToken indicates an expected call of GetAuthToken.
func (mr *MockApiMockRecorder) GetAuthToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthToken", reflect.TypeOf((*MockApi)(nil).GetAuthToken))
}

// ListRepositories mocks base method.
func (m *MockApi) ListRepositories(arg0 string) (*argocd.RepositoryList, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRepositories", arg0)
	ret0, _ := ret[0].(*argocd.RepositoryList)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListRepositories indicates an expected call of ListRepositories.
func (mr *MockApiMockRecorder) ListRepositories(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRepositories", reflect.TypeOf((*MockApi)(nil).ListRepositories), arg0)
}
