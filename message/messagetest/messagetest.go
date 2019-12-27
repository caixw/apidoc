// SPDX-License-Identifier: MIT

// Package messagetest 提供测试生成 message 相关的测试工具
package messagetest

import (
	"bytes"

	"github.com/caixw/apidoc/v5/message"
)

// MessageHandler 返回一个用于测试的 message.Handler 实例
func MessageHandler() (erro, succ *bytes.Buffer, h *message.Handler) {
	erro = new(bytes.Buffer)
	succ = new(bytes.Buffer)

	f := func(msg *message.Message) {
		switch msg.Type {
		case message.Erro:
			erro.WriteString(msg.Message)
		default:
			succ.WriteString(msg.Message)
		}
	}

	return erro, succ, message.NewHandler(f)
}
