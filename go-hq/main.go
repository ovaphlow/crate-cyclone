package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"ovaphlow/crate/hq/middleware"
	"ovaphlow/crate/hq/router"
	"ovaphlow/crate/hq/utility"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	utility.InitSlog()

	databaseType := os.Getenv("DATABASE_TYPE")
	if databaseType == "postgres" {
		utility.InitPostgres()
	} else {
		log.Panic("未设置数据库")
	}
}

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func main() {
	mux := http.NewServeMux()

	// Ping 路由
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	handler := applyMiddlewares(mux, middleware.APIVersionMiddleware, middleware.CORSMiddleware, middleware.SecurityHeadersMiddleware)
	log.Println("中间件已加载")

	// 静态文件服务
	fs := http.FileServer(http.Dir("./html"))
	mux.Handle("/html/", http.StripPrefix("/html", fs))
	log.Println("静态文件服务已加载至/html")

	// 注册服务路由
	router.RegisterServiceRouter(mux)

	// 动态代理
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

	// 定时健康检查
	sec := 15
	duration := time.Duration(sec) * time.Second
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			router.PerformHealthCheck(sec)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8421"
	}
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}

func determinTarget(r *http.Request) string {
	for _, service := range router.ServiceList {
		if strings.HasPrefix(r.URL.Path, "/proxy/"+service.Name) {
			return service.Protocol + "://" + service.Host + ":" + strconv.Itoa(service.Port)
		}
	}
	return ""
}
