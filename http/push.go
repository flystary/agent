package http

import (
	"encoding/json"
	"net/http"

	"github.com/flystary/aiops/model"
)

func push() {
	http.HandleFunc("/v1/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body is Blank", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var metrics []*model.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "connot decode Body", http.StatusBadRequest)
			return
		}

		// g.SendToTransfer(metrics)
		w.Write([]byte("success"))
	})
}
