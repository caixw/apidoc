// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具
//
// 从代码文件的注释中提取特定格式的内容，生成 RESTful API 文档，支持大部分的主流的编程语言。
package apidoc

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp"
)

const (
	// LSPVersion 当前支持的 language server protocol 版本
	LSPVersion = lsp.Version

	// DocVersion 文档的版本
	DocVersion = ast.Version
)

// Config 配置文件 apidoc.yaml 所表示的内容
type Config = build.Config

// SetLocale 设置当前的本地化 ID
//
// 如果不调用此函数，则默认会采用 internal/locale.DefaultLocaleID 的值。
// 如果想采用当前系统的本地化信息，可以使用
// github.com/issue9/localeutil.SystemLanguageTag 函数。
func SetLocale(tag language.Tag) { locale.SetTag(tag) }

// Locale 获取当前设置的本地化 ID
func Locale() language.Tag { return locale.Tag() }

// Locales 返回当前所有支持的本地化信息
func Locales() []language.Tag { return locale.Tags() }

// Version 当前程序的版本号
//
// full 表示是否需要在版本号中包含编译日期和编译时的 Git 记录 ID。
func Version(full bool) string {
	if full {
		return core.FullVersion()
	}
	return core.Version()
}

// Build 解析文档并输出文档内容
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *core.Error 类型返回错误信息。
//
// NOTE: 如果需要从配置文件进行构建文档，可以采用 Config.Build
func Build(h *core.MessageHandler, o *build.Output, i ...*build.Input) error {
	return build.Build(h, o, i...)
}

// Buffer 生成文档内容并返回
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *core.Error 类型返回错误信息。
//
// NOTE: 如果需要从配置文件进行构建文档，可以采用 Config.Buffer
func Buffer(h *core.MessageHandler, o *build.Output, i ...*build.Input) (*bytes.Buffer, error) {
	return build.Buffer(h, o, i...)
}

// CheckSyntax 测试文档语法
func CheckSyntax(h *core.MessageHandler, i ...*build.Input) error {
	return build.CheckSyntax(h, i...)
}

// ServeLSP 提供 language server protocol 服务
//
// header 表示传递内容是否带报头；
// t 表示允许连接的类型，目前可以是 tcp、udp、stdio 和 unix；
// timeout 表示服务端每次读取客户端时的超时时间，如果为 0 表示不会超时。
// 超时并不会出错，而是重新开始读取数据，防止被读取一直阻塞，无法结束进程；
func ServeLSP(header bool, t, addr string, timeout time.Duration, info, erro *log.Logger) error {
	return lsp.Serve(header, t, addr, timeout, info, erro)
}

// Static 为 dir 指向的路径内容搭建一个静态文件服务
//
// dir 为静态文件的根目录，一般指向 /docs
// 用于搭建一个本地版本的 https://apidoc.tools，默认页为 index.xml。
// 如果 dir 值为空，则会采用内置的文档内容作为静态文件服务的内容。
//
// stylesheet 表示是否只展示 XSL 及相关的内容。
//
// 用户可以通过以下代码搭建一个简易的 https://apidoc.tools 网站：
//  http.Handle("/apidoc", apidoc.Static(...))
func Static(dir core.URI, stylesheet bool, erro *log.Logger) http.Handler {
	return docs.Handler(dir, stylesheet, erro)
}

// Server 用于生成查看文档中间件的配置项
type Server struct {
	Status      int         // 默认值为 200
	Path        string      // 文档在路由中的地址，默认值为 apidoc.xml
	ContentType string      // 文档的 ContentType，为空表示采用 application/xml
	Dir         core.URI    // 除文档不之外的附加项，比如 xsl，css 等内容的所在位置，如果为空表示采用内嵌的数据；
	Stylesheet  bool        // 是否只采用 Dir 中的 xsl 和 css 等样式数据，而忽略其它文件
	Erro        *log.Logger // 服务出错时的错误信息输出通道，默认采用 log.Default()
}

func (srv *Server) sanitize() {
	if srv.Status == 0 {
		srv.Status = http.StatusOK
	}

	if srv.Path == "" {
		srv.Path = "/apidoc.xml"
	}

	if srv.ContentType == "" {
		srv.ContentType = "application/xml"
	}

	if srv.Erro == nil {
		srv.Erro = log.Default()
	}
}

// Buffer 将 buf 作为文档内容生成中间件
func (srv *Server) Buffer(buf []byte) http.Handler {
	srv.sanitize()

	buf = addStylesheet(buf)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == srv.Path {
			w.Header().Set("Content-Type", srv.ContentType)
			w.WriteHeader(srv.Status)
			w.Write(buf)
			return
		}

		Static(srv.Dir, srv.Stylesheet, srv.Erro).ServeHTTP(w, r)
	})
}

// File 将 path 指向的内容作为文档内容生成中间件
func (srv *Server) File(path core.URI) (http.Handler, error) {
	data, err := path.ReadAll(nil)
	if err != nil {
		return nil, err
	}

	if srv.Path == "" {
		file, err := path.File()
		if err != nil {
			return nil, err
		}
		srv.Path = "/" + filepath.Base(file)
	}

	return srv.Buffer(data), nil
}

// 用于查找 <?xml 指令
var procInst = regexp.MustCompile(`<\?xml .+ ?>`)

func addStylesheet(data []byte) []byte {
	pi := `
<?xml-stylesheet type="text/xsl" href="` + docs.StylesheetURL("./") + `"?>`

	if rslt := procInst.Find(data); len(rslt) > 0 {
		return procInst.ReplaceAll(data, append(rslt, []byte(pi)...))
	}

	ret := make([]byte, 0, len(data)+len(pi))
	return append(append(ret, pi...), data...)
}
