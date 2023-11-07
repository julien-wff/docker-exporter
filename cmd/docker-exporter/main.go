package main

import (
	"github.com/julien-wff/docker-exporter/internal/collectors"
	"github.com/julien-wff/docker-exporter/internal/config"
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

	cfg := config.GetConfig()
	log.Printf("----- Docker Exporter -----\n")
	log.Println("")
	log.Println("Config:")
	log.Printf("- Using a request timeout of %s\n", cfg.RequestTimeout.String())
	log.Println("")

	log.Printf("Starting server on port 9100\n")
	log.Fatal(http.ListenAndServe(":9100", mux))
}
