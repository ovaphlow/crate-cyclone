package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"ovaphlow/crate/hq/infrastructure"
	"ovaphlow/crate/hq/router"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infrastructure.InitSlog()

	databaseType := os.Getenv("DATABASE_TYPE")
	if databaseType == "postgres" {
		infrastructure.InitPostgres()
	} else if databaseType == "mysql" {
		// Initialize MySQL
	} else {
		log.Panic("未设置数据库")
	}
}

func main() {
	mux := http.NewServeMux()

	// CORS 中间件
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "accept, content-type, content-length, accept-encoding, x-csrf-token, authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	})

	// 安全头中间件
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
	})

	// API 版本中间件
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", "2024-01-06")
	})

	// Ping 路由
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// 静态文件服务
	fs := http.FileServer(http.Dir("./html"))
	mux.Handle("/html/", http.StripPrefix("/html", fs))

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
		port = "8080"
	}
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, mux))
}

func determinTarget(r *http.Request) string {
	if r.URL.Path == "/service1" {
		return "http://service1.example.com"
	} else if r.URL.Path == "/service2" {
		return "http://service2.example.com"
	}
	return "http://default.example.com"
}
