// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -source=handler.go -destination=handler_mocks_test.go -package=api_test
//

// Package api_test is a generated GoMock package.
package api_test

import (
	context "context"
	reflect "reflect"

	token "github.com/aspirin100/JWT-API/internal/token"
	gomock "go.uber.org/mock/gomock"
)

// MockTokenService is a mock of TokenService interface.
type MockTokenService struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceMockRecorder
	isgomock struct{}
}

// MockTokenServiceMockRecorder is the mock recorder for MockTokenService.
type MockTokenServiceMockRecorder struct {
	mock *MockTokenService
}

// NewMockTokenService creates a new mock instance.
func NewMockTokenService(ctrl *gomock.Controller) *MockTokenService {
	mock := &MockTokenService{ctrl: ctrl}
	mock.recorder = &MockTokenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenService) EXPECT() *MockTokenServiceMockRecorder {
	return m.recorder
}

// CreateNewTokensPair mocks base method.
func (m *MockTokenService) CreateNewTokensPair(ctx context.Context, params *token.PairParams) (*string, *string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewTokensPair", ctx, params)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(*string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateNewTokensPair indicates an expected call of CreateNewTokensPair.
func (mr *MockTokenServiceMockRecorder) CreateNewTokensPair(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewTokensPair", reflect.TypeOf((*MockTokenService)(nil).CreateNewTokensPair), ctx, params)
}

// RefreshTokenPair mocks base method.
func (m *MockTokenService) RefreshTokenPair(ctx context.Context, params *token.RefreshTokensParams) (*string, *string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTokenPair", ctx, params)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(*string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RefreshTokenPair indicates an expected call of RefreshTokenPair.
func (mr *MockTokenServiceMockRecorder) RefreshTokenPair(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokenPair", reflect.TypeOf((*MockTokenService)(nil).RefreshTokenPair), ctx, params)
}
