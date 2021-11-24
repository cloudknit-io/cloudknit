// Code generated by MockGen. DO NOT EDIT.
// Source: ./file.go

// Package mocks is a generated GoMock package.
package mocks

import (
	os "os"
	reflect "reflect"
	template "text/template"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateEmptyDirectory mocks base method.
func (m *MockService) CreateEmptyDirectory(folderName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmptyDirectory", folderName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEmptyDirectory indicates an expected call of CreateEmptyDirectory.
func (mr *MockServiceMockRecorder) CreateEmptyDirectory(folderName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmptyDirectory", reflect.TypeOf((*MockService)(nil).CreateEmptyDirectory), folderName)
}

// NewFile mocks base method.
func (m *MockService) NewFile(folderName, fileName string) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewFile", folderName, fileName)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewFile indicates an expected call of NewFile.
func (mr *MockServiceMockRecorder) NewFile(folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewFile", reflect.TypeOf((*MockService)(nil).NewFile), folderName, fileName)
}

// RemoveAll mocks base method.
func (m *MockService) RemoveAll(path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", path)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *MockServiceMockRecorder) RemoveAll(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockService)(nil).RemoveAll), path)
}

// SaveFileFromByteArray mocks base method.
func (m *MockService) SaveFileFromByteArray(input []byte, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromByteArray", input, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromByteArray indicates an expected call of SaveFileFromByteArray.
func (mr *MockServiceMockRecorder) SaveFileFromByteArray(input, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromByteArray", reflect.TypeOf((*MockService)(nil).SaveFileFromByteArray), input, folderName, fileName)
}

// SaveFileFromString mocks base method.
func (m *MockService) SaveFileFromString(input, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromString", input, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromString indicates an expected call of SaveFileFromString.
func (mr *MockServiceMockRecorder) SaveFileFromString(input, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromString", reflect.TypeOf((*MockService)(nil).SaveFileFromString), input, folderName, fileName)
}

// SaveFileFromTemplate mocks base method.
func (m *MockService) SaveFileFromTemplate(t *template.Template, vars interface{}, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromTemplate", t, vars, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromTemplate indicates an expected call of SaveFileFromTemplate.
func (mr *MockServiceMockRecorder) SaveFileFromTemplate(t, vars, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromTemplate", reflect.TypeOf((*MockService)(nil).SaveFileFromTemplate), t, vars, folderName, fileName)
}

// SaveVarsToFile mocks base method.
func (m *MockService) SaveVarsToFile(variables []*v1.Variable, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveVarsToFile", variables, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveVarsToFile indicates an expected call of SaveVarsToFile.
func (mr *MockServiceMockRecorder) SaveVarsToFile(variables, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveVarsToFile", reflect.TypeOf((*MockService)(nil).SaveVarsToFile), variables, folderName, fileName)
}

// SaveYamlFile mocks base method.
func (m *MockService) SaveYamlFile(obj interface{}, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveYamlFile", obj, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveYamlFile indicates an expected call of SaveYamlFile.
func (mr *MockServiceMockRecorder) SaveYamlFile(obj, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveYamlFile", reflect.TypeOf((*MockService)(nil).SaveYamlFile), obj, folderName, fileName)
}