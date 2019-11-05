// SPDX-License-Identifier: MIT

// Package docs docs 内容管理
package docs

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const indexPage = "index.xml"

var filterFiles = []string{
	"CNAME",
	"README.md",
}

// Handler 将 docs 当作一个网站进行管理
func Handler(docs string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pp := r.URL.Path

		if len(pp) > 0 && isFilterFile(path.Base(pp)) {
			errStatus(w, http.StatusNotFound)
			return
		}

		if len(pp) == 0 || pp[0] != '/' {
			pp = "/" + pp
			r.URL.Path = pp
		}

		pp = filepath.Clean(filepath.Join(docs, pp))
		info, err := os.Stat(pp)
		if err != nil {
			if os.IsNotExist(err) {
				errStatus(w, http.StatusNotFound)
				return
			}
			if os.IsPermission(err) {
				errStatus(w, http.StatusForbidden)
				return
			}

			errStatus(w, http.StatusInternalServerError)
			return
		}
		if info.IsDir() {
			pp = filepath.Clean(filepath.Join(pp, indexPage))
		}

		http.ServeFile(w, r, pp)
	})
}

func errStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func isFilterFile(filename string) bool {
	for _, file := range filterFiles {
		if file == filename {
			return true
		}
	}

	return false
}
