package main

// Import necessary packages
import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"ovaphlow/crate/hq/dbutil"
	"ovaphlow/crate/hq/middleware"
	"ovaphlow/crate/hq/router"
	"ovaphlow/crate/hq/utility"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize structured logging
	utility.InitSlog()

	// Initialize the corresponding database based on environment variables
	database_type := os.Getenv("DATABASE_TYPE")
	if database_type == "posgres" {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		database := os.Getenv("POSTGRES_DATABASE")
		utility.InitPostgres(user, password, host, port, database)
	} else if database_type == "mysql" {
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")
		database := os.Getenv("MYSQL_DATABASE")
		utility.InitMySQL(user, password, host, port, database)
	} else {
		log.Println("未设置有效的数据库类型 DATABASE_TYPE (postgres/mysql)")
		log.Println("不初始化 RDB 连接")
	}

	sqlite := os.Getenv("SQLITE_ENABLED")
	if sqlite == "true" {
		utility.InitSQLite()
	}
}

type Middleware func(http.Handler) http.Handler

// applyMiddlewares applies the given middlewares to an HTTP handler.
// Parameters:
//   - h: The initial http.Handler to which the subsequent middlewares will be applied.
//   - middlewares: Variadic parameter list of Middleware functions to be applied in sequence.
//
// Returns:
//   - An http.Handler with all middlewares applied.
func applyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define Ping route, returns "pong"
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Apply multiple middlewares to mux
	handler := applyMiddlewares(mux, middleware.APIVersionMiddleware, middleware.CORSMiddleware, middleware.SecurityHeadersMiddleware)
	log.Println("中间件已加载")

	// Set up static file service, path is /html
	fs := http.FileServer(http.Dir("./html"))
	mux.Handle("/html/", http.StripPrefix("/html", fs))
	log.Println("静态文件服务已加载至 /html")

	// Register service routes
	router.RegisterServiceRouter(mux)

	// Set up dynamic proxy route
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

	// Set up periodic health check, executed every 15 seconds
	sec := 15
	duration := time.Duration(sec) * time.Second
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			router.PerformHealthCheck(sec)
		}
	}()

	// Initialize database connection and create shared resource repository
	databaseType := os.Getenv("DATABASE_TYPE")
	var rdbRepo dbutil.Repo
	if databaseType == "postgres" {
		rdbRepo = dbutil.NewPostgresRepo(utility.Postgres)
	} else if databaseType == "mysql" {
		rdbRepo = dbutil.NewMySQLRepo(utility.MySQL)
	} else {
		log.Fatal("Unsupported DATABASE_TYPE")
	}

	// Create application service and load shared routes
	appService := dbutil.NewApplicationService(rdbRepo)
	router.LoadRDBUtilRouter(mux, "/cyclone-api", appService)

	edb := os.Getenv("SQLITE_ENABLED")
	if edb == "true" {
		edbRepo := dbutil.NewSQLiteRepo(utility.SQLite)
		edbService := dbutil.NewApplicationService(edbRepo)
		router.LoadEDBUtilRouter(mux, "/cyclone-api", edbService)
	}

	// Ensure SQLite database is saved to disk on program exit
	sqlite := os.Getenv("SQLITE_ENABLED")
	if sqlite == "true" {
		dsn := os.Getenv("SQLITE_DATABASE")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			utility.PersistSQLite(dsn)
			os.Exit(0)
		}()

		// Periodically persist SQLite database to disk at specific times
		go func() {
			for {
				now := time.Now()
				next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
				if now.Hour() >= 20 {
					next = next.Add(4 * time.Hour)
				} else if now.Hour() >= 16 {
					next = next.Add(4 * time.Hour)
				} else if now.Hour() >= 12 {
					next = next.Add(4 * time.Hour)
				} else if now.Hour() >= 8 {
					next = next.Add(4 * time.Hour)
				} else {
					next = next.Add(8 * time.Hour)
				}
				time.Sleep(time.Until(next))
				utility.PersistSQLite(dsn)
			}
		}()
	}

	// Get port number and start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8421"
	}
	log.Println("0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}

// determineTarget decides the target service address based on the request
func determinTarget(r *http.Request) string {
	for _, service := range router.ServiceList {
		if strings.HasPrefix(r.URL.Path, "/proxy/"+service.Name) {
			return service.Protocol + "://" + service.Host + ":" + strconv.Itoa(service.Port)
		}
	}
	return ""
}
