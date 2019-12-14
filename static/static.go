// SPDX-License-Identifier: MIT

// Package static 提供了对 docs 中内容的处理方式
package static

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Type 表示文件的类型
type Type int8

// 几种文件类型的定义
const (
	TypeNone = iota
	TypeAll
	TypeStylesheet
)

// 默认页面
const indexPage = "index.xml"

// 指定了 xml 文档所需的 xsl 内容。
//
// 可以以前缀的方式指定，比如：v5/ 表示以 v5/ 开头的所有文件。
var styles = []string{
	"icon.svg",
	"v5/",
}

// EmbeddedHandler 将由 Pack 打包的内容当作一个文件服务中间件
func EmbeddedHandler(data []*FileInfo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pp := r.URL.Path

		if len(pp) > 0 && pp[0] == '/' {
			pp = pp[1:]
		}
		r.URL.Path = pp
		indexPath := path.Join(pp, indexPage)

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

// FolderHandler 将 folder 当作文件服务中间件
func FolderHandler(folder string, t Type) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t == TypeNone {
			errStatus(w, http.StatusNotFound)
			return
		}

		path := r.URL.Path

		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
			r.URL.Path = path
		}

		if t == TypeStylesheet && !isStylesheetFile(path) {
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
