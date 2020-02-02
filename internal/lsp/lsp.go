// SPDX-License-Identifier: MIT

// Package lsp 提供 language server protocol 服务
package lsp

import (
	"log"

	"github.com/issue9/jsonrpc"
)

// Version lsp 的版本
const Version = "3.15.0"

// Conn 创建 jsonrpc.Conn 实例
func Conn(errlog *log.Logger) *jsonrpc.Conn {
	conn := jsonrpc.NewConn(errlog)
	srv := newServer()
	ws := newWorkspace(srv)

	conn.Registers(map[string]interface{}{
		"initialize": srv.initialize,

		"workspace/workspaceFolders": ws.workspaceFolders,
	})

	return conn
}
