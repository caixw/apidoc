// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具
//
// 从代码文件的注释中提取特定格式的内容，生成 API 文档，
// 支持大部分的主流的编程语言。
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
)

const (
	// LSPVersion 当前支持的 language server protocol 版本
	LSPVersion = lsp.Version

	// DocVersion 文档的版本
	DocVersion = ast.Version
)

type (
	// Config 生成文档时所需要的配置项
	Config = build.Config

	// PackOptions 指定了打包文档内容的参数
	PackOptions = build.PackOptions
)

// SetLocale 设置当前的本地化 ID
func SetLocale(tag language.Tag) {
	locale.SetTag(tag)
}

// Locale 获取当前设置的本地化 ID
func Locale() language.Tag {
	return locale.Tag()
}

// Locales 返回当前所有支持的本地化信息
func Locales() []language.Tag {
	return locale.Tags()
}

// Version 当前程序的版本号
//
// full 表示是否需要在版本号中包含编译日期和编译时的 Git 记录 ID。
func Version(full bool) string {
	if full {
		return core.FullVersion()
	}
	return core.Version
}

// Build 解析文档并输出文档内容
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 如果需要从配置文件进行构建文档，可以采用 Config.Build
func Build(h *core.MessageHandler, o *build.Output, i ...*build.Input) error {
	return build.Build(h, o, i...)
}

// Buffer 生成文档内容并返回
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 如果需要从配置文件进行构建文档，可以采用 Config.Buffer
func Buffer(h *core.MessageHandler, o *build.Output, i ...*build.Input) (*bytes.Buffer, error) {
	return build.Buffer(h, o, i...)
}

// CheckSyntax 测试文档语法
func CheckSyntax(h *core.MessageHandler, i ...*build.Input) {
	build.CheckSyntax(h, i...)
}

// Pack 将文档内容打包成一个 Go 文件
//
// opt 用于指定打包的设置项，如果为空，则会使用一个默认的设置项，
// 该默认设置项会在当前目录下创建一个包为 apidoc 的包，且公开文档数据为 APIDOC，
// 用户可以使用 Unpack 解包该常量的内容，即为一个合法的 apidoc 文档。
func Pack(h *core.MessageHandler, opt *PackOptions, o *build.Output, i ...*build.Input) error {
	return build.Pack(h, opt, o, i...)
}

// Unpack 用于解压由 Pack 输出的内容
func Unpack(buffer string) (string, error) {
	return build.Unpack(buffer)
}

// ServeLSP 提供 language server protocol 服务
//
// header 表示传递内容是否带报头；
// t 表示允许连接的类型，目前可以是 tcp、udp、stdio 和 ipc
func ServeLSP(header bool, t, addr string, infolog, errlog *log.Logger) error {
	return lsp.Serve(header, t, addr, infolog, errlog)
}

// Static 为 /docs 搭建一个静态文件服务
//
// 相当于本地版本的 https://apidoc.tools，默认页为 index.xml。
//
// 用户可以通过以下代码搭建一个简易的 https://apidoc.tools 网站：
//  http.Handle("/apidoc", apidoc.Static(...))
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
// url 表示文档在路由中的地址；
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

// ViewPack 返回查看文档的中间件
//
// 功能基本与 View 相同，但是第三个参数 unpackData 为 Pack() 函数打包之内的内容，
// 不需要调用 Unpack() 解包，直接由 ViewPack 自行解包。
func ViewPack(status int, url string, unpackData string, contentType string, dir core.URI, stylesheet bool) http.Handler {
	data, err := Unpack(unpackData)
	if err != nil {
		panic(err)
	}

	return View(status, url, []byte(data), contentType, dir, stylesheet)
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
