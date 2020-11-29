package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/opensatelliteproject/SatHelperApp/DSP"
	"net/http"
	"time"
)

func getFFT(w http.ResponseWriter, r *http.Request) {
	data := DSP.GetFFTImage()
	if data != nil {
		w.Header().Add("Content-Type", "image/jpeg")
		w.WriteHeader(200)
		_, _ = w.Write(data)
		return
	}

	w.WriteHeader(500)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	stats := DSP.GetStats()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	jsonData, _ := json.MarshalIndent(stats, "", "  ")
	_, _ = w.Write(jsonData)
}

func startHTTP() error {
	r := mux.NewRouter()
	r.HandleFunc("/fft", getFFT)
	r.HandleFunc("/stats", getStats)
	srv := &http.Server{
		Handler: r,
		Addr:    ":14123",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}
