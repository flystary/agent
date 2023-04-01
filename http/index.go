package http

import (
	"github.com/flystary/agent/g"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/toolkits/file"
)

func index() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			if !file.IsExist(filepath.Join(g.Root, "/static", r.URL.Path, "index.html")) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Root, "/static"))).ServeHTTP(w, r)
	})
}
