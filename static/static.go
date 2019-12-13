// SPDX-License-Identifier: MIT

// Package static 用于打包静态文件内容
package static

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const indexPage = "index.xml"

var styles = []string{
	"icon.svg",
	"v5/",
}

// EmbeddedHandler 将 data 作为一个静态文件服务内容进行管理
func EmbeddedHandler(data []*FileInfo, stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pp := r.URL.Path

		if len(pp) > 0 && pp[0] == '/' {
			pp = pp[1:]
		}
		r.URL.Path = pp
		indexPath := path.Join(pp, indexPage)

		if stylesheet && !isStylesheetFile(pp) {
			errStatus(w, http.StatusNotFound)
			return
		}

		for _, info := range data {
			if info.Name == pp || info.Name == indexPath {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", info.ContentType)
				w.Write(info.Content)
				return
			}
		}
		errStatus(w, http.StatusNotFound)
	})
}

// FolderHandler 将 folder 当作一个网站进行管理
func FolderHandler(folder string, stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
			r.URL.Path = path
		}

		if stylesheet && !isStylesheetFile(path) {
			errStatus(w, http.StatusNotFound)
			return
		}

		path = filepath.Clean(filepath.Join(folder, path))
		info, err := os.Stat(path)
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
			path = filepath.Clean(filepath.Join(path, indexPage))
		}

		http.ServeFile(w, r, path)
	})
}

func errStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func isStylesheetFile(filename string) bool {
	if len(filename) > 0 && filename[0] == '/' {
		filename = filename[1:]
	}

	for _, file := range styles {
		if file == filename || strings.HasPrefix(filename, file) {
			return true
		}
	}

	return false
}
