package storage

import (
	"github.com/soltanat/metrics/internal"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Store(metric *internal.Metric) error {
	args := m.Called(metric)
	return args.Error(0)
}

func (m *MockStorage) GetGauge(name string) (*internal.Metric, error) {
	args := m.Called(name)

	var r0 *internal.Metric
	if args.Get(0) == nil {
		r0 = nil
	} else {
		r0 = args.Get(0).(*internal.Metric)
	}
	return r0, args.Error(1)
}

func (m *MockStorage) GetCounter(name string) (*internal.Metric, error) {
	args := m.Called(name)
	var r0 *internal.Metric
	if args.Get(0) == nil {
		r0 = nil
	} else {
		r0 = args.Get(0).(*internal.Metric)
	}
	return r0, args.Error(1)
}

func (m *MockStorage) GetList() ([]internal.Metric, error) {
	args := m.Called()
	return args.Get(0).([]internal.Metric), args.Error(1)
}
