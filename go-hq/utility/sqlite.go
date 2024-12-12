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

	SQLite, err = sql.Open("sqlite3", "file:"+dsn+"?mode=memory&cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	if err := SQLite.Ping(); err != nil {
		log.Println("连接数据库失败 SQLite")
		log.Fatal(err.Error())
	}

	log.Println("连接数据库成功 SQLite")

	// Defer the call to save the in-memory database to disk on program exit
	defer SaveSQLiteToDisk(dsn)
}

// SaveSQLiteToDisk writes the in-memory SQLite database to the original file
func SaveSQLiteToDisk(dsn string) {
	diskDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer diskDB.Close()

	_, err = SQLite.Exec("VACUUM INTO ?", dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("内存数据库已写入磁盘 SQLite")
}
