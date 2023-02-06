// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/common-fate/common-fate/pkg/api (interfaces: TargetGroupService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	targetgroup "github.com/common-fate/common-fate/pkg/targetgroup"
	types "github.com/common-fate/common-fate/pkg/types"
	gomock "github.com/golang/mock/gomock"
)

// MockTargetGroupService is a mock of TargetGroupService interface.
type MockTargetGroupService struct {
	ctrl     *gomock.Controller
	recorder *MockTargetGroupServiceMockRecorder
}

// MockTargetGroupServiceMockRecorder is the mock recorder for MockTargetGroupService.
type MockTargetGroupServiceMockRecorder struct {
	mock *MockTargetGroupService
}

// NewMockTargetGroupService creates a new mock instance.
func NewMockTargetGroupService(ctrl *gomock.Controller) *MockTargetGroupService {
	mock := &MockTargetGroupService{ctrl: ctrl}
	mock.recorder = &MockTargetGroupServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTargetGroupService) EXPECT() *MockTargetGroupServiceMockRecorder {
	return m.recorder
}

// CreateTargetGroup mocks base method.
func (m *MockTargetGroupService) CreateTargetGroup(arg0 context.Context, arg1 types.CreateTargetGroupRequest) (*targetgroup.TargetGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTargetGroup", arg0, arg1)
	ret0, _ := ret[0].(*targetgroup.TargetGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTargetGroup indicates an expected call of CreateTargetGroup.
func (mr *MockTargetGroupServiceMockRecorder) CreateTargetGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTargetGroup", reflect.TypeOf((*MockTargetGroupService)(nil).CreateTargetGroup), arg0, arg1)
}
