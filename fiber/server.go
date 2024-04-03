package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt"
)

func Serve(addr string) {
	InitMySQL()

	InitPostgres()

	app := fiber.New(fiber.Config{
		Prefork:   true,
		BodyLimit: 16 * 1024 * 1024,
	})

	app.Use(compress.New())

	app.Use(cors.New())

	app.Use(etag.New())

	app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format:     "${ip} [${time}] ${status} - ${latency} ${method} ${path} ?${queryParams}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Shanghai",
	}))

	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		for _, item := range PUBLIC_URIS {
			match, _ := regexp.MatchString(item, c.Path())
			if match {
				slogger.Info("public uri", "match", c.Path())
				return c.Next()
			}
		}
		auth := c.Get("Authorization")
		auth = strings.Replace(auth, "Bearer ", "", 1)
		token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
			return []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")), nil
		})
		if err != nil {
			slogger.Error(err.Error())
			return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
		}
		if !token.Valid {
			return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
		}
		return c.Next()
	})

	app.Get("/crate-api/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "hola el mondo",
		})
	})

	app.Get("/crate-api/event", EventEndpointGet)
	app.Post("/crate-api/subscriber/refresh-jwt", endpointRefreshJwt)
	app.Post("/crate-api/subscriber/sign-in", endpointSignIn)
	app.Post("/crate-api/subscriber/sign-up", endpointSignUp)
	app.Post("/crate-api/subscriber/validate-token", endpointValidateToken)
	app.Get("/crate-api/subscriber/:uuid/:id", endpointGetWithParams)

	log.Fatal(app.Listen(addr))
}
