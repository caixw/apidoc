// SPDX-License-Identifier: MIT

package core

import "golang.org/x/text/message"

// SyntaxError 语法错误
//
// Deprecated: 请使用 Error 代替
type SyntaxError = Error

// NewSyntaxError 本地化的错误信息
//
// 其中的 msg 和 v 会被转换成本地化的内容保存。
//
// Deprecated: 请使用 NewError 代替
func NewSyntaxError(loc Location, field string, key message.Reference, v ...interface{}) *Error {
	return NewError(key, v...).WithLocation(loc).WithField(field)
}

// NewSyntaxErrorWithError 声明 SyntaxError 实例，其中的提示信息由 err 返回
//
// 若 err 本身就是 *SyntaxError 类型，则会更新其 location 和 Field 两个字段的信息。
//
// Deprecated: 请使用 WithError 代替
func NewSyntaxErrorWithError(loc Location, field string, err error) *Error {
	return WithError(err).WithLocation(loc).WithField(field)
}
