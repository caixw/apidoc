// SPDX-License-Identifier: MIT

// Package jsonrpc 实现简单的 json-rpc 2.0
//
// https://wiki.geekdream.com/Specification/json-rpc_2.0.html
package jsonrpc

import (
	"encoding/json"
)

// Version json-rpc 的版本
const Version = "2.0"

// LSP 定义的错误代码
const (
	// Defined by JSON RPC
	CodeParseError           = -32700
	CodeInvalidRequest       = -32600
	CodeMethodNotFound       = -32601
	CodeInvalidParams        = -32602
	CodeInternalError        = -32603
	CodeServerErrorStart     = -32099
	CodeServerErrorEnd       = -32000
	CodeServerNotInitialized = -32002
	CodeUnknownErrorCode     = -32001

	// Defined by the protocol.
	CodeRequestCancelled = -32800
	CodeContentModified  = -32801
)

// Request 请求对象
type Request struct {
	// 指定 JSON-RPC 协议版本的字符串
	Version string `json:"jsonrpc"`

	// 已建立客户端的唯一标识 id，值必须包含一个字符串、数值或 NULL 空值。
	// 如果不包含该成员则被认定为是一个通知。该值一般不为 NULL，若为数值则不应该包含小数。
	ID string `json:"id,omitempty"`

	// 包含所要调用方法名称的字符串
	//
	// 以 rpc 开头的方法名，用英文句号（U+002E or ASCII 46）
	// 连接的为预留给 rpc 内部的方法名及扩展名，且不能在其他地方使用。
	Method string `json:"method"`

	// 调用方法所需要的结构化参数值，该成员参数可以被省略。
	Params *json.RawMessage `json:"params,omitempty"`
}

// Response 响应对象
type Response struct {
	// 指定 JSON-RPC 协议版本的字符串
	Version string `json:"jsonrpc"`

	// 成功时的返回结果，如果不成功，则不应该输出该对象。
	Result *json.RawMessage `json:"result,omitempty"`

	// 失败时的返回结果，如果成功，则不应该输出该对象。
	Error *Error `json:"error,omitempty"`

	// ID 返回请求端的 ID，如果检查 ID 失败时，返回空值
	ID string `json:"id"`
}

// Error 错误类型
type Error struct {
	// 错误代码
	Code int `json:"code"`

	// 错误的简短描述
	Message string `json:"message"`

	// 详细的错误描述信息，可以为空
	Data interface{} `json:"data,omitempty"`
}

// NewError 新的 Error 对象
func NewError(code int, msg string) *Error {
	return NewErrorWithData(code, msg, nil)
}

// NewErrorWithData 新的 Error 对象
func NewErrorWithData(code int, msg string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func (err *Error) Error() string {
	return err.Message
}
