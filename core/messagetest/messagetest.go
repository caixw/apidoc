// SPDX-License-Identifier: MIT

// Package messagetest 提供测试生成 message 相关的测试工具
package messagetest

import (
	"github.com/caixw/apidoc/v7/core"
)

// Result NewMessageHandler 返回的对象
type Result struct {
	Errors, Warns, Infos, Successes []interface{}
	Handler                         *core.MessageHandler
}

// NewMessageHandler 返回一个用于测试的 message.MessageHandler 实例
func NewMessageHandler() *Result {
	rslt := &Result{
		Errors:    []interface{}{},
		Warns:     []interface{}{},
		Infos:     []interface{}{},
		Successes: []interface{}{},
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
