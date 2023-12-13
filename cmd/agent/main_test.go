package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/soltanat/metrics/internal/model"
)

type PollerMock struct {
	mock.Mock
}

func (m *PollerMock) Get() ([]model.Metric, error) {
	return []model.Metric{}, nil
}

func (m *PollerMock) Poll() error {
	args := m.Called()
	return args.Error(0)
}

type ReporterMock struct {
	mock.Mock
}

func (m *ReporterMock) Report() error {
	args := m.Called()
	return args.Error(0)
}

func TestRun(t *testing.T) {
	pollerMock := new(PollerMock)
	pollerMock.On("Poll").Return(nil).Times(1)

	reporterMock := new(ReporterMock)
	reporterMock.On("Report").Return(nil).Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	Run(ctx, time.Second-time.Millisecond*100, time.Second-time.Millisecond*100, pollerMock, reporterMock)

	assert.True(t, pollerMock.AssertExpectations(t))
	assert.True(t, reporterMock.AssertExpectations(t))
}
