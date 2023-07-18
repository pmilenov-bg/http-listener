package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function with the request and response recorder
	handleRequest(rr, req)

	// Verify the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Verify the response body
	expectedResponse := "GET request received"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, but got %q", expectedResponse, rr.Body.String())
	}
}

func TestHandlePostRequest(t *testing.T) {
	// Create a JSON request body
	data := map[string]interface{}{
		"message": "Hello, World!",
	}
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create a new HTTP request with the JSON body
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function with the request and response recorder
	handleRequest(rr, req)

	// Verify the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Verify the response body
	expectedResponse := "POST request received"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q, but got %q", expectedResponse, rr.Body.String())
	}
}

// Add similar tests for other HTTP methods (PUT, DELETE) if needed

func TestCalculateBufferSize(t *testing.T) {
	// Initialize the buffer with some data
	buffer = []interface{}{
		map[string]interface{}{
			"field1": "value1",
		},
		map[string]interface{}{
			"field2": "value2",
		},
	}

	// Call the calculateBufferSize function
	bufferSize := calculateBufferSize()

	// Verify the buffer size
	expectedSize := len(`[{"field1":"value1"},{"field2":"value2"}]`)
	if bufferSize != expectedSize {
		t.Errorf("Expected buffer size %d, but got %d", expectedSize, bufferSize)
	}
}

// Add additional tests for other functions as needed

func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()

	// Clean up any test-specific resources here

	// Exit the tests
	os.Exit(code)
}
