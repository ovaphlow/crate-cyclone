package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func RegisterServiceRouter(r *gin.Engine) {
	r.POST("/crate-hq-api/service", func(c *gin.Context) {
		var body Service
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(err)
			return
		}
		flag := true
		for _, service := range ServiceList {
			if service.Name == body.Name && service.Host == c.ClientIP() && service.Port == body.Port {
				flag = false
				break
			}
		}
		if !flag {
			c.Status(409)
			return
		}
		protocol := "http"
		if c.Request.TLS != nil {
			protocol = "https"
		}
		body.Protocol = protocol
		body.Host = c.ClientIP()
		t := time.Now()
		body.Time = t
		body.HealthCheck.LastCheck = t
		body.HealthCheck.Failed = 0
		ServiceList = append(ServiceList, body)
		c.Status(201)
	})

	r.DELETE("/crate-hq-api/service", func(c *gin.Context) {
		name := c.Query("name")
		host := c.Query("host")
		for i, service := range ServiceList {
			if service.Name == name && service.Host == host {
				ServiceList = append(ServiceList[:i], ServiceList[i+1:]...)
			}
		}
	})

	r.GET("/crate-hq-api/service", func(c *gin.Context) {
		c.JSON(200, ServiceList)
	})
}
