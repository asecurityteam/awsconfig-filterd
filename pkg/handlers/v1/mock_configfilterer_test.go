// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asecurityteam/awsconfig-filterd/pkg/domain (interfaces: ConfigFilterer)

// Package v1 is a generated GoMock package.
package v1

import (
	reflect "reflect"

	domain "github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockConfigFilterer is a mock of ConfigFilterer interface
type MockConfigFilterer struct {
	ctrl     *gomock.Controller
	recorder *MockConfigFiltererMockRecorder
}

// MockConfigFiltererMockRecorder is the mock recorder for MockConfigFilterer
type MockConfigFiltererMockRecorder struct {
	mock *MockConfigFilterer
}

// NewMockConfigFilterer creates a new mock instance
func NewMockConfigFilterer(ctrl *gomock.Controller) *MockConfigFilterer {
	mock := &MockConfigFilterer{ctrl: ctrl}
	mock.recorder = &MockConfigFiltererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigFilterer) EXPECT() *MockConfigFiltererMockRecorder {
	return m.recorder
}

// FilterConfig mocks base method
func (m *MockConfigFilterer) FilterConfig(arg0 domain.ConfigurationItem) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterConfig", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// FilterConfig indicates an expected call of FilterConfig
func (mr *MockConfigFiltererMockRecorder) FilterConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterConfig", reflect.TypeOf((*MockConfigFilterer)(nil).FilterConfig), arg0)
}
