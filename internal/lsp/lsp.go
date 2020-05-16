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

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Version lsp 的版本
const Version = "3.15.0"

// Serve 执行 LSP 服务
func Serve(header bool, t string, addr string, infolog, errlog *log.Logger) error {
	switch strings.ToLower(t) {
	case "pipe":
	case "stdio":
		return serveStdio(header, infolog, errlog)
	case "ipc":
		return serveStdio(header, infolog, errlog)
	case "udp":
		return serveUDP(header, addr, infolog, errlog)
	case "tcp", "unix":
		return serveTCP(header, t, addr, infolog, errlog)
	}

	return core.NewSyntaxError(core.Location{}, "", 0, locale.ErrInvalidValue)
}

func serveStdio(header bool, infolog, errlog *log.Logger) error {
	return serve(jsonrpc.NewStreamTransport(header, os.Stdin, os.Stdout, nil), infolog, errlog)
}

func serveUDP(header bool, addr string, infolog, errlog *log.Logger) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	return serve(jsonrpc.NewSocketTransport(header, conn), infolog, errlog)
}

// t 可以是 tcp 和 unix
func serveTCP(header bool, t string, addr string, infolog, errlog *log.Logger) error {
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
		return serve(jsonrpc.NewSocketTransport(header, conn), infolog, errlog)
	}
}

func serve(t jsonrpc.Transport, infolog, errlog *log.Logger) error {
	jsonrpcServer := jsonrpc.NewServer()

	jsonrpcServer.RegisterBefore(func(method string) error {
		if strings.HasPrefix(method, "$/") && !jsonrpcServer.Exists(method) {
			return newError(ErrMethodNotFound, locale.UnimplementedRPC, method)
		}

		log.Println(locale.Sprintf(locale.RequestRPC, method))
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
		"initialize":      srv.initialize,
		"initialized":     srv.initialized,
		"shutdown":        srv.shutdown,
		"exit":            srv.exit,
		"$/cancelRequest": srv.cancel,

		// workspace
		"workspace/didChangeWorkspaceFolders": srv.workspaceDidChangeWorkspaceFolders,

		// textDocument
		"textDocument/didOpen":   srv.textDocumentDidOpen,
		"textDocument/didChange": srv.textDocumentDidChange,
		"textDocument/hover":     srv.textDocumentHover,
	})

	return srv.Serve(ctx)
}

// 分析 path 的内容，并将其中的文档解析至 doc
func parseFile(doc *ast.APIDoc, h *core.MessageHandler, uri core.URI, i *build.Input) {
	doc.ParseBlocks(h, func(blocks chan core.Block) {
		i.ParseFile(blocks, h, uri)
	})
}
