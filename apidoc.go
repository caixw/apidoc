// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具
//
// 可以从代码文件的注释中提取文档内容，生成 API 文档，
// 支持大部分的主流的编程语言。
//
// 在生成文档之前，请确保已经调用 Init() 用于初始化环境，
// Init() 可以确保能以你指定的本地化信息显示提示信息。
package apidoc

import (
	"bytes"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp"
	"github.com/caixw/apidoc/v7/internal/mock"
	"github.com/caixw/apidoc/v7/internal/vars"
)

// Config 配置文件映身的结构
type Config = build.Config

// SetLocale 设置当前的本地化 ID
//
// 如果不成功返回 false，比如设置了个不支持的本地化 ID。
func SetLocale(tag language.Tag) {
	locale.SetLanguageTag(tag)
}

// Locale 获取当前设置的本地化 ID
func Locale() language.Tag {
	return locale.LanguageTag()
}

// Locales 返回当前所有支持的本地信息
func Locales() map[language.Tag]string {
	return locale.DisplayNames()
}

// Version 当前程序的版本号
//
// 包含了版本号，编译日期以及编译是的 Git 记录 ID。
//
// 为一个正常的 semver(https://semver.org/lang/zh-CN/) 格式字符串。
func Version(full bool) string {
	if full {
		return vars.FullVersion()
	}
	return vars.Version
}

// LSPVersion 获取当前支持的 LSP 版本
func LSPVersion() string {
	return lsp.Version
}

// DocVersion 获取文档的版本信息
func DocVersion() string {
	return ast.Version
}

// Build 解析文档并输出文档内容
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 需要先调用 Init() 初始化本地化信息
func Build(h *core.MessageHandler, o *build.Output, i ...*build.Input) error {
	return build.Build(h, o, i...)
}

// Buffer 生成文档内容并返回
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 需要先调用 Init() 初始化本地化信息
func Buffer(h *core.MessageHandler, o *build.Output, i ...*build.Input) (*bytes.Buffer, error) {
	return build.Buffer(h, o, i...)
}

// Test 测试文档语法，并将结果输出到 h
//
// NOTE: 需要先调用 Init() 初始化本地化信息
func Test(h *core.MessageHandler, i ...*build.Input) {
	build.Test(h, i...)
}

// Static 为 /docs 搭建一个静态文件服务
//
// 相当于本地版本的 https://apidoc.tools，默认页为 index.xml。
//
// 用户可以通过诸如：
//  http.Handle("/apidoc", apidoc.Static(...))
// 的代码搭建一个简易的 https://apidoc.tools 网站。
//
// /docs 存放了整个项目的文档内容。其中根目录中包含网站的相关内容，
// 而 /v7 这些以版本号开头的则是查看 xml 文档的工具代码。
// 同时这一份代码也被编译在代码中。如果你不需要修改文档内容，
// 则可以直接传递空的 dir，表示采用内置的文档，否则指向指定的目录，
// 如果指向了自定义的目录，需要保证目录结构和文件名与 /docs 相同。
// stylesheet 则指定了是否需要根目录的内容，如果为 true，只会提供转换工具的代码。
func Static(dir core.URI, stylesheet bool) http.Handler {
	return docs.Handler(dir, stylesheet)
}

// View 返回查看文档的中间件
//
// 提供了与 Static 相同的功能，同时又可以额外添加一个文件。
// 与 Buffer 结合，可以提供一个完整的文档查看功能。
//
// status 是新文档的返回的状态码；
// url 表示文档在路由中的地址，必须以 / 开头；
// data 表示文档的实际内容，会添加 xml-stylesheet 指令，并指向当前的 apidoc.xsl；
// contentType 表示文档的 Content-Type 报头值；
// dir 和 stylesheet 则和 Static 相同。
func View(status int, url string, data []byte, contentType string, dir core.URI, stylesheet bool) http.Handler {
	data = addStylesheet(data)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == url {
			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(status)
			w.Write(data)
			return
		}

		Static(dir, stylesheet).ServeHTTP(w, r)
	})
}

// ServeLSP 提供 language server protocol 服务
//
// header 表示传递内容是否带报头；
// t 表示允许连接的类型，目前可以是 tcp、udp、stdio 和 ipc
func ServeLSP(header bool, t, addr string, infolog, errlog *log.Logger) error {
	return lsp.Serve(header, t, addr, infolog, errlog)
}

// ViewFile 返回查看文件的中间件
//
// 功能等同于 View，但是将 data 参数换成了文件地址。
// url 可以为空值，表示接受 path 的文件名部分作为其值。
//
// path 可以是远程文件，也可以是本地文件。
func ViewFile(status int, url string, path core.URI, contentType string, dir core.URI, stylesheet bool) (http.Handler, error) {
	data, err := path.ReadAll(nil)
	if err != nil {
		return nil, err
	}

	file, err := path.File()
	if err != nil {
		return nil, err
	}
	if url == "" {
		url = "/" + filepath.Base(file)
	}

	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(file))
	}

	return View(status, url, data, contentType, dir, stylesheet), nil
}

// Valid 验证文档内容的正确性
func Valid(b core.Block) error {
	return (&ast.APIDoc{}).Parse(b)
}

// Mock 生成 Mock 中间件
//
// data 为文档内容；
// servers 为文档中所有 server 以及对应的路由前缀。
func Mock(h *core.MessageHandler, data []byte, servers map[string]string) (http.Handler, error) {
	d := &ast.APIDoc{}
	if err := d.Parse(core.Block{Data: data}); err != nil {
		return nil, err
	}

	return mock.New(h, d, servers)
}

// MockFile 生成 Mock 中间件
//
// path 为文档路径，可以是本地路径也可以是 URL，根据是否为 http 或是 https 开头做判断；
// servers 为文档中所有 server 以及对应的路由前缀。
func MockFile(h *core.MessageHandler, path core.URI, servers map[string]string) (http.Handler, error) {
	return mock.Load(h, path, servers)
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
