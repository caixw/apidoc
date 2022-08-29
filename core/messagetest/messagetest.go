// SPDX-License-Identifier: MIT

// Package messagetest 提供测试生成 message 相关的测试工具
package messagetest

import "github.com/caixw/apidoc/v7/core"

type Result struct {
	Errors, Warns, Infos, Successes []any
	Handler                         *core.MessageHandler
}

// NewMessageHandler 返回一个用于测试的 core.MessageHandler 实例
func NewMessageHandler() *Result {
	rslt := &Result{
		Errors:    []any{},
		Warns:     []any{},
		Infos:     []any{},
		Successes: []any{},
	}

	rslt.Handler = core.NewMessageHandler(func(msg *core.Message) {
		switch msg.Type {
		case core.Erro:
			rslt.Errors = append(rslt.Errors, msg.Message)
		case core.Warn:
			rslt.Warns = append(rslt.Warns, msg.Message)
		case core.Info:
			rslt.Infos = append(rslt.Infos, msg.Message)
		case core.Succ:
			rslt.Successes = append(rslt.Successes, msg.Message)
		default:
			panic("unreached")
		}
	})

	return rslt
}
