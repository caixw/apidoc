// SPDX-License-Identifier: MIT

// Package lsp 提供 language server protocol 服务
package lsp

import (
	"net/http"

	"github.com/gorilla/rpc/v2"
)

// Version lsp 的版本
const Version = "3.14.0"

// Server 执行 LSP 服务
func Server() (http.Handler, error) {
	srv := rpc.NewServer()
	srv.RegisterCodec(newJSONCodec(), "application/json")

	if err := srv.RegisterService(hello, "hello"); err != nil {
		return nil, err
	}
	return srv, nil
}

var hello = &HelloService{}

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}
