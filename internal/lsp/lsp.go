// SPDX-License-Identifier: MIT

// Package lsp 提供 language server protocol 服务
package lsp

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Version lsp 的版本
const Version = "3.16.0"

// Serve 执行 LSP 服务
//
// t 表示服务的类型，可以是 stdio、udp、tcp 和 unix。
func Serve(header bool, t string, addr string, timeout time.Duration, infolog, errlog *log.Logger) error {
	switch strings.ToLower(t) {
	case "stdio":
		return serveStdio(header, infolog, errlog)
	case "udp":
		return serveUDP(header, addr, timeout, infolog, errlog)
	case "tcp", "unix":
		return serveTCP(header, t, addr, timeout, infolog, errlog)
	}

	return core.NewError(locale.ErrInvalidValue)
}

func serveStdio(header bool, infolog, errlog *log.Logger) error {
	return serve(jsonrpc.NewStreamTransport(header, os.Stdin, os.Stdout, nil), infolog, errlog)
}

func serveUDP(header bool, addr string, timeout time.Duration, infolog, errlog *log.Logger) error {
	t, err := jsonrpc.NewUDPServerTransport(header, addr, timeout)
	if err != nil {
		return err
	}
	return serve(t, infolog, errlog)
}

// t 可以是 tcp 和 unix
func serveTCP(header bool, t string, addr string, timeout time.Duration, infolog, errlog *log.Logger) error {
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
		return serve(jsonrpc.NewSocketTransport(header, conn, timeout), infolog, errlog)
	}
}

func serve(t jsonrpc.Transport, infolog, errlog *log.Logger) error {
	return newServe(t, infolog, errlog).serve()
}
