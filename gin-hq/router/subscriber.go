package router

import (
	"log"
	"os"
	"ovaphlow/crate/hq/subscriber"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RegisterSubscriberRouter(r *gin.Engine, s *subscriber.SubscriberService) {
	type body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	r.POST("/crate-hq-api/subscriber/sign-up", func(c *gin.Context) {
		var b body
		if err := c.ShouldBindJSON(&b); err != nil {
			c.Error(err)
			return
		}
		err := s.SignUp(&subscriber.Subscriber{
			Email:  b.Username,
			Name:   b.Username,
			Detail: b.Password,
		})
		if err != nil {
			c.Error(err)
			return
		}
		c.Status(200)
	})

	r.POST("/crate-hq-api/subscriber/log-in", func(c *gin.Context) {
		var b body
		if err := c.ShouldBindJSON(&b); err != nil {
			c.Error(err)
			return
		}
		err := godotenv.Load()
		if err != nil {
			log.Fatal("加载环境变量失败")
		}
		jwtKey := []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", ""))
		token, err := s.LogIn(b.Username, b.Password, string(jwtKey))
		if err != nil {
			c.Error(err)
			return
		}
		log.Println(1111)
		log.Println(token)
		c.JSON(200, gin.H{
			"token": token,
		})
	})
}
