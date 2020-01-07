// SPDX-License-Identifier: MIT

// Package docs 打包文档内容
package docs

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	xpath "github.com/caixw/apidoc/v6/internal/path"
	"github.com/caixw/apidoc/v6/internal/vars"
)

// FileInfo 被打包文件的信息
type FileInfo struct {
	// 相对于打包根目录的地址，同时也会被作为路由地址
	Name string

	ContentType string
	Content     []byte
}

// 默认页面
const indexPage = "index.xml"

// 指定在 TypeStylesheet 下需要的文件列表。
//
// 可以以前缀的方式指定，比如：v5/ 表示以 v5/ 开头的所有文件。
var styles = []string{
	"icon.svg",
	vars.DocVersion() + "/", // v5/ 仅支持当前的文档版本
}

// Dir 指向 /docs 的路径
func Dir() string {
	return Path("")
}

// Path 指赂 /docs 下的 p 路径
func Path(p string) string {
	return filepath.Join(xpath.CurrPath("../../docs"), p)
}

// Handler 返回文件服务中间件
func Handler(folder string, stylesheet bool) http.Handler {
	if folder == "" {
		return embeddedHandler(stylesheet)
	}
	return folderHandler(folder, stylesheet)
}

func embeddedHandler(stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pp := r.URL.Path

		if len(pp) > 0 && pp[0] == '/' {
			pp = pp[1:]
		}
		r.URL.Path = pp
		indexPath := path.Join(pp, indexPage)

		for _, info := range data {
			if info.Name == pp || info.Name == indexPath {
				if stylesheet && !isStylesheetFile(info.Name) {
					errStatus(w, http.StatusNotFound)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", info.ContentType)
				w.Write(info.Content)
				return
			}
		}
		errStatus(w, http.StatusNotFound)
	})
}

func folderHandler(folder string, stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path

		if len(p) > 0 && p[0] == '/' {
			p = p[1:]
			r.URL.Path = p
		}

		if stylesheet && !isStylesheetFile(p) {
			errStatus(w, http.StatusNotFound)
			return
		}

		p = filepath.Clean(filepath.Join(folder, p))
		info, err := os.Stat(p)
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
			p = filepath.Clean(filepath.Join(p, indexPage))
		}

		http.ServeFile(w, r, p)
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
