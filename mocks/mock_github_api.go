// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/compuzest/zlifecycle-il-operator/controllers/util/github (interfaces: GitAPI,RepositoryAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v32/github"
)

// MockGitAPI is a mock of GitAPI interface.
type MockGitAPI struct {
	ctrl     *gomock.Controller
	recorder *MockGitAPIMockRecorder
}

// MockGitAPIMockRecorder is the mock recorder for MockGitAPI.
type MockGitAPIMockRecorder struct {
	mock *MockGitAPI
}

// NewMockGitAPI creates a new mock instance.
func NewMockGitAPI(ctrl *gomock.Controller) *MockGitAPI {
	mock := &MockGitAPI{ctrl: ctrl}
	mock.recorder = &MockGitAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitAPI) EXPECT() *MockGitAPIMockRecorder {
	return m.recorder
}

// CreateCommit mocks base method.
func (m *MockGitAPI) CreateCommit(arg0, arg1 string, arg2 *github.Commit) (*github.Commit, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommit", arg0, arg1, arg2)
	ret0, _ := ret[0].(*github.Commit)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateCommit indicates an expected call of CreateCommit.
func (mr *MockGitAPIMockRecorder) CreateCommit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommit", reflect.TypeOf((*MockGitAPI)(nil).CreateCommit), arg0, arg1, arg2)
}

// CreateTree mocks base method.
func (m *MockGitAPI) CreateTree(arg0, arg1, arg2 string, arg3 []*github.TreeEntry) (*github.Tree, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTree", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*github.Tree)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateTree indicates an expected call of CreateTree.
func (mr *MockGitAPIMockRecorder) CreateTree(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTree", reflect.TypeOf((*MockGitAPI)(nil).CreateTree), arg0, arg1, arg2, arg3)
}

// GetCommit mocks base method.
func (m *MockGitAPI) GetCommit(arg0, arg1, arg2 string) (*github.Commit, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommit", arg0, arg1, arg2)
	ret0, _ := ret[0].(*github.Commit)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCommit indicates an expected call of GetCommit.
func (mr *MockGitAPIMockRecorder) GetCommit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommit", reflect.TypeOf((*MockGitAPI)(nil).GetCommit), arg0, arg1, arg2)
}

// GetRef mocks base method.
func (m *MockGitAPI) GetRef(arg0, arg1, arg2 string) (*github.Reference, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRef", arg0, arg1, arg2)
	ret0, _ := ret[0].(*github.Reference)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRef indicates an expected call of GetRef.
func (mr *MockGitAPIMockRecorder) GetRef(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRef", reflect.TypeOf((*MockGitAPI)(nil).GetRef), arg0, arg1, arg2)
}

// GetTree mocks base method.
func (m *MockGitAPI) GetTree(arg0, arg1, arg2 string, arg3 bool) (*github.Tree, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTree", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*github.Tree)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetTree indicates an expected call of GetTree.
func (mr *MockGitAPIMockRecorder) GetTree(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTree", reflect.TypeOf((*MockGitAPI)(nil).GetTree), arg0, arg1, arg2, arg3)
}

// UpdateRef mocks base method.
func (m *MockGitAPI) UpdateRef(arg0, arg1 string, arg2 *github.Reference, arg3 bool) (*github.Reference, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRef", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*github.Reference)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateRef indicates an expected call of UpdateRef.
func (mr *MockGitAPIMockRecorder) UpdateRef(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRef", reflect.TypeOf((*MockGitAPI)(nil).UpdateRef), arg0, arg1, arg2, arg3)
}

// MockRepositoryAPI is a mock of RepositoryAPI interface.
type MockRepositoryAPI struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryAPIMockRecorder
}

// MockRepositoryAPIMockRecorder is the mock recorder for MockRepositoryAPI.
type MockRepositoryAPIMockRecorder struct {
	mock *MockRepositoryAPI
}

// NewMockRepositoryAPI creates a new mock instance.
func NewMockRepositoryAPI(ctrl *gomock.Controller) *MockRepositoryAPI {
	mock := &MockRepositoryAPI{ctrl: ctrl}
	mock.recorder = &MockRepositoryAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryAPI) EXPECT() *MockRepositoryAPIMockRecorder {
	return m.recorder
}

// CreateHook mocks base method.
func (m *MockRepositoryAPI) CreateHook(arg0, arg1 string, arg2 *github.Hook) (*github.Hook, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHook", arg0, arg1, arg2)
	ret0, _ := ret[0].(*github.Hook)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateHook indicates an expected call of CreateHook.
func (mr *MockRepositoryAPIMockRecorder) CreateHook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHook", reflect.TypeOf((*MockRepositoryAPI)(nil).CreateHook), arg0, arg1, arg2)
}

// CreateRepository mocks base method.
func (m *MockRepositoryAPI) CreateRepository(arg0, arg1 string) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRepository", arg0, arg1)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRepository indicates an expected call of CreateRepository.
func (mr *MockRepositoryAPIMockRecorder) CreateRepository(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRepository", reflect.TypeOf((*MockRepositoryAPI)(nil).CreateRepository), arg0, arg1)
}

// DownloadContents mocks base method.
func (m *MockRepositoryAPI) DownloadContents(arg0, arg1, arg2, arg3 string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadContents", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadContents indicates an expected call of DownloadContents.
func (mr *MockRepositoryAPIMockRecorder) DownloadContents(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadContents", reflect.TypeOf((*MockRepositoryAPI)(nil).DownloadContents), arg0, arg1, arg2, arg3)
}

// GetRepository mocks base method.
func (m *MockRepositoryAPI) GetRepository(arg0, arg1 string) (*github.Repository, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepository", arg0, arg1)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRepository indicates an expected call of GetRepository.
func (mr *MockRepositoryAPIMockRecorder) GetRepository(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepository", reflect.TypeOf((*MockRepositoryAPI)(nil).GetRepository), arg0, arg1)
}

// ListHooks mocks base method.
func (m *MockRepositoryAPI) ListHooks(arg0, arg1 string, arg2 *github.ListOptions) ([]*github.Hook, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListHooks", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*github.Hook)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListHooks indicates an expected call of ListHooks.
func (mr *MockRepositoryAPIMockRecorder) ListHooks(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListHooks", reflect.TypeOf((*MockRepositoryAPI)(nil).ListHooks), arg0, arg1, arg2)
}