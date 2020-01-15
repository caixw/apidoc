// SPDX-License-Identifier: MIT

package lsp

import "github.com/gorilla/rpc/v2/json2"

// LSP 专有错误
const (
	RequestCancelled json2.ErrorCode = -32800
	ContentModified                  = -32801
)
