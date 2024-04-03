package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Postgres *sql.DB

func InitPostgres() {
	err := godotenv.Load()
	if err != nil {
		slogger.Error("加载环境变量失败")
	}
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DATABASE")
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		database,
	)
	Postgres, err = sql.Open("postgres", dsn)
	if err != nil {
		slogger.Error(err.Error())
	}
	Postgres.SetConnMaxLifetime(time.Second * 30)
	Postgres.SetMaxIdleConns(2)
	if err = Postgres.Ping(); err != nil {
		slogger.Error("连接数据库失败")
		slogger.Error(err.Error())
	}
}
