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

	// 根据环境变量初始化对应的数据库
	databaseType := os.Getenv("DATABASE_TYPE")
	if databaseType == "postgres" {
		utility.InitPostgres()
	} else if databaseType == "mysql" {
		utility.InitMySQL()
	} else {
		log.Panic("未设置数据库")
	}
}

type Middleware func(http.Handler) http.Handler

// applyMiddlewares 将给定的中间件应用到一个 HTTP 处理器上。
// 参数:
//   - h: 初始的 http.Handler，后续中间件将应用到它上面。
//   - middlewares: 可变参数列表的 Middleware 函数，依次应用。
//
// 返回值:
//   - 应用了所有中间件的 http.Handler。
func applyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func main() {
	// 创建一个新的ServeMux
	mux := http.NewServeMux()

	// 定义Ping路由，返回"pong"
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 应用多个中间件到mux
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

	// 设置定时健康检查，每15秒执行一次
	sec := 15
	duration := time.Duration(sec) * time.Second
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			router.PerformHealthCheck(sec)
		}
	}()

	// 初始化数据库连接并创建共享资源仓库
	databaseType := os.Getenv("DATABASE_TYPE")
	var rdbRepo dbutil.RDBRepo
	if databaseType == "postgres" {
		rdbRepo = dbutil.NewPostgresRepo(utility.Postgres)
	} else if databaseType == "mysql" {
		rdbRepo = dbutil.NewMySQLRepo(utility.MySQL)
	} else {
		log.Fatal("Unsupported DATABASE_TYPE")
	}

	// 创建应用服务并加载共享路由
	appService := dbutil.NewApplicationService(rdbRepo)
	router.LoadRDBUtilRouter(mux, "/cyclone-api", appService)

	// 获取端口号并启动HTTP服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8421"
	}
	log.Println("0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}

// 根据请求决定目标服务的地址
func determinTarget(r *http.Request) string {
	for _, service := range router.ServiceList {
		if strings.HasPrefix(r.URL.Path, "/proxy/"+service.Name) {
			return service.Protocol + "://" + service.Host + ":" + strconv.Itoa(service.Port)
		}
	}
	return ""
}
