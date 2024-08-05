package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"ovaphlow/crate/hq/infrastructure"
	"ovaphlow/crate/hq/router"
	"ovaphlow/crate/hq/subscriber"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infrastructure.InitSlog()

	database_type := os.Getenv("DATABASE_TYPE")
	if database_type == "postgres" {
		infrastructure.InitPostgres()
	} else if database_type == "mysql" {
		//
	} else {
		log.Panic("未设置数据库")
	}
}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("access-control-allow-origin", "*")
		c.Header("access-control-allow-credentials", "true")
		c.Header("access-control-allow-methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("access-control-allow-headers", "accept, content-type, content-length, accept-encoding, x-csrf-token, authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Use(func(c *gin.Context) {
		c.Header("x-frame-options", "DENY")
		c.Header("x-content-type-options", "nosniff")
		c.Header("x-xss-protection", "1; mode=block")
		c.Next()
	})

	r.Use(func(c *gin.Context) {
		c.Header("x-api-version", "2024-01-06")
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Use(func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			err := ctx.Errors[0].Err
			ctx.JSON(500, gin.H{
				"type":     "about:blank",
				"status":   500,
				"title":    "服务器错误",
				"detail":   err.Error(),
				"instance": ctx.Request.URL.Path,
			})
			ctx.Abort()
		}
	})

	r.Static("/html", "./html")

	schemaRepo := infrastructure.NewSchemaRepoImpl(infrastructure.Postgres)
	schemaService := infrastructure.NewSchemaService(schemaRepo)

	subscriberRepo := subscriber.NewSubscriberRepoImpl(infrastructure.Postgres)
	subscriberService := subscriber.NewSubscriberService(subscriberRepo, schemaService)
	router.RegisterSubscriberRouter(r, subscriberService)

	determinTarget := func(c *gin.Context) string {
		if c.Request.URL.Path == "service1" {
			return "http://service1.example.com"
		} else if c.Request.URL.Path == "service2" {
			return "http://service2.example.com"
		}
		return "http://default.example.com"
	}

	dynamicProxy := func(c *gin.Context) {
		target := determinTarget(c)
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		c.Request.URL.Host = remote.Host
		c.Request.URL.Scheme = remote.Scheme
		c.Request.Header.Set("x-forwarded-host", c.Request.Header.Get("Host"))
		c.Request.Host = remote.Host
		c.Request.Header.Set("x-auth", "1123")
		proxy.ServeHTTP(c.Writer, c.Request)
	}

	r.Any("/proxy/*proxyPath", dynamicProxy)

	r.Run("0.0.0.0:" + os.Getenv("PORT"))
}
