package router

import (
	"encoding/json"
	"net/http"
	"ovaphlow/crate/hq/subscriber"
)

func RegisterSubscriberRouter(mux *http.ServeMux, subscriberService *subscriber.SubscriberService) {
	mux.HandleFunc("/crate-hq-api/subscriber/sign-up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body subscriber.Subscriber
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 处理 body
	})
}
