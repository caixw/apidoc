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
func Serve(t string, addr string, infolog, errlog *log.Logger) error {
	switch strings.ToLower(t) {
	case "pipe":
	case "stdio":
		return serveStdio(infolog, errlog)
	case "ipc":
		return serveStdio(infolog, errlog)
	case "udp":
		return serveUDP(addr, infolog, errlog)
	case "tcp", "unix":
		return serveTCP(t, addr, infolog, errlog)
	}

	return message.NewLocaleError("", "", 0, locale.ErrInvalidValue)
}

func serveStdio(infolog, errlog *log.Logger) error {
	return serve(jsonrpc.NewStreamTransport(os.Stdin, os.Stdout), infolog, errlog)
}

func serveUDP(addr string, infolog, errlog *log.Logger) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	return serve(jsonrpc.NewSocketTransport(conn), infolog, errlog)
}

// t 可以是 tcp 和 unix
func serveTCP(t string, addr string, infolog, errlog *log.Logger) error {
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
		return serve(jsonrpc.NewSocketTransport(conn), infolog, errlog)
	}
}

func serve(t jsonrpc.Transport, infolog, errlog *log.Logger) error {
	jsonrpcServer := jsonrpc.NewServer()

	jsonrpcServer.RegisterBefore(func(method string) error {
		infolog.Println(locale.Sprintf(locale.RequestRPC, method))
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())

	srv := &server{
		Conn:       jsonrpcServer.NewConn(t, errlog),
		state:      serverCreated,
		cancelFunc: cancel,
	}

	jsonrpcServer.Registers(map[string]interface{}{
		"initialize":  srv.initialize,
		"initialized": srv.initialized,
		"shutdown":    srv.shutdown,

		"workspace/didChangeWorkspaceFolders": srv.workspaceDidChangeWorkspaceFolders,
	})

	return srv.Serve(ctx)
}
