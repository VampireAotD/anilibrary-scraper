// Code generated by MockGen. DO NOT EDIT.
// Source: controller.go
//
// Generated by this command:
//
//	mockgen -source=controller.go -destination=./mocks.go -package=anime
//

// Package anime is a generated GoMock package.
package anime

import (
	context "context"
	reflect "reflect"

	scraper "github.com/VampireAotD/anilibrary-scraper/internal/application/usecase/scraper"
	entity "github.com/VampireAotD/anilibrary-scraper/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockScraperUseCase is a mock of ScraperUseCase interface.
type MockScraperUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockScraperUseCaseMockRecorder
	isgomock struct{}
}

// MockScraperUseCaseMockRecorder is the mock recorder for MockScraperUseCase.
type MockScraperUseCaseMockRecorder struct {
	mock *MockScraperUseCase
}

// NewMockScraperUseCase creates a new mock instance.
func NewMockScraperUseCase(ctrl *gomock.Controller) *MockScraperUseCase {
	mock := &MockScraperUseCase{ctrl: ctrl}
	mock.recorder = &MockScraperUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScraperUseCase) EXPECT() *MockScraperUseCaseMockRecorder {
	return m.recorder
}

// Scrape mocks base method.
func (m *MockScraperUseCase) Scrape(ctx context.Context, dto scraper.DTO) (entity.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scrape", ctx, dto)
	ret0, _ := ret[0].(entity.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scrape indicates an expected call of Scrape.
func (mr *MockScraperUseCaseMockRecorder) Scrape(ctx, dto any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scrape", reflect.TypeOf((*MockScraperUseCase)(nil).Scrape), ctx, dto)
}
