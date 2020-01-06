// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具
//
// 可以从代码文件的注释中提取文档内容，生成 API 文档，
// 支持大部分的主流的编程语言。
//
// 在生成文档之前，请确保已经调用 Init() 用于初始化环境，
// Init() 可以确保能以你指定的本地化信息显示提示信息。
//
// 生成的文档，可以调用 Do() 输出为文件；也可以通过 Buffer()
// 返回 bytes.Buffer 实例；或者通过 Pack() 直接将文档与其依赖的 XSL
// 打包成一个 Go 源码文件，这样可以直接编译在二进制文件中。
package apidoc

import (
	"bytes"
	"net/http"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/mock"
	"github.com/caixw/apidoc/v5/internal/static"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
	"github.com/caixw/apidoc/v5/output"
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

// Do 解析文档并输出文档内容
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 需要先调用 Init() 初始化本地化信息
func Do(h *message.Handler, o *output.Options, i ...*input.Options) error {
	doc, err := input.Parse(h, i...)
	if err != nil {
		return err
	}

	return output.Render(doc, o)
}

// Buffer 生成文档内容并返回
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 需要先调用 Init() 初始化本地化信息
func Buffer(h *message.Handler, o *output.Options, i ...*input.Options) (*bytes.Buffer, error) {
	doc, err := input.Parse(h, i...)
	if err != nil {
		return nil, err
	}

	return output.Buffer(doc, o)
}

// Test 测试文档语法，并将结果输出到 h
func Test(h *message.Handler, i ...*input.Options) {
	if _, err := input.Parse(h, i...); err != nil {
		h.Error(message.Erro, err)
		return
	}
	h.Message(message.Succ, locale.TestSuccess)
}

// Static 返回文件服务中间件
//
// 相当于本地版本的 https://apidoc.tools，默认页为 index.xml。
//
// 用户可以通过诸如：
//  http.Handle("/apidoc", apidoc.Static(...))
// 的代码搭建一个简易的 https://apidoc.tools 网站。
//
// dir 表示文档的根目录，会在该目录下查找用户请求的文档内容。
// 当 dir 为空时，表示使用内置的静态文件作为文件服务的内容 ，
// 如果是普通的文件静态服务，可以直接采用 http.FileServer 会更通用；
// t 表示可以访问的文件类型。
//
// NOTE: 只要 dir 不为空，则只会采用 dir 文件夹的内容作为文件服务的主体内容。
// dir 是否为空的区别在于：dir 指同一个本地目录，方便在运行时进行修改；
// 而 dir 为空则直接将 /docs 内容内嵌到代码中，无法修改。
func Static(dir string, t static.Type) http.Handler {
	if dir == "" {
		return static.EmbeddedHandler(t)
	}

	return static.FolderHandler(dir, t)
}

// Valid 验证文档内容的正确性
func Valid(content []byte) error {
	return doc.Valid(content)
}

// Mock 生成 Mock 中间件
//
// data 为文档内容；
// servers 为文档中所有 server 以及对应的路由前缀。
func Mock(h *message.Handler, data []byte, servers map[string]string) (http.Handler, error) {
	d := doc.New()
	if err := d.FromXML("", 0, data); err != nil {
		return nil, err
	}

	return mock.New(h, d, servers)
}

// MockFile 生成 Mock 中间件
//
// path 为文档路径，可以是本地路径也可以是 URL，根据是否为 http 或是 https 开头做判断；
// servers 为文档中所有 server 以及对应的路由前缀。
func MockFile(h *message.Handler, path string, servers map[string]string) (http.Handler, error) {
	return mock.Load(h, path, servers)
}
