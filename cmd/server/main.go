package main

import (
	"github.com/soltanat/metrics/internal/handlers"
	"github.com/soltanat/metrics/internal/storage"
	"net/http"
)

func main() {
	s := storage.NewMemStorage()
	h := handlers.New(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", h.HandleMetric)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}

}
