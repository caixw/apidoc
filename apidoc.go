// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具。
package apidoc

import (
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5/doc"
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

// Version 获取当前程序的版本号
func Version() string {
	return vars.Version()
}

// Output 按 output 的要求输出内容。
func Output(doc *doc.Doc, opt *output.Options) error {
	return output.Render(doc, opt)
}

// Parse 分析从 input 中获取的代码块
//
// 所有与解析有关的错误均通过 h 输出。
// 如果 input 参数有误，会通过 error 参数返回。
func Parse(h *message.Handler, opt ...*input.Options) (*doc.Doc, error) {
	return input.Parse(h, opt...)
}
