// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/common-fate/granted-approvals/pkg/service/grantsvc (interfaces: AHClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"

	types "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	gomock "github.com/golang/mock/gomock"
)

// MockAHClient is a mock of AHClient interface.
type MockAHClient struct {
	ctrl     *gomock.Controller
	recorder *MockAHClientMockRecorder
}

// MockAHClientMockRecorder is the mock recorder for MockAHClient.
type MockAHClientMockRecorder struct {
	mock *MockAHClient
}

// NewMockAHClient creates a new mock instance.
func NewMockAHClient(ctrl *gomock.Controller) *MockAHClient {
	mock := &MockAHClient{ctrl: ctrl}
	mock.recorder = &MockAHClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAHClient) EXPECT() *MockAHClientMockRecorder {
	return m.recorder
}

// FetchArgGroupValuesWithBodyWithResponse mocks base method.
func (m *MockAHClient) FetchArgGroupValuesWithBodyWithResponse(arg0 context.Context, arg1, arg2, arg3 string, arg4 io.Reader, arg5 ...types.RequestEditorFn) (*types.FetchArgGroupValuesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3, arg4}
	for _, a := range arg5 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchArgGroupValuesWithBodyWithResponse", varargs...)
	ret0, _ := ret[0].(*types.FetchArgGroupValuesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchArgGroupValuesWithBodyWithResponse indicates an expected call of FetchArgGroupValuesWithBodyWithResponse.
func (mr *MockAHClientMockRecorder) FetchArgGroupValuesWithBodyWithResponse(arg0, arg1, arg2, arg3, arg4 interface{}, arg5 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3, arg4}, arg5...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchArgGroupValuesWithBodyWithResponse", reflect.TypeOf((*MockAHClient)(nil).FetchArgGroupValuesWithBodyWithResponse), varargs...)
}

