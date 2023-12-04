package handlers

import (
	"github.com/soltanat/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	storage storage.Storage
}

func New(s storage.Storage) *Handlers {
	return &Handlers{storage: s}
}

func (h *Handlers) HandleMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	if len(path) != 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType, metricName, metricValue := path[2], path[3], path[4]
	if metricType == storage.Gauge || metricType == storage.Counter {
		if metricName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch metricType {
		case storage.Counter:
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = h.storage.StoreCounter(metricName, value)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case storage.Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = h.storage.StoreGauge(metricName, value)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

}
