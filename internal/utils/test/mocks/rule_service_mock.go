// Code generated by MockGen. DO NOT EDIT.
// Source: ./rule_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/nicopozo/mockserver/internal/model"
	reflect "reflect"
)

// MockIRuleService is a mock of IRuleService interface
type MockIRuleService struct {
	ctrl     *gomock.Controller
	recorder *MockIRuleServiceMockRecorder
}

// MockIRuleServiceMockRecorder is the mock recorder for MockIRuleService
type MockIRuleServiceMockRecorder struct {
	mock *MockIRuleService
}

// NewMockIRuleService creates a new mock instance
func NewMockIRuleService(ctrl *gomock.Controller) *MockIRuleService {
	mock := &MockIRuleService{ctrl: ctrl}
	mock.recorder = &MockIRuleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIRuleService) EXPECT() *MockIRuleServiceMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockIRuleService) Save(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, rule)
	ret0, _ := ret[0].(*model.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save
func (mr *MockIRuleServiceMockRecorder) Save(ctx, rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIRuleService)(nil).Save), ctx, rule)
}

// Update mocks base method
func (m *MockIRuleService) Update(ctx context.Context, key string, rule *model.Rule) (*model.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, key, rule)
	ret0, _ := ret[0].(*model.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockIRuleServiceMockRecorder) Update(ctx, key, rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIRuleService)(nil).Update), ctx, key, rule)
}

// Get mocks base method
func (m *MockIRuleService) Get(ctx context.Context, key string) (*model.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(*model.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockIRuleServiceMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIRuleService)(nil).Get), ctx, key)
}

// Search mocks base method
func (m *MockIRuleService) Search(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, params, paging)
	ret0, _ := ret[0].(*model.RuleList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockIRuleServiceMockRecorder) Search(ctx, params, paging interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIRuleService)(nil).Search), ctx, params, paging)
}

// SearchByMethodAndPath mocks base method
func (m *MockIRuleService) SearchByMethodAndPath(ctx context.Context, method, path string) (*model.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByMethodAndPath", ctx, method, path)
	ret0, _ := ret[0].(*model.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchByMethodAndPath indicates an expected call of SearchByMethodAndPath
func (mr *MockIRuleServiceMockRecorder) SearchByMethodAndPath(ctx, method, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByMethodAndPath", reflect.TypeOf((*MockIRuleService)(nil).SearchByMethodAndPath), ctx, method, path)
}

// Delete mocks base method
func (m *MockIRuleService) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockIRuleServiceMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIRuleService)(nil).Delete), ctx, key)
}
