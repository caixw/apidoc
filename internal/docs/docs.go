// SPDX-License-Identifier: MIT

// Package docs 打包文档内容
package docs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/issue9/source"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/docs"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// FileHeader 表示输出文件的文件头内容
const FileHeader = "该文件由工具自动生成，请勿手动修改！"

// 默认页面
const indexPage = "index.xml"

// 指定在 Handler 中，folder 不为空时，可以访问的文件列表。
//
// 可以以前缀的方式指定，比如：v5/ 表示以 v5/ 开头的所有文件。
var styles = []string{
	"icon.svg",
	ast.MajorVersion + "/",
}

var docsDir = core.FileURI(source.CurrentPath("../../docs"))

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
//
// 如果 folder 为空，表示采用内嵌的数据作为文件服务；
// stylesheet 是否只返回最基本的样式表相关文件。
func Handler(folder core.URI, stylesheet bool) http.Handler {
	if folder == "" {
		return fsHandler(docs.FS, stylesheet)
	}

	switch scheme, path := folder.Parse(); scheme {
	case core.SchemeFile, "":
		return fsHandler(os.DirFS(path), stylesheet)
	case core.SchemeHTTP, core.SchemeHTTPS:
		return remoteHandler(folder, stylesheet)
	default:
		panic(locale.NewError(locale.ErrInvalidURIScheme, scheme))
	}
}

func fsHandler(fsys fs.FS, stylesheet bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pp := r.URL.Path
		if pp == "" || pp == "/" {
			pp = indexPage
		}
		if pp[0] == '/' {
			pp = pp[1:]
		}
		if pp[len(pp)-1] == '/' {
			pp = path.Join(pp, indexPage)
		}

	READ:
		if stylesheet && !isStylesheetFile(pp) {
			errStatus(w, http.StatusNotFound)
			return
		}

		f, err := fsys.Open(pp)
		if errors.Is(err, fs.ErrNotExist) {
			errStatus(w, http.StatusNotFound)
			return
		} else if errors.Is(err, fs.ErrPermission) {
			errStatus(w, http.StatusForbidden)
			return
		} else if err != nil {
			errStatusWithError(w, err)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			errStatusWithError(w, err)
			return
		}
		if stat.IsDir() {
			pp = path.Join(pp, indexPage)
			goto READ
		}

		data, err := io.ReadAll(f)
		if err != nil {
			errStatusWithError(w, err)
			return
		}

		http.ServeContent(w, r, r.URL.Path, stat.ModTime(), bytes.NewReader(data))
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

func errStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func errStatusWithError(w http.ResponseWriter, err error) {
	herr, ok := err.(*core.HTTPError)
	if ok {
		http.Error(w, herr.Err.Error(), herr.Code)
		return
	}

	fmt.Println("ERR:", err)
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
