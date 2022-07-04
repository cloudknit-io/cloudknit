// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/compuzest/zlifecycle-il-operator/controller/codegen/file (interfaces: API)

// Package file is a generated GoMock package.
package file

import (
	fs "io/fs"
	os "os"
	reflect "reflect"
	template "text/template"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
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

// CleanDir mocks base method.
func (m *MockAPI) CleanDir(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanDir", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanDir indicates an expected call of CleanDir.
func (mr *MockAPIMockRecorder) CleanDir(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanDir", reflect.TypeOf((*MockAPI)(nil).CleanDir), arg0, arg1)
}

// CopyDirContent mocks base method.
func (m *MockAPI) CopyDirContent(arg0, arg1 string, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyDirContent", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyDirContent indicates an expected call of CopyDirContent.
func (mr *MockAPIMockRecorder) CopyDirContent(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyDirContent", reflect.TypeOf((*MockAPI)(nil).CopyDirContent), arg0, arg1, arg2)
}

// CopyFile mocks base method.
func (m *MockAPI) CopyFile(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyFile indicates an expected call of CopyFile.
func (mr *MockAPIMockRecorder) CopyFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyFile", reflect.TypeOf((*MockAPI)(nil).CopyFile), arg0, arg1)
}

// CreateEmptyDirectory mocks base method.
func (m *MockAPI) CreateEmptyDirectory(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmptyDirectory", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEmptyDirectory indicates an expected call of CreateEmptyDirectory.
func (mr *MockAPIMockRecorder) CreateEmptyDirectory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmptyDirectory", reflect.TypeOf((*MockAPI)(nil).CreateEmptyDirectory), arg0)
}

// FileExistsInDir mocks base method.
func (m *MockAPI) FileExistsInDir(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileExistsInDir", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FileExistsInDir indicates an expected call of FileExistsInDir.
func (mr *MockAPIMockRecorder) FileExistsInDir(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileExistsInDir", reflect.TypeOf((*MockAPI)(nil).FileExistsInDir), arg0, arg1)
}

// IsDir mocks base method.
func (m *MockAPI) IsDir(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDir", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDir indicates an expected call of IsDir.
func (mr *MockAPIMockRecorder) IsDir(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDir", reflect.TypeOf((*MockAPI)(nil).IsDir), arg0)
}

// IsFile mocks base method.
func (m *MockAPI) IsFile(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFile", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFile indicates an expected call of IsFile.
func (mr *MockAPIMockRecorder) IsFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFile", reflect.TypeOf((*MockAPI)(nil).IsFile), arg0)
}

// NewFile mocks base method.
func (m *MockAPI) NewFile(arg0, arg1 string) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewFile", arg0, arg1)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewFile indicates an expected call of NewFile.
func (mr *MockAPIMockRecorder) NewFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewFile", reflect.TypeOf((*MockAPI)(nil).NewFile), arg0, arg1)
}

// ReadDir mocks base method.
func (m *MockAPI) ReadDir(arg0 string) ([]fs.DirEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadDir", arg0)
	ret0, _ := ret[0].([]fs.DirEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadDir indicates an expected call of ReadDir.
func (mr *MockAPIMockRecorder) ReadDir(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDir", reflect.TypeOf((*MockAPI)(nil).ReadDir), arg0)
}

// RemoveAll mocks base method.
func (m *MockAPI) RemoveAll(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *MockAPIMockRecorder) RemoveAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockAPI)(nil).RemoveAll), arg0)
}

// SaveFileFromByteArray mocks base method.
func (m *MockAPI) SaveFileFromByteArray(arg0 []byte, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromByteArray", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromByteArray indicates an expected call of SaveFileFromByteArray.
func (mr *MockAPIMockRecorder) SaveFileFromByteArray(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromByteArray", reflect.TypeOf((*MockAPI)(nil).SaveFileFromByteArray), arg0, arg1, arg2)
}

// SaveFileFromString mocks base method.
func (m *MockAPI) SaveFileFromString(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromString", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromString indicates an expected call of SaveFileFromString.
func (mr *MockAPIMockRecorder) SaveFileFromString(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromString", reflect.TypeOf((*MockAPI)(nil).SaveFileFromString), arg0, arg1, arg2)
}

// SaveFileFromTemplate mocks base method.
func (m *MockAPI) SaveFileFromTemplate(arg0 *template.Template, arg1 interface{}, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromTemplate", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromTemplate indicates an expected call of SaveFileFromTemplate.
func (mr *MockAPIMockRecorder) SaveFileFromTemplate(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromTemplate", reflect.TypeOf((*MockAPI)(nil).SaveFileFromTemplate), arg0, arg1, arg2, arg3)
}

// SaveVarsToFile mocks base method.
func (m *MockAPI) SaveVarsToFile(arg0 []*v1.Variable, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveVarsToFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveVarsToFile indicates an expected call of SaveVarsToFile.
func (mr *MockAPIMockRecorder) SaveVarsToFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveVarsToFile", reflect.TypeOf((*MockAPI)(nil).SaveVarsToFile), arg0, arg1, arg2)
}

// SaveYamlFile mocks base method.
func (m *MockAPI) SaveYamlFile(arg0 interface{}, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveYamlFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveYamlFile indicates an expected call of SaveYamlFile.
func (mr *MockAPIMockRecorder) SaveYamlFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveYamlFile", reflect.TypeOf((*MockAPI)(nil).SaveYamlFile), arg0, arg1, arg2)
}
