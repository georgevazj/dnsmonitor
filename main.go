package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/shirou/gopsutil/mem"
)

// DNSAlive command
const DNSAlive = "dig @%s %s +short"

// Server JSON model
type Server struct {
	Hostname string  `json:"hostname,omitempty"`
	FreeMem  uint64  `json:"freemem,omitempty"`
	Average  float64 `json:"cpuUsage,omitempty"`
	DNSAlive string  `json:"dnsalive,omitempty"`
}

var servers []Server

// execCmd executes a shell command and returns its output
func execCmd(cmd string) string {
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	output := strings.TrimSuffix(string(out), "\n")
	return output
}

// getStatus returns DNS server stats
func getStatus(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	memory, _ := mem.VirtualMemory()
	avg, _ := loadavg.Get()
	dnsAlive := execCmd(DNSAlive)
	json.NewEncoder(w).Encode(Server{Hostname: hostname, FreeMem: memory.Free, Average: avg.Loadavg1, DNSAlive: dnsAlive})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/status", getStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
