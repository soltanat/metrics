package main

import (
	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/storage"
	"net/http"
)

func main() {
	s := storage.NewMemStorage()
	h := handler.New(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", h.HandleMetric)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
