// main_test.go
package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
)

func init() {
    // Disable logging during tests
    log.SetOutput(io.Discard)
}

func TestGetCpuHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/cpu", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getCpuHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedContentType := "application/json"
    if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
        t.Errorf("content type header does not match: got %v want %v", ct, expectedContentType)
    }

    var response map[string]string
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("could not parse JSON response: %v", err)
    }

    if _, ok := response["CPU Usage"]; !ok {
        t.Errorf("response JSON does not contain 'CPU Usage' key")
    }
}

func TestGetMemoryHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/memory", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getMemoryHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedContentType := "application/json"
    if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
        t.Errorf("content type header does not match: got %v want %v", ct, expectedContentType)
    }

    var response map[string]string
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("could not parse JSON response: %v", err)
    }

    if _, ok := response["Memory Usage"]; !ok {
        t.Errorf("response JSON does not contain 'Memory Usage' key")
    }
}

func TestGetDiskHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/disk", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getDiskHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedContentType := "application/json"
    if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
        t.Errorf("content type header does not match: got %v want %v", ct, expectedContentType)
    }

    var response map[string]string
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("could not parse JSON response: %v", err)
    }

    if _, ok := response["Disk Usage"]; !ok {
        t.Errorf("response JSON does not contain 'Disk Usage' key")
    }
}

func TestGetBandwidthHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/bandwidth", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getBandwidthHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedContentType := "application/json"
    if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
        t.Errorf("content type header does not match: got %v want %v", ct, expectedContentType)
    }

    var stats []NetworkStats
    err = json.Unmarshal(rr.Body.Bytes(), &stats)
    if err != nil {
        t.Errorf("could not parse JSON response: %v", err)
    }

    if len(stats) == 0 {
        t.Errorf("expected at least one network interface in response")
    }
}

func TestCrashHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/crash", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(crashHandler)

    defer func() {
        if r := recover(); r == nil {
            t.Errorf("expected panic but did not get one")
        }
    }()

    handler.ServeHTTP(rr, req)
}

func TestReadNetworkStats(t *testing.T) {
    stats, err := readNetworkStats()
    if err != nil {
        t.Errorf("readNetworkStats returned error: %v", err)
    }

    if len(stats) == 0 {
        t.Errorf("expected at least one network interface")
    }

    for _, stat := range stats {
        if stat.Interface == "" {
            t.Errorf("found network interface with empty name")
        }
    }
}
