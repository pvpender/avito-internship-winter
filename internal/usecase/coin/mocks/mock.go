// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen --source=deps.go --destination=mocks/mock.go
//

// Package mock_coin is a generated GoMock package.
package mock_coin

import (
	context "context"
	reflect "reflect"

	models "github.com/pvpender/avito-shop/internal/models"
	coin "github.com/pvpender/avito-shop/internal/usecase/coin"
	gomock "go.uber.org/mock/gomock"
)

// MockCoinRepository is a mock of CoinRepository interface.
type MockCoinRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCoinRepositoryMockRecorder
	isgomock struct{}
}

// MockCoinRepositoryMockRecorder is the mock recorder for MockCoinRepository.
type MockCoinRepositoryMockRecorder struct {
	mock *MockCoinRepository
}

// NewMockCoinRepository creates a new mock instance.
func NewMockCoinRepository(ctrl *gomock.Controller) *MockCoinRepository {
	mock := &MockCoinRepository{ctrl: ctrl}
	mock.recorder = &MockCoinRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoinRepository) EXPECT() *MockCoinRepositoryMockRecorder {
	return m.recorder
}

// CreateTransmission mocks base method.
func (m *MockCoinRepository) CreateTransmission(ctx context.Context, request *models.CoinOperationWithIds) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransmission", ctx, request)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransmission indicates an expected call of CreateTransmission.
func (mr *MockCoinRepositoryMockRecorder) CreateTransmission(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransmission", reflect.TypeOf((*MockCoinRepository)(nil).CreateTransmission), ctx, request)
}

// GetUserTransmissions mocks base method.
func (m *MockCoinRepository) GetUserTransmissions(ctx context.Context, userId uint32, transmissionType coin.TransmissionType) ([]*models.CoinOperationWithUsernames, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserTransmissions", ctx, userId, transmissionType)
	ret0, _ := ret[0].([]*models.CoinOperationWithUsernames)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserTransmissions indicates an expected call of GetUserTransmissions.
func (mr *MockCoinRepositoryMockRecorder) GetUserTransmissions(ctx, userId, transmissionType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserTransmissions", reflect.TypeOf((*MockCoinRepository)(nil).GetUserTransmissions), ctx, userId, transmissionType)
}

// MockCoinUseCase is a mock of CoinUseCase interface.
type MockCoinUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCoinUseCaseMockRecorder
	isgomock struct{}
}

// MockCoinUseCaseMockRecorder is the mock recorder for MockCoinUseCase.
type MockCoinUseCaseMockRecorder struct {
	mock *MockCoinUseCase
}

// NewMockCoinUseCase creates a new mock instance.
func NewMockCoinUseCase(ctrl *gomock.Controller) *MockCoinUseCase {
	mock := &MockCoinUseCase{ctrl: ctrl}
	mock.recorder = &MockCoinUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoinUseCase) EXPECT() *MockCoinUseCaseMockRecorder {
	return m.recorder
}

// SendCoin mocks base method.
func (m *MockCoinUseCase) SendCoin(ctx context.Context, userId uint32, request *models.SendCoinRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoin", ctx, userId, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoin indicates an expected call of SendCoin.
func (mr *MockCoinUseCaseMockRecorder) SendCoin(ctx, userId, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoin", reflect.TypeOf((*MockCoinUseCase)(nil).SendCoin), ctx, userId, request)
}
