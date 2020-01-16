// SPDX-License-Identifier: MIT

// Package lsp 提供 language server protocol 服务
package lsp

import (
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// Version lsp 的版本
const Version = "3.14.0"

// Serve 执行 LSP 服务
func Serve() (http.Handler, error) {
	return NewServer()
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
	return &json2.Error{
		Code:    json2.E_INVALID_REQ,
		Message: "uninitialized TODO locale",
	}
}
