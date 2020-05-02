// SPDX-License-Identifier: MIT

// Package messagetest 提供测试生成 message 相关的测试工具
package messagetest

import (
	"bytes"
	"fmt"

	"github.com/caixw/apidoc/v7/core"
)

// MessageHandler 返回一个用于测试的 message.MessageHandler 实例
func MessageHandler() (erro, succ *bytes.Buffer, h *core.MessageHandler) {
	erro = new(bytes.Buffer)
	succ = new(bytes.Buffer)

	f := func(msg *core.Message) {
		switch msg.Type {
		case core.Erro:
			erro.WriteString(fmt.Sprint(msg.Message))
		default:
			succ.WriteString(fmt.Sprint(msg.Message))
		}
	}

	return erro, succ, core.NewMessageHandler(f)
}
