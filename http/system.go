package http

import (
	"net/http"
	"time"
)

func system() {
	http.HandleFunc("/system/date", func(w http.ResponseWriter, r *http.Request) {
		RenderJson(w, time.Now().Format("2006-01-02 15:04:05"))
	})
}
