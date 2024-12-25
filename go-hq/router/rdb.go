package router

import (
	"encoding/json" // Added import for JSON handling
	"log"
	"net/http"
	"ovaphlow/crate/hq/dbutil"
	"ovaphlow/crate/hq/utility"
	"strings"
)

// LoadSharedRouter 加载路由
//
// 参数:
//   - mux: HTTP请求多路复用器
//   - prefix: 路由前缀
//   - service: 数据库服务实现
//
// 返回值:
//
//	无
func LoadRDBUtilRouter(mux *http.ServeMux, prefix string, service *dbutil.ApplicationServiceImpl) {
	mux.HandleFunc("DELETE "+prefix+"/rdb-util/{st}/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		st := r.PathValue("st")
		id := r.PathValue("id")

		err := service.Remove(st, "id='"+id+"'")
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			response := utility.CreateHTTPResponseRFC9457("删除失败", http.StatusInternalServerError, r)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := utility.CreateHTTPResponseRFC9457("删除成功", http.StatusOK, r)
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("PUT "+prefix+"/rdb-util/{st}/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		st := r.PathValue("st")
		id := r.PathValue("id")
		d := r.URL.Query().Get("d")

		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			response := utility.CreateHTTPResponseRFC9457("无效的请求体", http.StatusBadRequest, r)
			json.NewEncoder(w).Encode(response)
			return
		}
		data["id"] = id

		deprecated := false
		if d == "1" || d == "true" {
			deprecated = true
		}
		err := service.Update(st, data, "id='"+id+"'", deprecated)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			response := utility.CreateHTTPResponseRFC9457("更新失败", http.StatusInternalServerError, r)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := utility.CreateHTTPResponseRFC9457("更新成功", http.StatusOK, r)
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("GET "+prefix+"/rdb-util/{st}", func(w http.ResponseWriter, r *http.Request) {
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
		log.Println(f)
		columns := r.URL.Query().Get("c")
		var c []string
		if columns == "" {
			c = []string{}
		} else {
			c = strings.Split(columns, ",")
		}

		result, err := service.GetMany(st, c, f, last)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			response := utility.CreateHTTPResponseRFC9457("内部服务器错误", http.StatusInternalServerError, r)
			json.NewEncoder(w).Encode(response)
			return
		}

		json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("POST "+prefix+"/rdb-util/{st}", func(w http.ResponseWriter, r *http.Request) {
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
