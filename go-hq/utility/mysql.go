package utility

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var MySQL *sql.DB

func InitMySQL(user, password, host, port, database string) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		password,
		host,
		port,
		database,
	)
	var err error
	MySQL, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	MySQL.SetConnMaxLifetime(time.Minute * 3)
	cpuCount := runtime.NumCPU()
	MySQL.SetMaxOpenConns(cpuCount*2 + 1)
	MySQL.SetMaxIdleConns(cpuCount*2 + 1)
	if err = MySQL.Ping(); err != nil {
		log.Println("连接数据库失败 MySQL")
		log.Fatal(err.Error())
	}
	log.Println("连接数据库成功 MySQL")
}
