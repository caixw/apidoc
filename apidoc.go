// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具
//
// 可以从代码文件的注释中提取文档内容，生成 API 文档，
// 支持大部分的主流的编程语言。
//
// apidoc 采用了多协程处理各个文件，所有的语法错误都是以异步的方式发送给
// message.Handler 进行处理的。用户需要自行实现 message.HandlerFunc
// 类型的方法交给 message.Handler，以实现自已的消息处理功能。
package apidoc

import (
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5/input"
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

// Do 执行分析操作
//
// 如果是 o 和 i 的配置内容有问题，error 返回的实际类型为 *message.SyntaxError
func Do(h *message.Handler, o *output.Options, i ...*input.Options) error {
	doc, err := input.Parse(h, i...)
	if err != nil {
		return err
	}

	return output.Render(doc, o)
}

// Test 测试语法的正确性
//
// 错误信息依然输出到 h，配置文件的错误则直接返回。
func Test(h *message.Handler, i ...*input.Options) error {
	_, err := input.Parse(h, i...)
	return err
}
