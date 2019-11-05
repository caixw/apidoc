// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具
//
// 可以从代码文件的注释中提取文档内容，生成 API 文档，
// 支持大部分的主流的编程语言。
package apidoc

import (
	"bytes"
	"net/http"
	"time"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/config"
	"github.com/caixw/apidoc/v5/internal/docs"
	"github.com/caixw/apidoc/v5/internal/locale"
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
// 如果需要控制详细的操作步骤，可以自行调用 input 和 output 的相关函数实现。
//
// 如果是文档语法错误，则相关的错误信息会反馈给 h，由 h 处理错误信息；
// 如果是配置项（o 和 i）有问题，则以 *message.SyntaxError 类型返回错误信息。
//
// NOTE: 需要先调用 Init() 初始化本地化信息。
func Do(h *message.Handler, o *output.Options, i ...*input.Options) error {
	doc, err := input.Parse(h, i...)
	if err != nil {
		return err
	}

	return output.Render(doc, o)
}

// Buffer 生成文档内容，并保存在内存。
func Buffer(h *message.Handler, o *output.Options, i ...*input.Options) (*bytes.Buffer, error) {
	doc, err := input.Parse(h, i...)
	if err != nil {
		return nil, err
	}

	return output.Buffer(doc, o)
}

// Make 根据 wd 目录下的配置文件生成文档
//
// test 表示是否仅测试该配置文件和文档的正确性。
//
// NOTE: 需要先调用 Init() 初始化本地化信息。
func Make(h *message.Handler, wd string, test bool) {
	now := time.Now()

	cfg, err := config.Load(wd)
	if err != nil {
		h.Error(message.Erro, err)
		return
	}

	if test {
		if _, err := input.Parse(h, cfg.Inputs...); err != nil {
			h.Error(message.Erro, err)
			return
		}
		h.Message(message.Succ, locale.TestSuccess)
	} else {
		if err := Do(h, cfg.Output, cfg.Inputs...); err != nil {
			h.Error(message.Erro, err)
			return
		}
		h.Message(message.Succ, locale.Complete, cfg.Output.Path, time.Now().Sub(now))
	}
}

// MakeBuffer 根据 wd 目录下的配置文件生成文档内容并保存至内存
//
// NOTE: 需要先调用 Init() 初始化本地化信息。
func MakeBuffer(h *message.Handler, wd string) (*bytes.Buffer, time.Duration) {
	now := time.Now()

	cfg, err := config.Load(wd)
	if err != nil {
		h.Error(message.Erro, err)
		return nil, 0
	}

	buf, err := Buffer(h, cfg.Output, cfg.Inputs...)
	if err != nil {
		h.Error(message.Erro, err)
		return nil, 0
	}

	return buf, time.Now().Sub(now)
}

// Detect 检测 wd 目录，并根据目录情况写入配置文件到该目录。
func Detect(wd string, recursive bool) error {
	return config.Write(wd, recursive)
}

// Site 将 dir 作为静态文件服务内容
//
// 默认页为 index.xml，同时会过滤 CNAME，
// 如果将 dir 指同 docs 目录，相当于本地版本的 https://apidoc.tools
//
// 用户可以通过诸如：
//  http.Handle("/apidoc", apidoc.Site("./docs"))
// 的代码搭建一个简易的 https://apidoc.tools 网站。
func Site(dir string) http.Handler {
	return docs.Handler(dir)
}
