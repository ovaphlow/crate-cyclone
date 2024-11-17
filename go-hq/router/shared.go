package router

import (
	"net/http"
	"ovaphlow/crate/hq/dbutil"
)

func LoadSharedRouter(mux *http.ServeMux, prefix string, service *dbutil.ApplicationServiceImpl) {
	mux.HandleFunc("GET "+prefix+"/shared", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("shared"))
	})

	mux.HandleFunc("POST "+prefix+"/shared", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("shared"))
	})
}
