package http

import (
	"github.com/flystary/agent/g"
	"io"
	"net/http"

	"github.com/toolkits/sys"
)

func run() {
	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if !g.Config().Http.Backdoor {
			w.Write([]byte("/run Disabled"))
			return
		}

		if g.IsTrustable(r.RemoteAddr) {
			if r.ContentLength == 0 {
				http.Error(w, "body is Blank", http.StatusBadRequest)
				return
			}

			bys, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			body := string(bys)
			out, err := sys.CmdOutBytes("sh", "-c", body)
			if err != nil {
				w.Write([]byte("exec fail: " + err.Error()))
				return
			}
			w.Write(out)
		} else {
			w.Write([]byte("no privilege"))
		}
	})
}
