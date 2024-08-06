package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct {
	Endpoint  string    `json:"endpoint"`
	LastCheck time.Time `json:"lastCheck"`
	Interval  int       `json:"interval"`
}

type Service struct {
	Name        string      `json:"name"`
	Host        string      `json:"host"`
	Time        time.Time   `json:"time"`
	HealthCheck HealthCheck `json:"healthCheck"`
}

var ServiceList = []Service{}

func RegisterServiceRouter(r *gin.Engine) {
	r.POST("/crate-hq-api/service", func(c *gin.Context) {
		var body Service
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(err)
			return
		}
		flag := true
		for _, service := range ServiceList {
			if service.Name == body.Name && service.Host == body.Host {
				flag = false
				break
			}
		}
		if !flag {
			c.Status(409)
			return
		}
		t := time.Now()
		body.Time = t
		body.HealthCheck.LastCheck = t
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
