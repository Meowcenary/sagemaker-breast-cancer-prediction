package main

import (
    // "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestApiStatusHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/status", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(apiStatusHandler)

    handler.ServeHTTP(rr, req)

    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check the response body
    var response map[string]string
    if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
        t.Fatal(err)
    }

    expected := map[string]string{"code": "200"}
    if !equal(response, expected) {
        t.Errorf("handler returned unexpected body: got %v want %v", response, expected)
    }
}

// Helper function to compare two maps
func equal(a, b map[string]string) bool {
    if len(a) != len(b) {
        return false
    }
    for k, v := range a {
        if b[k] != v {
						return false
				}
    }
		return true
}
