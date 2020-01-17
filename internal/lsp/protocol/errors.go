// SPDX-License-Identifier: MIT

package protocol

import "github.com/gorilla/rpc/v2/json2"

// 错误代码，rpc 定义并不全，全部重新定义
const (
	ErrParseError                           = json2.E_PARSE
	ErrInvalidRequest                       = json2.E_INVALID_REQ
	ErrMethodNotFound                       = json2.E_NO_METHOD
	ErrInvalidParams                        = json2.E_BAD_PARAMS
	ErrInternalError                        = json2.E_INTERNAL
	ErrServerErrorStart     json2.ErrorCode = -32099
	ErrServerErrorEnd       json2.ErrorCode = -32000
	ErrServerNotInitialized json2.ErrorCode = -32002
	ErrUnknownErrorCode     json2.ErrorCode = -32001

	// Defined by the protocol.
	ErrRequestCancelled json2.ErrorCode = -32800
	ErrContentModified  json2.ErrorCode = -32801
)

// NewError 声明 LSP 错误信息
func NewError(code json2.ErrorCode, message string, data []byte) error {
	return &json2.Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
