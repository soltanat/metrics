package storage

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/soltanat/metrics/internal/model"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			"success new",
			&MemStorage{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
				mu:      &sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
		mu      *sync.RWMutex
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get existed",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: map[string]int64{"metric-name": 1},
				mu:      &sync.RWMutex{},
			},
			args:    args{name: "metric-name"},
			want:    model.NewCounter("metric-name", 1),
			wantErr: assert.NoError,
		},
		{
			name: "get not existed",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
				mu:      &sync.RWMutex{},
			},
			args: args{name: "metric-name"},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, model.ErrMetricNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
				mu:      tt.fields.mu,
			}
			got, err := s.GetCounter(tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("GetCounter(%v)", tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetCounter(%v)", tt.args.name)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
		mu      *sync.RWMutex
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get existed",
			fields: fields{
				gauge:   map[string]float64{"metric-name": 1.1},
				counter: map[string]int64{},
				mu:      &sync.RWMutex{},
			},
			args:    args{name: "metric-name"},
			want:    model.NewGauge("metric-name", 1.1),
			wantErr: assert.NoError,
		},
		{
			name: "get not existed",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
				mu:      &sync.RWMutex{},
			},
			args: args{name: "metric-name"},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, model.ErrMetricNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
				mu:      tt.fields.mu,
			}
			got, err := s.GetGauge(tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("GetGauge(%v)", tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetGauge(%v)", tt.args.name)
		})
	}
}

func TestMemStorage_GetList(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
		mu      *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get list not empty",
			fields: fields{
				gauge:   map[string]float64{"gauge-metric": 1.1},
				counter: map[string]int64{"counter-metric": 1},
				mu:      &sync.RWMutex{},
			},
			want: []model.Metric{
				*model.NewCounter("counter-metric", 1),
				*model.NewGauge("gauge-metric", 1.1),
			},
			wantErr: assert.NoError,
		},
		{
			name: "get list empty",
			fields: fields{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
				mu:      &sync.RWMutex{},
			},
			want:    []model.Metric{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
				mu:      tt.fields.mu,
			}
			got, err := s.GetList()
			if !tt.wantErr(t, err, "GetList()") {
				return
			}
			assert.Equalf(t, tt.want, got, "GetList()")
		})
	}
}

func TestMemStorage_Store(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
		mu      *sync.RWMutex
	}
	type args struct {
		metric *model.Metric
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		assertStored func(t assert.TestingT, fields fields) bool
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "store counter",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
				mu:      &sync.RWMutex{},
			},
			args: args{
				metric: model.NewCounter("metric", 1),
			},
			assertStored: func(t assert.TestingT, fields fields) bool {
				return assert.NotEqual(t, fields.counter["metric"], 1)
			},
			wantErr: assert.NoError,
		},
		{
			name: "store gauge",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
				mu:      &sync.RWMutex{},
			},
			args: args{
				metric: model.NewGauge("metric", 1.1),
			},
			assertStored: func(t assert.TestingT, fields fields) bool {
				return assert.NotNil(t, fields.gauge["metric"], 1.1)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
				mu:      tt.fields.mu,
			}
			tt.wantErr(t, s.Store(tt.args.metric), fmt.Sprintf("Store(%v)", tt.args.metric))
			tt.assertStored(t, tt.fields)
		})
	}
}
