package utility

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var MySQL *sql.DB

func InitMySQL() {
	err := godotenv.Load()
	if err != nil {
		Slogger.Error("加载环境变量失败")
		log.Fatal(err.Error())
	}
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		password,
		host,
		port,
		database,
	)
	MySQL, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	MySQL.SetConnMaxLifetime(time.Minute * 3)
	MySQL.SetMaxOpenConns(10)
	MySQL.SetMaxIdleConns(10)
	if err = MySQL.Ping(); err != nil {
		log.Println("连接MySQL数据库失败")
		log.Fatal(err.Error())
	}
	log.Println("连接MySQL数据库成功")
}
