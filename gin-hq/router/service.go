package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct {
	Endpoint  string
	LastCheck time.Time
	Interval  int
}

type Service struct {
	Name        string
	URL         string
	Time        time.Time
	HealthCheck HealthCheck
}

var ServiceList = []Service{}

func RegisterServiceRouter(r *gin.Engine) {
	r.POST("/crate-hq-api/service", func(c *gin.Context) {
		var body Service
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(err)
			return
		}
		body.Time = time.Now()
		for _, service := range ServiceList {
			if service.Name != body.Name {
				ServiceList = append(ServiceList, body)
			}
		}
		c.Status(201)
	})

	r.DELETE("/crate-hq-api/service", func(c *gin.Context) {
		name := c.Query("name")
		url := c.Query("url")
		for i, service := range ServiceList {
			if service.Name == name && service.URL == url {
				ServiceList = append(ServiceList[:i], ServiceList[i+1:]...)
			}
		}
	})

	r.GET("/crate-hq-api/service", func(c *gin.Context) {
		c.JSON(200, ServiceList)
	})
}
