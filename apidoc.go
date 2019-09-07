// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具。
package apidoc

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5/errors"
	"github.com/caixw/apidoc/v5/internal/locale"
	o "github.com/caixw/apidoc/v5/internal/output"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/options"
)

// InitLocale 初始化语言环境
//
// NOTE: 必须保证在第一时间调用。
//
// 如果 tag 的值为 language.Und，则表示采用系统语言
func InitLocale(tag language.Tag) error {
	return locale.Init(tag)
}

// Version 获取当前程序的版本号
func Version() string {
	return vars.Version()
}

// Do 分析输入信息，并最终输出到指定的文件。
//
// h 表示处理语法错误的处理器。
// output 输出设置项；
// inputs 输入设置项。
func Do(ctx context.Context, h *errors.Handler, output *options.Output, inputs ...*options.Input) error {
	doc, err := Parse(ctx, h, inputs...)
	if err != nil {
		return err
	}

	// TODO doc.check

	return o.Render(doc, output)
}
