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
func Conn(t jsonrpc.Transport, errlog *log.Logger) *jsonrpc.Conn {
	jsonrpcServer := jsonrpc.NewServer()
	conn := jsonrpcServer.NewConn(t, errlog)
	srv := newServer(conn)

	jsonrpcServer.Registers(map[string]interface{}{
		"initialize":  srv.initialize,
		"initialized": srv.initialized,
		"shutdown":    srv.shutdown,

		"workspace/didChangeWorkspaceFolders": srv.workspaceDidChangeWorkspaceFolders,
	})

	return conn
}
