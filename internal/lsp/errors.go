// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/issue9/jsonrpc"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/internal/locale"
)

// 错误代码，部分为 lsp 特有
const (
	ErrParseError           = jsonrpc.CodeParseError
	ErrInvalidRequest       = jsonrpc.CodeInvalidRequest
	ErrMethodNotFound       = jsonrpc.CodeMethodNotFound
	ErrInvalidParams        = jsonrpc.CodeInvalidParams
	ErrInternalError        = jsonrpc.CodeInternalError
	ErrServerNotInitialized = -32002
	ErrUnknownErrorCode     = -32001
	ErrRequestCancelled     = -32800
	ErrContentModified      = -32801
)

func newError(code int, key message.Reference, v ...any) *jsonrpc.Error {
	return jsonrpc.NewError(code, locale.Sprintf(key, v...))
}
