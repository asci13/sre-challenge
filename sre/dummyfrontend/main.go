package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-playground/validator/v10"
	"github.com/heptiolabs/healthcheck"

	"github.com/gorilla/mux"
)

var settings *ServerSettings
var validate *validator.Validate

type Request struct {
	Port int `validate:"required"`
}

type ServerSettings struct {
	ListenPort     string
	BackendAddress string
}

func replyWithError(w http.ResponseWriter, statusCode int, statusMessage string) {
	log.Printf(statusMessage)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = statusMessage
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func respondToDefaultRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Default request received.\n")
	requestId := validatePortRequest(r.URL.Path[1:])
	log.Printf("RequestID is %d\n", requestId)
	if requestId < 0 {
		replyWithError(w, http.StatusBadRequest, "RequestID wasn't a valid int.")
		return
	}
	log.Printf("Requesting file from %s\n", settings.BackendAddress)
	requestFile, fileError := requestFileFromBackend(w, requestId, settings.BackendAddress)
	if fileError != nil {
		var errorMessage string
		errorMessage = "An unknown error occurred. " + fileError.Error()
		replyWithError(w, http.StatusBadRequest, errorMessage)
		return
	}
	defer os.Remove(requestFile)

	log.Printf("Received file from %s at %s\n", settings.BackendAddress, requestFile)

	http.ServeFile(w, r, requestFile)
}

func requestFileFromBackend(w http.ResponseWriter, requestId int, requestUri string) (string, error) {
	log.Printf("GET %s", requestUri)
	resp, err := http.Get(requestUri)
	if err != nil {
		replyWithError(w, http.StatusFailedDependency, "Fileserver returned an error.")
		return "", errors.New("Fileserver returned an error.")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		replyWithError(w, http.StatusFailedDependency, "File from Fileserver couldn't be read")
		return "", errors.New("File from Fileserver couldn't be read.")
	}
	mimeType := http.DetectContentType(bytes)
	fileExtension := "pdf"
	if mimeType == "image/png" {
		fileExtension = "png"
	}
	requestFile := fmt.Sprintf("%d.%s", requestId, fileExtension)
	fileWriteError := ioutil.WriteFile(requestFile, bytes, 0644)
	if fileWriteError != nil {
		replyWithError(w, http.StatusFailedDependency, "File from Fileserver couldn't be processed.")
		return "", errors.New("File from Fileserver couldn't be processed.")
	}
	return requestFile, nil
}

func validatePortRequest(id string) int {
	validate = validator.New()
	i, convErr := strconv.Atoi(id)
	if convErr != nil {
		return -1
	}
	requestPort := &Request{
		Port: i,
	}

	err := validate.Struct(requestPort)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
			return -1
		}
	}
	return requestPort.Port
}

func AddServiceRoutes(rtr *mux.Router) {

}

func main() {
	listenPort := os.Getenv("ListenPort")
	backendAddress := os.Getenv("BackendAddress")
	log.Printf("ListenPort is %s", listenPort)
	log.Printf("BackendAddress is %s", backendAddress)
	settings = &ServerSettings{
		ListenPort:     listenPort,
		BackendAddress: backendAddress,
	}

	if len(settings.ListenPort) == 0 {
		settings.ListenPort = ":5555"

	}
	if len(settings.BackendAddress) == 0 {
		log.Fatalf("BackendAddress is empty - please set environment variable.")
	}

	adminMux := http.NewServeMux()
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

	adminMux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// Expose a liveness check on /live
	adminMux.HandleFunc("/live", health.LiveEndpoint)

	// Expose a readiness check on /ready
	adminMux.HandleFunc("/ready", health.ReadyEndpoint)
	adminMux.HandleFunc("/", respondToDefaultRequest)

	log.Printf("Handlers started.\n")
	err := http.ListenAndServe(settings.ListenPort, adminMux)
	if err != nil {
		log.Fatal(err)
	}
}
