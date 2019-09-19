// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具。
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

// Version 获取当前程序的版本号
func Version() string {
	return vars.Version()
}

// Do 执行分析操作
//
// 需要确保已经调用 o 和 i 的各个 Sanitize 方法。
func Do(h *message.Handler, o *output.Options, i ...*input.Options) error {
	return output.Render(input.Parse(h, i...), o)
}
