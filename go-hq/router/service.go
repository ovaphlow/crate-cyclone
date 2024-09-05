package router

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body Service
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 处理 body
	})
}
