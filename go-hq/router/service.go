package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HealthCheck struct {
	Endpoint    string    `json:"endpoint"`
	LastCheck   time.Time `json:"lastCheck"`
	LastSuccess time.Time `json:"lastSuccess"`
	Interval    int       `json:"interval"`
	Failed      int       `json:"failed"`
}

type Service struct {
	Name        string      `json:"name"`
	Protocol    string      `json:"protocol"`
	Host        string      `json:"host"`
	Port        int         `json:"port"`
	Time        time.Time   `json:"time"`
	HealthCheck HealthCheck `json:"healthCheck"`
}

var ServiceList = []Service{}

func PerformHealthCheck(sec int) {
	for i := range ServiceList {
		service := &ServiceList[i]
		service.HealthCheck.LastCheck = time.Now()
		service.HealthCheck.Interval = sec
		url := fmt.Sprintf("%s://%s:%d%s", service.Protocol, service.Host, service.Port, service.HealthCheck.Endpoint)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			service.HealthCheck.Failed++
			continue
		}
		defer resp.Body.Close()
		service.HealthCheck.Failed = 0
		service.HealthCheck.LastSuccess = time.Now()
	}
}

func RegisterServiceRouter(mux *http.ServeMux) {
	mux.HandleFunc("/crate-hq-api/service", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(ServiceList); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			var body Service
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			flag := true
			clientIP := strings.Split(r.RemoteAddr, ":")[0]
			for _, service := range ServiceList {
				if service.Name == body.Name && service.Host == clientIP && service.Port == body.Port {
					flag = false
					break
				}
			}
			if !flag {
				http.Error(w, "Service already exists", http.StatusConflict)
				return
			}
			protocol := "http"
			if r.TLS != nil {
				protocol = "https"
			}
			body.Protocol = protocol
			body.Host = clientIP
			t := time.Now()
			body.Time = t
			body.HealthCheck.LastCheck = t
			body.HealthCheck.Failed = 0
			ServiceList = append(ServiceList, body)
			w.WriteHeader(http.StatusCreated)
		} else if r.Method == http.MethodDelete {
			name := r.URL.Query().Get("name")
			host := r.URL.Query().Get("host")
			for i, service := range ServiceList {
				if service.Name == name && service.Host == host {
					ServiceList = append(ServiceList[:i], ServiceList[i+1:]...)
				}
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
