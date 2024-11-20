package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type NetworkStats struct {
    Interface string `json:"interface"`
    RXBytes   uint64 `json:"rx_bytes"`
    TXBytes   uint64 `json:"tx_bytes"`
}

func readNetworkStats() ([]NetworkStats, error) {
    file, err := os.Open("/proc/net/dev")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var stats []NetworkStats

    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, ":") { // Valid lines include a colon
            fields := strings.Fields(strings.TrimSpace(line))
            iface := strings.TrimSuffix(fields[0], ":")
            rxBytes, _ := strconv.ParseUint(fields[1], 10, 64)
            txBytes, _ := strconv.ParseUint(fields[9], 10, 64)
            stats = append(stats, NetworkStats{
                Interface: iface,
                RXBytes:   rxBytes,
                TXBytes:   txBytes,
            })
        }
    }
    return stats, scanner.Err()
}

func getBandwidthHandler(w http.ResponseWriter, r *http.Request) {
    logRequestDetails(r)

    stats, err := readNetworkStats()
    if err != nil {
        http.Error(w, "Error reading network statistics: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats) 
}

func getCpuHandler(w http.ResponseWriter, r *http.Request) {
    logRequestDetails(r)
    out, err := exec.Command("/bin/sh", "-c", "top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\\([0-9.]*\\)%* id.*/\\1/' | awk '{print 100 - $1\"%\"}'").Output()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"CPU Usage": string(out)})
}

func getMemoryHandler(w http.ResponseWriter, r *http.Request) {
    logRequestDetails(r)
    out, err := exec.Command("/bin/sh", "-c", "free -m | grep Mem | awk '{print $3\" MB used out of \"$2\" MB\"}'").Output()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"Memory Usage": string(out)})
}

func getDiskHandler(w http.ResponseWriter, r *http.Request) {
    logRequestDetails(r)
    out, err := exec.Command("df", "-h").Output()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"Disk Usage": string(out)})
}

func crashHandler(w http.ResponseWriter, r *http.Request) {
    
    panic("Intentional crash triggered")
}

func logRequestDetails(r *http.Request) {
   
    log.Printf("Received %s request from %s", r.Method, r.RemoteAddr)
}

func main() {
    http.HandleFunc("/cpu", getCpuHandler)
    http.HandleFunc("/memory", getMemoryHandler)
    http.HandleFunc("/disk", getDiskHandler)
    http.HandleFunc("/bandwidth", getBandwidthHandler)
    http.HandleFunc("/crash", crashHandler)
    port := "8080"
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
