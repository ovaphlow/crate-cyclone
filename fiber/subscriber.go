package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

/*
CREATE TABLE crate.subscriber (
	id int8 NOT NULL,
	email varchar NOT NULL,
	"name" varchar NOT NULL,
	phone varchar NOT NULL,
	tags jsonb NOT NULL,
	detail jsonb NOT NULL,
	"time" timestamptz NOT NULL,
	CONSTRAINT subscriber_pk PRIMARY KEY (id)
);
CREATE INDEX subscriber_email_idx ON crate.subscriber USING btree (email);
CREATE INDEX subscriber_name_idx ON crate.subscriber USING btree (name);
CREATE INDEX subscriber_phone_idx ON crate.subscriber USING btree (phone);
*/

var subscriberColumns = []string{"id", "email", "name", "phone", "tags", "detail", "time"}

type Subscriber struct {
	Id     int64
	Email  string
	Name   string
	Phone  string
	Tags   string
	Detail string
	Time   string
}

func endpointGetWithParams(c *fiber.Ctx) error {
	c.Set(HEADER_API_VERSION, "2024-01-06")
	id := c.Params("id", "")
	uuid := c.Params("uuid", "")
	if id == "" || uuid == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	id_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	subscriber, err := repoRetrieveSubscriberById(id_, uuid)
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if subscriber == nil {
		return c.Status(404).JSON(fiber.Map{"message": "用户不存在"})
	}
	return c.JSON(fiber.Map{
		"id":     subscriber.Id,
		"email":  subscriber.Email,
		"name":   subscriber.Name,
		"phone":  subscriber.Phone,
		"tags":   subscriber.Tags,
		"detail": subscriber.Detail,
		"time":   subscriber.Time,
		"_id":    strconv.FormatInt(subscriber.Id, 10),
	})
}

func endpointRefreshJwt(c *fiber.Ctx) error {
	c.Set(HEADER_API_VERSION, "2024-01-06")
	type Body struct {
		Token string `json:"token"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if body.Token == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")), nil
	})
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slogger.Error("token claims is not jwt.MapClaims")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if claims["exp"] == nil {
		slogger.Error("token claims exp is nil")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return c.Status(401).JSON(fiber.Map{"message": "token 已过期"})
	}
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")))
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.JSON(fiber.Map{"token": tokenString})
}

func endpointSignIn(c *fiber.Ctx) error {
	c.Set(HEADER_API_VERSION, "2024-01-06")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	jwtKey := []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", ""))
	type SignInBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body SignInBody
	if err := c.BodyParser(&body); err != nil {
		slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if body.Username == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	subscriber, err := repoRetrieveSubscriberByUsername(body.Username)
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if subscriber == nil {
		return c.Status(401).JSON(fiber.Map{"message": "用户名或密码错误"})
	}
	var detail map[string]interface{}
	if err := json.Unmarshal([]byte(subscriber.Detail), &detail); err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	salt, ok := detail["salt"].(string)
	if !ok {
		slogger.Error("salt is not string")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	key := []byte(salt)
	r := hmac.New(sha256.New, key)
	r.Write([]byte(body.Password))
	sha := hex.EncodeToString(r.Sum(nil))
	if sha != detail["sha"] {
		return c.Status(401).JSON(fiber.Map{"message": "用户名或密码错误"})
	}
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Issuer:    "crate",
		Subject:   strconv.FormatInt(subscriber.Id, 10),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.JSON(fiber.Map{"token": tokenString})
}

func endpointSignUp(c *fiber.Ctx) error {
	c.Set(HEADER_API_VERSION, "2024-01-06")
	type SignUpBody struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	var body SignUpBody
	if err := c.BodyParser(&body); err != nil {
		slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if (body.Email == "" && body.Name == "" && body.Phone == "") || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	subscriber, err := repoRetrieveSubscriberByUsername(body.Email)
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if subscriber != nil {
		return c.Status(401).JSON(fiber.Map{"message": "用户名已存在"})
	}
	subscriber = &Subscriber{
		Email: body.Email,
		Name:  body.Name,
		Phone: body.Phone,
		Tags:  "[]",
	}
	bytes := make([]byte, 8)
	_, err = rand.Read(bytes)
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	key := []byte(hex.EncodeToString(bytes))
	r := hmac.New(sha256.New, key)
	r.Write([]byte(body.Password))
	sha := hex.EncodeToString(r.Sum(nil))
	subscriber.Detail = fmt.Sprintf(`{"salt": "%s", "sha": "%s", "uuid": "%s"}`, key, sha, uuid.New())
	if err := repoCreateSubscriber(subscriber); err != nil {
		slogger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.JSON(fiber.Map{"message": "注册成功"})
}

func endpointValidateToken(c *fiber.Ctx) error {
	c.Set(HEADER_API_VERSION, "2024-01-06")
	type Body struct {
		Token string `json:"token"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		slogger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	if body.Token == "" {
		return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
	}
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(strings.ReplaceAll(os.Getenv("JWT_KEY"), " ", "")), nil
	})
	if err != nil {
		slogger.Error(err.Error())
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "用户凭证异常"})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slogger.Error("token claims is not jwt.MapClaims")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if claims["exp"] == nil {
		slogger.Error("token claims exp is nil")
		return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
	}
	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return c.Status(401).JSON(fiber.Map{"message": "token 已过期"})
	}
	return c.JSON(claims)
}

func repoCreateSubscriber(subscriber *Subscriber) error {
	q := fmt.Sprintf(
		"insert into crate.subscriber (%s) values ($1, $2, $3, $4, $5, $6, $7)",
		strings.Join(subscriberColumns, ", "),
	)
	statement, err := Postgres.Prepare(q)
	if err != nil {
		return err
	}
	node, err := snowflake.NewNode(1)
	if err != nil {
		slogger.Error(err.Error())
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(
		node.Generate(),
		subscriber.Email,
		subscriber.Name,
		subscriber.Phone,
		subscriber.Tags,
		subscriber.Detail,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}
	return nil
}

func repoRetrieveSubscriberById(id int64, uuid string) (*Subscriber, error) {
	q := fmt.Sprintf(
		`
		select %s from crate.subscriber
		where id = $1 and json_contains(detail, json_object('uuid', $2))
		limit 1
		`,
		strings.Join(subscriberColumns, ", "),
	)
	statement, err := Postgres.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	subscriber := &Subscriber{}
	err = statement.QueryRow(id, uuid).Scan(
		&subscriber.Id,
		&subscriber.Email,
		&subscriber.Name,
		&subscriber.Phone,
		&subscriber.Tags,
		&subscriber.Detail,
		&subscriber.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return subscriber, nil
}

func repoRetrieveSubscriberByUsername(username string) (*Subscriber, error) {
	q := fmt.Sprintf(
		"select %s from crate.subscriber where email = $1 or name = $2 or phone = $3",
		strings.Join(subscriberColumns, ", "),
	)
	statement, err := Postgres.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Query(username, username, username)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	var rows []Subscriber
	for result.Next() {
		var subscriber Subscriber
		err = result.Scan(
			&subscriber.Id,
			&subscriber.Email,
			&subscriber.Name,
			&subscriber.Phone,
			&subscriber.Tags,
			&subscriber.Detail,
			&subscriber.Time,
		)
		if err != nil {
			return nil, err
		}
		rows = append(rows, subscriber)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	if len(rows) > 1 {
		return nil, fmt.Errorf("duplicate subscriber")
	}
	return &rows[0], nil
}
