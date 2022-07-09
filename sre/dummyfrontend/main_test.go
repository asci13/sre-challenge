package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var portChecks = []struct {
	in       string
	expected int
}{
	{"abc", -1},
	{"1", 1},
}

// Determine request id validity correctly
func Test_validatePortRequest(t *testing.T) {
	for _, tt := range portChecks {
		t.Run(tt.in, func(t *testing.T) {

			//request, _ := http.NewRequest(http.MethodGet, tt.in, nil)

			response := validatePortRequest(tt.in)

			got := response
			if got != tt.expected {
				t.Errorf("got %d, want %d", got, tt.expected)
			}

		})
	}
}

// Return bad request in case of invalid request id
func Test_respondToDefaultRequest_IncorrectRequestId(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/abc", nil)
	response := httptest.NewRecorder()

	respondToDefaultRequest(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode, "Should return bad request.")
	assert.Equal(t, "{\"message\":\"RequestID wasn't a valid int.\"}", response.Body.String(), "Should provide the correct error output as json.")
}

// Return bad request in case of invalid request id
func Test_respondToDefaultRequest_CorrectRequestId(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://127.0.0.1:5555/123", nil)
	//response := httptest.NewRecorder()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	//respondToDefaultRequest(response, request)
	assert.NotNil(t, response)
	//assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode, "Should return bad request.")
	//assert.Equal(t, "{\"message\":\"RequestID wasn't a valid int.\"}", response.Body.String(), "Should provide the correct error output as json.")
}

var fileChecks = []struct {
	responseFile      string
	requestId         int
	expected_filename string
}{
	{"./testfiles/dummy.pdf", 123, "123.pdf"},
	{"./testfiles/dummy.png", 1255, "1255.png"},
}

func Test_respondToDefaultRequest_ReturnsFile(t *testing.T) {

	for _, tt := range fileChecks {
		t.Run(tt.responseFile, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				http.ServeFile(w, r, tt.responseFile)
			}))
			defer server.Close()
			response := httptest.NewRecorder()

			result, _ := requestFileFromBackend(response, tt.requestId, server.URL)
			defer os.Remove(result)

			assert.Equal(t, tt.expected_filename, result)

		})
	}
}
