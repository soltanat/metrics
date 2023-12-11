package storage

import (
	"github.com/soltanat/metrics/internal"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMemStorage_StoreCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"store counter",
			fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args{
				name:  "test name",
				value: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			m := internal.Metric{
				Type:    internal.CounterType,
				Name:    tt.args.name,
				Counter: tt.args.value,
			}
			if err := s.Store(&m); (err != nil) != tt.wantErr {
				assert.Equal(t, tt.fields.gauge[tt.args.name], tt.args.value)
				assert.NoError(t, err)
			}
		})
	}
}

func TestMemStorage_StoreGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"store gauge",
			fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args{
				name:  "test name",
				value: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			m := internal.Metric{
				Type:  internal.GaugeType,
				Name:  tt.args.name,
				Gauge: tt.args.value,
			}
			if err := s.Store(&m); (err != nil) != tt.wantErr {
				assert.Equal(t, tt.fields.counter[tt.args.name], tt.args.value)
				assert.NoError(t, err)
			}
		})
	}
}

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
