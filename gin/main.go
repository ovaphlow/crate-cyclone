package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"ovaphlow.crate/cyclone/infrastruction"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infrastruction.InitSlog()
	infrastruction.InitPostgres()
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

	schemaRepo := NewSchemaRepoImpl(infrastruction.Postgres)
	schemaService := NewSchemaService(schemaRepo)

	AddSchemaEndpoints(r, schemaService)

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

	r.Run("0.0.0.0:" + os.Getenv("PORT"))
}
