package storage

import (
	"github.com/stretchr/testify/mock"

	"github.com/soltanat/metrics/internal/model"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Store(metric *model.Metric) error {
	args := m.Called(metric)
	return args.Error(0)
}

func (m *MockStorage) GetGauge(name string) (*model.Metric, error) {
	args := m.Called(name)

	var r0 *model.Metric
	if args.Get(0) == nil {
		r0 = nil
	} else {
		r0 = args.Get(0).(*model.Metric)
	}
	return r0, args.Error(1)
}

func (m *MockStorage) GetCounter(name string) (*model.Metric, error) {
	args := m.Called(name)
	var r0 *model.Metric
	if args.Get(0) == nil {
		r0 = nil
	} else {
		r0 = args.Get(0).(*model.Metric)
	}
	return r0, args.Error(1)
}

func (m *MockStorage) GetList() ([]model.Metric, error) {
	args := m.Called()
	return args.Get(0).([]model.Metric), args.Error(1)
}
