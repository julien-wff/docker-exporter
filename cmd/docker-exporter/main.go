package main

import (
	"github.com/julien-wff/docker-exporter/internal/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	reg := prometheus.NewRegistry()

	dockerCollector := collectors.NewDockerCollector()
	reg.MustRegister(dockerCollector)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	log.Printf("Starting server on port 9100")
	log.Fatal(http.ListenAndServe(":9100", mux))
}
