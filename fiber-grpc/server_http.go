package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"ovaphlow/cratecyclone/configuration"
	"ovaphlow/cratecyclone/schema"
	"ovaphlow/cratecyclone/subscriber"
	"ovaphlow/cratecyclone/utility"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func registerService(app string, port int, url string) {
	body := map[string]interface{}{
		"name": app,
		"port": port,
		"healthCheck": map[string]interface{}{
			"endpoint": "/crate-api/health",
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Print(err.Error())
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		log.Print("注册服务失败")
		log.Print(err.Error())
		return
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("注册服务失败")
		log.Print(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Print("注册服务失败")
		log.Print(resp.Status)
	}
}

func HTTPServe(port string) {
	err := godotenv.Load()
	if err != nil {
		log.Println("加载环境变量失败")
		log.Println(err.Error())
	}
	a := os.Getenv("APP_NAME")
	h := os.Getenv("HQ_HOST")
	p := os.Getenv("HQ_PORT")
	e := os.Getenv("HQ_REGISTER_ENDPOINT")
	_port, err := strconv.Atoi(port)
	if err != nil {
		log.Print("端口异常 " + err.Error())
	}
	registerService(a, _port, h+":"+p+e)

	utility.InitPostgres()

	app := fiber.New(fiber.Config{
		// Prefork 模式与 gRPC 服务冲突
		// Prefork:   true,
		BodyLimit: 16 * 1024 * 1024,
	})

	app.Use(compress.New())

	app.Use(cors.New())

	app.Use(etag.New())

	app.Use(helmet.New())

	app.Use(func(c *fiber.Ctx) error {
		utility.Slogger.Info(c.Path(), "method", c.Method(), "query", c.Queries(), "ip", c.IP())
		return c.Next()
	})

	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Set(configuration.HeaderAPIVersion, "2024-02-03")
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		for _, item := range configuration.PublicUris {
			match, _ := regexp.MatchString(item, c.Path())
			if match {
				return c.Next()
			}
		}
		// auth := c.Get("Authorization")
		// auth = strings.Replace(auth, "Bearer ", "", 1)
		// token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")), nil
		// })
		// if err != nil {
		// 	utility.Slogger.Error(err.Error())
		// 	return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
		// }
		// if !token.Valid {
		// 	return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
		// }
		return c.Next()
	})

	app.Get("/crate-api/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/crate-api/subscriber/refresh-token", subscriber.RefreshJwt)
	app.Post("/crate-api/subscriber/log-in", subscriber.LogIn)
	app.Post("/crate-api/subscriber/sign-up", subscriber.SignUp)
	app.Post("/crate-api/subscriber/validate-token", subscriber.ValidateToken)
	app.Get("/crate-api/subscriber/:uuid/:id", subscriber.GetWithParams)

	app.Get("/crate-api/:schema/:table/:uuid/:id", schema.GetWithParams)
	// app.Post("/crate-api/:schema/:table", schema.Post)
	app.Put("/crate-api/:schema/:table/:uuid/:id", schema.Put)
	app.Delete("/crate-api/:schema/:table/:uuid/:id", schema.Delete)

	schemaRepo := NewSchemaRepoImpl(utility.Postgres)
	schemaService := NewSchemaService(schemaRepo)

	subscriberRepo := NewSubscriberRepoImpl(utility.Postgres)
	subscriberService := NewSubscriberService(subscriberRepo)

	AddSchemaEndpoints(app, schemaService)
	AddSubscriberEndpoints(app, subscriberService)

	log.Fatal(app.Listen(":" + port))
}