// FetchArgGroupValuesWithResponse mocks base method.
func (m *MockAHClient) FetchArgGroupValuesWithResponse(arg0 context.Context, arg1, arg2 string, arg3 types.FetchArgGroupValuesJSONRequestBody, arg4 ...types.RequestEditorFn) (*types.FetchArgGroupValuesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchArgGroupValuesWithResponse", varargs...)
	ret0, _ := ret[0].(*types.FetchArgGroupValuesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchArgGroupValuesWithResponse indicates an expected call of FetchArgGroupValuesWithResponse.
func (mr *MockAHClientMockRecorder) FetchArgGroupValuesWithResponse(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchArgGroupValuesWithResponse", reflect.TypeOf((*MockAHClient)(nil).FetchArgGroupValuesWithResponse), varargs...)
}

// GetAccessInstructionsWithResponse mocks base method.
func (m *MockAHClient) GetAccessInstructionsWithResponse(arg0 context.Context, arg1 string, arg2 *types.GetAccessInstructionsParams, arg3 ...types.RequestEditorFn) (*types.GetAccessInstructionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccessInstructionsWithResponse", varargs...)
	ret0, _ := ret[0].(*types.GetAccessInstructionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccessInstructionsWithResponse indicates an expected call of GetAccessInstructionsWithResponse.
func (mr *MockAHClientMockRecorder) GetAccessInstructionsWithResponse(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessInstructionsWithResponse", reflect.TypeOf((*MockAHClient)(nil).GetAccessInstructionsWithResponse), varargs...)
}

// GetGrantsWithResponse mocks base method.
func (m *MockAHClient) GetGrantsWithResponse(arg0 context.Context, arg1 ...types.RequestEditorFn) (*types.GetGrantsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetGrantsWithResponse", varargs...)
	ret0, _ := ret[0].(*types.GetGrantsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGrantsWithResponse indicates an expected call of GetGrantsWithResponse.
func (mr *MockAHClientMockRecorder) GetGrantsWithResponse(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGrantsWithResponse", reflect.TypeOf((*MockAHClient)(nil).GetGrantsWithResponse), varargs...)
}

// GetHealthWithResponse mocks base method.
func (m *MockAHClient) GetHealthWithResponse(arg0 context.Context, arg1 ...types.RequestEditorFn) (*types.GetHealthResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetHealthWithResponse", varargs...)
	ret0, _ := ret[0].(*types.GetHealthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHealthWithResponse indicates an expected call of GetHealthWithResponse.
func (mr *MockAHClientMockRecorder) GetHealthWithResponse(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHealthWithResponse", reflect.TypeOf((*MockAHClient)(nil).GetHealthWithResponse), varargs...)
}

// GetProviderArgsWithResponse mocks base method.
func (m *MockAHClient) GetProviderArgsWithResponse(arg0 context.Context, arg1 string, arg2 ...types.RequestEditorFn) (*types.GetProviderArgsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProviderArgsWithResponse", varargs...)
	ret0, _ := ret[0].(*types.GetProviderArgsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderArgsWithResponse indicates an expected call of GetProviderArgsWithResponse.
func (mr *MockAHClientMockRecorder) GetProviderArgsWithResponse(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderArgsWithResponse", reflect.TypeOf((*MockAHClient)(nil).GetProviderArgsWithResponse), varargs...)
}

// GetProviderWithResponse mocks base method.
func (m *MockAHClient) GetProviderWithResponse(arg0 context.Context, arg1 string, arg2 ...types.RequestEditorFn) (*types.GetProviderResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProviderWithResponse", varargs...)
	ret0, _ := ret[0].(*types.GetProviderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderWithResponse indicates an expected call of GetProviderWithResponse.
func (mr *MockAHClientMockRecorder) GetProviderWithResponse(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderWithResponse", reflect.TypeOf((*MockAHClient)(nil).GetProviderWithResponse), varargs...)
}

// ListProviderArgOptionsWithResponse mocks base method.
func (m *MockAHClient) ListProviderArgOptionsWithResponse(arg0 context.Context, arg1, arg2 string, arg3 ...types.RequestEditorFn) (*types.ListProviderArgOptionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListProviderArgOptionsWithResponse", varargs...)
	ret0, _ := ret[0].(*types.ListProviderArgOptionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProviderArgOptionsWithResponse indicates an expected call of ListProviderArgOptionsWithResponse.
func (mr *MockAHClientMockRecorder) ListProviderArgOptionsWithResponse(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProviderArgOptionsWithResponse", reflect.TypeOf((*MockAHClient)(nil).ListProviderArgOptionsWithResponse), varargs...)
}

// ListProvidersWithResponse mocks base method.
func (m *MockAHClient) ListProvidersWithResponse(arg0 context.Context, arg1 ...types.RequestEditorFn) (*types.ListProvidersResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListProvidersWithResponse", varargs...)
	ret0, _ := ret[0].(*types.ListProvidersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProvidersWithResponse indicates an expected call of ListProvidersWithResponse.
func (mr *MockAHClientMockRecorder) ListProvidersWithResponse(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProvidersWithResponse", reflect.TypeOf((*MockAHClient)(nil).ListProvidersWithResponse), varargs...)
}

// PostGrantsRevokeWithBodyWithResponse mocks base method.
func (m *MockAHClient) PostGrantsRevokeWithBodyWithResponse(arg0 context.Context, arg1, arg2 string, arg3 io.Reader, arg4 ...types.RequestEditorFn) (*types.PostGrantsRevokeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PostGrantsRevokeWithBodyWithResponse", varargs...)
	ret0, _ := ret[0].(*types.PostGrantsRevokeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostGrantsRevokeWithBodyWithResponse indicates an expected call of PostGrantsRevokeWithBodyWithResponse.
func (mr *MockAHClientMockRecorder) PostGrantsRevokeWithBodyWithResponse(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostGrantsRevokeWithBodyWithResponse", reflect.TypeOf((*MockAHClient)(nil).PostGrantsRevokeWithBodyWithResponse), varargs...)
}

// PostGrantsRevokeWithResponse mocks base method.
func (m *MockAHClient) PostGrantsRevokeWithResponse(arg0 context.Context, arg1 string, arg2 types.PostGrantsRevokeJSONRequestBody, arg3 ...types.RequestEditorFn) (*types.PostGrantsRevokeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PostGrantsRevokeWithResponse", varargs...)
	ret0, _ := ret[0].(*types.PostGrantsRevokeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostGrantsRevokeWithResponse indicates an expected call of PostGrantsRevokeWithResponse.
func (mr *MockAHClientMockRecorder) PostGrantsRevokeWithResponse(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostGrantsRevokeWithResponse", reflect.TypeOf((*MockAHClient)(nil).PostGrantsRevokeWithResponse), varargs...)
}

// PostGrantsWithBodyWithResponse mocks base method.
func (m *MockAHClient) PostGrantsWithBodyWithResponse(arg0 context.Context, arg1 string, arg2 io.Reader, arg3 ...types.RequestEditorFn) (*types.PostGrantsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PostGrantsWithBodyWithResponse", varargs...)
	ret0, _ := ret[0].(*types.PostGrantsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostGrantsWithBodyWithResponse indicates an expected call of PostGrantsWithBodyWithResponse.
func (mr *MockAHClientMockRecorder) PostGrantsWithBodyWithResponse(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostGrantsWithBodyWithResponse", reflect.TypeOf((*MockAHClient)(nil).PostGrantsWithBodyWithResponse), varargs...)
}

// PostGrantsWithResponse mocks base method.
func (m *MockAHClient) PostGrantsWithResponse(arg0 context.Context, arg1 types.CreateGrant, arg2 ...types.RequestEditorFn) (*types.PostGrantsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PostGrantsWithResponse", varargs...)
	ret0, _ := ret[0].(*types.PostGrantsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostGrantsWithResponse indicates an expected call of PostGrantsWithResponse.
func (mr *MockAHClientMockRecorder) PostGrantsWithResponse(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostGrantsWithResponse", reflect.TypeOf((*MockAHClient)(nil).PostGrantsWithResponse), varargs...)
}

// RefreshAccessProvidersWithResponse mocks base method.
func (m *MockAHClient) RefreshAccessProvidersWithResponse(arg0 context.Context, arg1 ...types.RequestEditorFn) (*types.RefreshAccessProvidersResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RefreshAccessProvidersWithResponse", varargs...)
	ret0, _ := ret[0].(*types.RefreshAccessProvidersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshAccessProvidersWithResponse indicates an expected call of RefreshAccessProvidersWithResponse.
func (mr *MockAHClientMockRecorder) RefreshAccessProvidersWithResponse(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshAccessProvidersWithResponse", reflect.TypeOf((*MockAHClient)(nil).RefreshAccessProvidersWithResponse), varargs...)
}

// ValidateGrantWithBodyWithResponse mocks base method.
func (m *MockAHClient) ValidateGrantWithBodyWithResponse(arg0 context.Context, arg1 string, arg2 io.Reader, arg3 ...types.RequestEditorFn) (*types.ValidateGrantResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateGrantWithBodyWithResponse", varargs...)
	ret0, _ := ret[0].(*types.ValidateGrantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateGrantWithBodyWithResponse indicates an expected call of ValidateGrantWithBodyWithResponse.
func (mr *MockAHClientMockRecorder) ValidateGrantWithBodyWithResponse(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateGrantWithBodyWithResponse", reflect.TypeOf((*MockAHClient)(nil).ValidateGrantWithBodyWithResponse), varargs...)
}

// ValidateGrantWithResponse mocks base method.
func (m *MockAHClient) ValidateGrantWithResponse(arg0 context.Context, arg1 types.CreateGrant, arg2 ...types.RequestEditorFn) (*types.ValidateGrantResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateGrantWithResponse", varargs...)
	ret0, _ := ret[0].(*types.ValidateGrantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateGrantWithResponse indicates an expected call of ValidateGrantWithResponse.
func (mr *MockAHClientMockRecorder) ValidateGrantWithResponse(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateGrantWithResponse", reflect.TypeOf((*MockAHClient)(nil).ValidateGrantWithResponse), varargs...)
}

// ValidateSetupWithBodyWithResponse mocks base method.
func (m *MockAHClient) ValidateSetupWithBodyWithResponse(arg0 context.Context, arg1 string, arg2 io.Reader, arg3 ...types.RequestEditorFn) (*types.ValidateSetupResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateSetupWithBodyWithResponse", varargs...)
	ret0, _ := ret[0].(*types.ValidateSetupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateSetupWithBodyWithResponse indicates an expected call of ValidateSetupWithBodyWithResponse.
func (mr *MockAHClientMockRecorder) ValidateSetupWithBodyWithResponse(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSetupWithBodyWithResponse", reflect.TypeOf((*MockAHClient)(nil).ValidateSetupWithBodyWithResponse), varargs...)
}

// ValidateSetupWithResponse mocks base method.
func (m *MockAHClient) ValidateSetupWithResponse(arg0 context.Context, arg1 types.ValidateSetupJSONRequestBody, arg2 ...types.RequestEditorFn) (*types.ValidateSetupResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateSetupWithResponse", varargs...)
	ret0, _ := ret[0].(*types.ValidateSetupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateSetupWithResponse indicates an expected call of ValidateSetupWithResponse.
func (mr *MockAHClientMockRecorder) ValidateSetupWithResponse(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSetupWithResponse", reflect.TypeOf((*MockAHClient)(nil).ValidateSetupWithResponse), varargs...)
}
