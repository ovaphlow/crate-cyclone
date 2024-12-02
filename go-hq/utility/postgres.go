package utility

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"time"

	_ "github.com/lib/pq"
)

var Postgres *sql.DB

func InitPostgres(user, password, host, port, database string) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		database,
	)
	var err error
	Postgres, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	Postgres.SetConnMaxLifetime(time.Second * 30)
	cpuCount := runtime.NumCPU()
	Postgres.SetMaxOpenConns(cpuCount*2 + 1)
	Postgres.SetMaxIdleConns(cpuCount*2 + 1)
	if err = Postgres.Ping(); err != nil {
		log.Println("连接数据库失败 Postgres")
		log.Fatal(err.Error())
	}
	log.Println("连接数据库成功 Postgres")
}
