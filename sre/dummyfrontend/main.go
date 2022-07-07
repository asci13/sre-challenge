package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/heptiolabs/healthcheck"
)

func respondRandomly(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
}

func main() {
	listenPort := os.Getenv("ListenPort")
	backendAddress := os.Getenv("BackendAddress")

	if len(listenPort) == 0 {
		listenPort = ":5555"
	}

	// Create a new Prometheus registry (you'd likely already have one of these).
	registry := prometheus.NewRegistry()

	// Create a metrics-exposing Handler for the Prometheus registry
	// The healthcheck related metrics will be prefixed with the provided namespace
	health := healthcheck.NewMetricsHandler(registry, "dummyfrontend")

	// Add a simple readiness check that always fails.
	health.AddReadinessCheck("provider-service-available", healthcheck.HTTPGetCheck(backendAddress, 500*time.Millisecond))

	// Add a liveness check that always succeeds
	health.AddLivenessCheck("service-started", func() error {
		if len(backendAddress) == 0 {
			return fmt.Errorf("backendAddress environment not configured.")
		}
		return nil
	})

	adminMux := http.NewServeMux()

	// Expose prometheus metrics on /metrics
	adminMux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// Expose a liveness check on /live
	adminMux.HandleFunc("/live", health.LiveEndpoint)

	// Expose a readiness check on /ready
	adminMux.HandleFunc("/ready", health.ReadyEndpoint)
	http.HandleFunc("/", respondRandomly)

	fmt.Printf("Handlers started.\n")
	err := http.ListenAndServe(listenPort, adminMux)
	if err != nil {
		log.Fatal(err)
	}
}
