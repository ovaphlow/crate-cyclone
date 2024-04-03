package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var MySQL *sql.DB

func InitMySQL() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("加载环境变量失败")
	}
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	MySQL, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err.Error())
	}
	if err = MySQL.Ping(); err != nil {
		log.Println("连接数据库失败")
		log.Fatal(err.Error())
	}
	MySQL.SetConnMaxLifetime(time.Second * 30)
	MySQL.SetMaxIdleConns(runtime.NumCPU()*2 + 1)
	// MySQL.SetMaxOpenConns(runtime.NumCPU()*2 + 1)
}
