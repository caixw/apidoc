// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/issue9/jsonrpc"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/internal/locale"
)

// 错误代码，rpc 定义并不全，全部重新定义
const (
	ErrParseError           = jsonrpc.CodeParseError
	ErrInvalidRequest       = jsonrpc.CodeInvalidRequest
	ErrMethodNotFound       = jsonrpc.CodeMethodNotFound
	ErrInvalidParams        = jsonrpc.CodeInvalidParams
	ErrInternalError        = jsonrpc.CodeInternalError
	ErrServerErrorStart     = -32099
	ErrServerErrorEnd       = -32000
	ErrServerNotInitialized = -32002
	ErrUnknownErrorCode     = -32001

	// Defined by the protocol.
	ErrRequestCancelled = -32800
	ErrContentModified  = -32801
)

func newError(code int, key message.Reference, v ...interface{}) *jsonrpc.Error {
	return jsonrpc.NewError(code, locale.Sprintf(key, v...))
}
