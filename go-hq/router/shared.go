package router

import (
	"encoding/json" // Added import for JSON handling
	"log"
	"net/http"
	"ovaphlow/crate/hq/dbutil"
	"ovaphlow/crate/hq/utility"
)

// LoadSharedRouter 设置 GET 和 POST 路由。
// 该函数加载共享路由并配置 GET 和 POST 请求的处理程序。
//
// 参数:
//   - mux: 用于注册路由和处理程序的 HTTP 请求多路复用器。
//   - prefix: 定义路由基本路径的路由前缀。
//   - service: 用于处理业务逻辑的应用服务实例。
//
// 返回值:
//   - 无
func LoadSharedRouter(mux *http.ServeMux, prefix string, service *dbutil.ApplicationServiceImpl) {
	mux.HandleFunc("GET "+prefix+"/dbutil/{st}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		st := r.PathValue("st")
		last := r.URL.Query().Get("l")
		filter := r.URL.Query().Get("f")
		f, err := utility.ConvertQueryStringToDefaultFilter(filter)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			response := utility.CreateHTTPResponseRFC9457("无效的查询参数", http.StatusBadRequest, r)
			json.NewEncoder(w).Encode(response)
			return
		}

		result, err := service.GetMany(st, f, last)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			response := utility.CreateHTTPResponseRFC9457("内部服务器错误", http.StatusInternalServerError, r)
			json.NewEncoder(w).Encode(response)
			return
		}

		json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("POST "+prefix+"/dbutil/{st}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		st := r.PathValue("st")

		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			response := utility.CreateHTTPResponseRFC9457("无效的请求体", http.StatusBadRequest, r)
			json.NewEncoder(w).Encode(response)
			return
		}

		err := service.Create(st, data)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			response := utility.CreateHTTPResponseRFC9457("创建失败", http.StatusInternalServerError, r)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusCreated)
		response := utility.CreateHTTPResponseRFC9457("创建成功", http.StatusCreated, r)
		json.NewEncoder(w).Encode(response)
	})
}
