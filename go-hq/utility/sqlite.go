package utility

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var SQLite *sql.DB

func InitSQLite() {
	err := godotenv.Load()
	if err != nil {
		Slogger.Error("环境变量未设置 SQLite")
		log.Fatal(err.Error())
	}

	dsn := os.Getenv("SQLITE_DATABASE")
	if dsn == "" {
		log.Fatal(errors.New("环境变量未设置 SQLite"))
	}

	SQLite, err = sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := SQLite.Ping(); err != nil {
		log.Println("连接数据库失败 SQLite")
		log.Fatal(err.Error())
	}

	log.Println("连接数据库成功 SQLite")
}
