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
	"mime"
	"net/http"
	"path/filepath"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v6/doc"
	"github.com/caixw/apidoc/v6/input"
	"github.com/caixw/apidoc/v6/internal/docs"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/mock"
	xpath "github.com/caixw/apidoc/v6/internal/path"
	"github.com/caixw/apidoc/v6/internal/vars"
	"github.com/caixw/apidoc/v6/message"
	"github.com/caixw/apidoc/v6/output"
)

// Init 初始化包
//
// 如果传递了 language.Und，则采用系统当前的本地化信息。
// 如果获取系统的本地化信息依然失败，则会失放 zh-Hans 作为默认值。
func Init(tag language.Tag) error {
	return locale.Init(tag)
}

// Version 当前程序的版本号
//
// 为一个正常的 semver(https://semver.org/lang/zh-CN/) 格式字符串。
func Version() string {
	return vars.Version()
}

// Build 解析文档并输出文档内容
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 需要先调用 Init() 初始化本地化信息
func Build(h *message.Handler, o *output.Options, i ...*input.Options) error {
	d, err := input.Parse(h, i...)
	if err != nil {
		return err
	}

	return output.Render(d, o)
}

// Buffer 生成文档内容并返回
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 需要先调用 Init() 初始化本地化信息
func Buffer(h *message.Handler, o *output.Options, i ...*input.Options) (*bytes.Buffer, error) {
	d, err := input.Parse(h, i...)
	if err != nil {
		return nil, err
	}

	return output.Buffer(d, o)
}

// Test 测试文档语法，并将结果输出到 h
func Test(h *message.Handler, i ...*input.Options) {
	if _, err := input.Parse(h, i...); err != nil {
		h.Error(message.Erro, err)
		return
	}
	h.Message(message.Succ, locale.TestSuccess)
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
// 而 /v6 这些以版本号开头的则是查看 xml 文档的工具代码。
// 同时这一份代码也被编译在代码中。如果你不需要修改文档内容，
// 则可以直接传递空的 dir，表示采用内置的文档，否则指向指定的目录，
// 如果指向了自定义的目录，需要保证目录结构和文件名与 /docs 相同。
// stylesheet 则指定了是否需要根目录的内容，如果为 true，只会提供转换工具的代码。
func Static(dir string, stylesheet bool) http.Handler {
	return docs.Handler(dir, stylesheet)
}

// View 返回查看文档的中间件
//
// 提供了与 Static 相同的功能，同时又可以额外添加一个文件。
// 与 Buffer 结合，可以提供一个完整的文档查看功能。
//
// status 是新文档的返回的状态码；
// url 表示文档在路由中的地址，必须以 / 开头；
// data 表示文档的实际内容；
// contentType 表示文档的 Content-Type 报头值；
// dir 和 stylesheet 则和 Static 相同。
func View(status int, url string, data []byte, contentType, dir string, stylesheet bool) http.Handler {
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

// ViewFile 返回查看文件的中间件
//
// 功能等同于 View，但是将 data 参数换成了文件地址。
// url 可以为空值，表示接受 path 的文件名部分作为其值。
//
// path 可以是远程文件 (http 开头)，也可以是本地文件。
func ViewFile(status int, url, path, contentType, dir string, stylesheet bool) (http.Handler, error) {
	data, err := xpath.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if url == "" {
		url = "/" + filepath.Base(path)
	}

	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(path))
	}

	return View(status, url, data, contentType, dir, stylesheet), nil
}

// Valid 验证文档内容的正确性
func Valid(content []byte) error {
	return doc.Valid(content)
}

// Mock 生成 Mock 中间件
//
// 调用者需要保证 d 的正确性。
func Mock(h *message.Handler, d *doc.Doc, servers map[string]string) (http.Handler, error) {
	return mock.New(h, d, servers)
}

// MockBuffer 生成 Mock 中间件
//
// data 为文档内容；
// servers 为文档中所有 server 以及对应的路由前缀。
func MockBuffer(h *message.Handler, data []byte, servers map[string]string) (http.Handler, error) {
	d := doc.New()
	if err := d.FromXML("", 0, data); err != nil {
		return nil, err
	}

	return Mock(h, d, servers)
}

// MockFile 生成 Mock 中间件
//
// path 为文档路径，可以是本地路径也可以是 URL，根据是否为 http 或是 https 开头做判断；
// servers 为文档中所有 server 以及对应的路由前缀。
func MockFile(h *message.Handler, path string, servers map[string]string) (http.Handler, error) {
	return mock.Load(h, path, servers)
}
