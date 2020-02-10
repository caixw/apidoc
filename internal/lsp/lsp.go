// SPDX-License-Identifier: MIT

// Package lsp 提供 language server protocol 服务
package lsp

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
)

// Version lsp 的版本
const Version = "3.15.0"

// Serve 执行 LSP 服务
func Serve(t string, addr string, errlog *log.Logger) error {
	switch strings.ToLower(t) {
	case "pipe":
	case "stdio":
		return serveStdio(errlog)
	case "ipc":
		return serveStdio(errlog)
	case "udp":
		return serveUDP(addr, errlog)
	case "tcp", "unix":
		return serveTCP(t, addr, errlog)
	}

	return message.NewLocaleError("", "", 0, locale.ErrInvalidValue)
}

func serveStdio(errlog *log.Logger) error {
	jsonrpcServer := jsonrpc.NewServer()
	return serve(jsonrpcServer, jsonrpc.NewStreamTransport(os.Stdin, os.Stdout), errlog)
}

func serveUDP(addr string, errlog *log.Logger) error {
	jsonrpcServer := jsonrpc.NewServer()
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	return serve(jsonrpcServer, jsonrpc.NewSocketTransport(conn), errlog)
}

// t 可以是 tcp 和 unix
func serveTCP(t string, addr string, errlog *log.Logger) error {
	jsonrpcServer := jsonrpc.NewServer()

	l, err := net.Listen(t, addr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			errlog.Println(err)
			continue
		}
		return serve(jsonrpcServer, jsonrpc.NewSocketTransport(conn), errlog)
	}
}

func serve(jsonrpcServer *jsonrpc.Server, t jsonrpc.Transport, errlog *log.Logger) error {
	conn := jsonrpcServer.NewConn(t, errlog)
	ctx, cancel := context.WithCancel(context.Background())
	srv := newServer(conn, cancel)

	jsonrpcServer.Registers(map[string]interface{}{
		"initialize":  srv.initialize,
		"initialized": srv.initialized,
		"shutdown":    srv.shutdown,

		"workspace/didChangeWorkspaceFolders": srv.workspaceDidChangeWorkspaceFolders,
	})

	return conn.Serve(ctx)
}
