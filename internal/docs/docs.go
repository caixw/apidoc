// SPDX-License-Identifier: MIT

// Package docs 打包文档内容
package docs

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/issue9/utils"
)

// 默认页面
const indexPage = "index.xml"

// 指定在 Handler 中，folder 不为空时，可以访问的文件列表。
//
// 可以以前缀的方式指定，比如：v5/ 表示以 v5/ 开头的所有文件。
var styles = []string{
	"icon.svg",
	ast.MajorVersion + "/",
}

var docsDir = core.FileURI(utils.CurrentPath("../../docs"))

// FileInfo 被打包文件的信息
type FileInfo struct {
	// 相对于打包根目录的地址，同时也会被作为路由地址
	Name string

	ContentType string
	Content     []byte
}

// Dir 指向 /docs 的路径
func Dir() core.URI {
	return docsDir
}

// StylesheetURL 生成 apidoc.xsl 文件的 URL 地址
//
// 相对于 docs 目录
func StylesheetURL(prefix string) string {
	if prefix == "" {
		return ast.MajorVersion + "/apidoc.xsl"
	}
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}
	return prefix + ast.MajorVersion + "/apidoc.xsl"
}

// Handler 返回文件服务中间件
func Handler(folder core.URI, stylesheet bool) http.Handler {
	if folder == "" {
		return embeddedHandler(stylesheet)
	}

	switch scheme, _ := folder.Parse(); scheme {
	case core.SchemeFile, "":
		return localHandler(folder, stylesheet)
	case core.SchemeHTTP, core.SchemeHTTPS:
		return remoteHandler(folder, stylesheet)
	default:
		panic(locale.NewError(locale.ErrInvalidURIScheme))
	}
}

func embeddedHandler(stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pp := r.URL.Path

		if len(pp) > 0 && pp[0] == '/' {
			pp = pp[1:]
			r.URL.Path = pp
		}
		indexPath := path.Join(pp, indexPage)

		for _, info := range data {
			if info.Name == pp || info.Name == indexPath {
				if stylesheet && !isStylesheetFile(info.Name) {
					errStatus(w, http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", info.ContentType)
				w.WriteHeader(http.StatusOK)
				w.Write(info.Content)
				return
			}
		}
		errStatus(w, http.StatusNotFound)
	})
}

func remoteHandler(url core.URI, stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path

		if stylesheet && !isStylesheetFile(p) {
			errStatus(w, http.StatusNotFound)
			return
		}

		uri := url.Append(p)
		data, err := uri.ReadAll(nil)
		if err != nil {
			httpError, ok := err.(*core.HTTPError)
			if !ok {
				errStatus(w, http.StatusInternalServerError)
				return
			}

			if httpError.Code != http.StatusNotFound {
				errStatusWithError(w, httpError)
				return
			}

			data, err = uri.Append(indexPage).ReadAll(nil)
			if err != nil {
				errStatusWithError(w, err)
				return
			}
		}

		w.Write(data)
	})
}

func localHandler(folder core.URI, stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path

		if stylesheet && !isStylesheetFile(p) {
			errStatus(w, http.StatusNotFound)
			return
		}

		p, err := folder.Append(p).File()
		if err != nil {
			errStatus(w, http.StatusInternalServerError)
			return
		}

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

func errStatusWithError(w http.ResponseWriter, err error) {
	herr, ok := err.(*core.HTTPError)
	if ok {
		http.Error(w, herr.Err.Error(), herr.Code)
		return
	}

	errStatus(w, http.StatusInternalServerError)
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
