package main

// 导入必要的包
import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"ovaphlow/crate/hq/dbutil"
	"ovaphlow/crate/hq/middleware"
	"ovaphlow/crate/hq/router"
	"ovaphlow/crate/hq/utility"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// 加载环境变量文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 初始化结构化日志
	utility.InitSlog()

	// 初始化 PostgreSQL 数据库
	postgres_enabled := os.Getenv("POSTGRES_ENABLED")
	if postgres_enabled == "true" || postgres_enabled == "1" {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		database := os.Getenv("POSTGRES_DATABASE")
		utility.InitPostgres(user, password, host, port, database)
	}

	// 初始化 MySQL 数据库
	mysql_enabled := os.Getenv("MYSQL_ENABLED")
	if mysql_enabled == "true" || mysql_enabled == "1" {
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")
		database := os.Getenv("MYSQL_DATABASE")
		utility.InitMySQL(user, password, host, port, database)
	}

	// 初始化 SQLite 数据库
	sqlite_enabled := os.Getenv("SQLITE_ENABLED")
	if sqlite_enabled == "true" || sqlite_enabled == "1" {
		utility.InitSQLite()
	}
}

type Middleware func(http.Handler) http.Handler

// applyMiddlewares 应用给定的中间件到 HTTP 处理器。
// 参数:
//   - h: 初始的 http.Handler，后续的中间件将应用于此。
//   - middlewares: 可变参数列表，包含依次应用的中间件函数。
//
// 返回值:
//   - 一个应用了所有中间件的 http.Handler。
func applyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func main() {
	// 创建一个新的 ServeMux
	mux := http.NewServeMux()

	// 定义 Ping 路由，返回 "pong"
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 应用多个中间件到 mux
	handler := applyMiddlewares(mux, middleware.APIVersionMiddleware, middleware.CORSMiddleware, middleware.SecurityHeadersMiddleware)
	log.Println("中间件已加载")

	// 设置静态文件服务，路径为 /html
	fs := http.FileServer(http.Dir("./html"))
	mux.Handle("/html/", http.StripPrefix("/html", fs))
	log.Println("静态文件服务已加载至 /html")

	// 注册服务路由
	router.RegisterServiceRouter(mux)

	// 设置动态代理路由
	mux.HandleFunc("/proxy/", func(w http.ResponseWriter, r *http.Request) {
		target := determinTarget(r)
		remote, err := url.Parse(target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/proxy")
		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = remote.Host
		r.Header.Set("X-Auth", "1123")
		proxy.ServeHTTP(w, r)
	})

	// 设置定期健康检查，每 15 秒执行一次
	sec := 15
	duration := time.Duration(sec) * time.Second
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			router.PerformHealthCheck(sec)
		}
	}()

	// 加载 PostgreSQL 路由
	postgres_enabled := os.Getenv("POSTGRES_ENABLED")
	if postgres_enabled == "true" || postgres_enabled == "1" {
		postgresRepo := dbutil.NewPostgresRepo(utility.Postgres)
		postgresService := dbutil.NewApplicationService(postgresRepo)
		router.LoadPostgresRouter(mux, "/crate-api-data", postgresService)
	}

	// 加载 MySQL 路由
	mysql_enabled := os.Getenv("MYSQL_ENABLED")
	if mysql_enabled == "true" || mysql_enabled == "1" {
		mysqlRepo := dbutil.NewMySQLRepo(utility.MySQL)
		mysqlService := dbutil.NewApplicationService(mysqlRepo)
		router.LoadMySQLRouter(mux, "/crate-api-data", mysqlService)
	}

	// 加载 SQLite 路由
	sqlite_enabled := os.Getenv("SQLITE_ENABLED")
	if sqlite_enabled == "true" || sqlite_enabled == "1" {
		sqliteRepo := dbutil.NewSQLiteRepo(utility.SQLite)
		sqliteService := dbutil.NewApplicationService(sqliteRepo)
		router.LoadSQLiteRouter(mux, "/crate-api-data", sqliteService)
	}

	// 获取端口号并启动 HTTP 服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8421"
	}
	log.Println("0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}

// determineTarget 根据请求决定目标服务地址
func determinTarget(r *http.Request) string {
	for _, service := range router.ServiceList {
		if strings.HasPrefix(r.URL.Path, "/proxy/"+service.Name) {
			return service.Protocol + "://" + service.Host + ":" + strconv.Itoa(service.Port)
		}
	}
	return ""
}
