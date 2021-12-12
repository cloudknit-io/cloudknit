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

// CreateEmptyDirectory mocks base method.
func (m *MockAPI) CreateEmptyDirectory(folderName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmptyDirectory", folderName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEmptyDirectory indicates an expected call of CreateEmptyDirectory.
func (mr *MockAPIMockRecorder) CreateEmptyDirectory(folderName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmptyDirectory", reflect.TypeOf((*MockAPI)(nil).CreateEmptyDirectory), folderName)
}

// NewFile mocks base method.
func (m *MockAPI) NewFile(folderName, fileName string) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewFile", folderName, fileName)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewFile indicates an expected call of NewFile.
func (mr *MockAPIMockRecorder) NewFile(folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewFile", reflect.TypeOf((*MockAPI)(nil).NewFile), folderName, fileName)
}

// RemoveAll mocks base method.
func (m *MockAPI) RemoveAll(path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", path)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *MockAPIMockRecorder) RemoveAll(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockAPI)(nil).RemoveAll), path)
}

// SaveFileFromByteArray mocks base method.
func (m *MockAPI) SaveFileFromByteArray(input []byte, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromByteArray", input, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromByteArray indicates an expected call of SaveFileFromByteArray.
func (mr *MockAPIMockRecorder) SaveFileFromByteArray(input, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromByteArray", reflect.TypeOf((*MockAPI)(nil).SaveFileFromByteArray), input, folderName, fileName)
}

// SaveFileFromString mocks base method.
func (m *MockAPI) SaveFileFromString(input, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromString", input, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromString indicates an expected call of SaveFileFromString.
func (mr *MockAPIMockRecorder) SaveFileFromString(input, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromString", reflect.TypeOf((*MockAPI)(nil).SaveFileFromString), input, folderName, fileName)
}

// SaveFileFromTemplate mocks base method.
func (m *MockAPI) SaveFileFromTemplate(t *template.Template, vars interface{}, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFileFromTemplate", t, vars, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFileFromTemplate indicates an expected call of SaveFileFromTemplate.
func (mr *MockAPIMockRecorder) SaveFileFromTemplate(t, vars, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFileFromTemplate", reflect.TypeOf((*MockAPI)(nil).SaveFileFromTemplate), t, vars, folderName, fileName)
}

// SaveVarsToFile mocks base method.
func (m *MockAPI) SaveVarsToFile(variables []*v1.Variable, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveVarsToFile", variables, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveVarsToFile indicates an expected call of SaveVarsToFile.
func (mr *MockAPIMockRecorder) SaveVarsToFile(variables, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveVarsToFile", reflect.TypeOf((*MockAPI)(nil).SaveVarsToFile), variables, folderName, fileName)
}

// SaveYamlFile mocks base method.
func (m *MockAPI) SaveYamlFile(obj interface{}, folderName, fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveYamlFile", obj, folderName, fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveYamlFile indicates an expected call of SaveYamlFile.
func (mr *MockAPIMockRecorder) SaveYamlFile(obj, folderName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveYamlFile", reflect.TypeOf((*MockAPI)(nil).SaveYamlFile), obj, folderName, fileName)
}
