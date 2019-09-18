// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具。
package apidoc

import (
	"context"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
	"github.com/caixw/apidoc/v5/output"
)

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
func Parse(ctx context.Context, h *message.Handler, opt ...*input.Options) (*doc.Doc, error) {
	return input.Parse(ctx, h, opt...)
}
