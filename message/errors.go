// SPDX-License-Identifier: MIT

package message

import (
	"encoding/xml"

	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// SyntaxError 表示语法错误信息
//
// 无论是配置文件的错误，还是文档的语法错误，都将返回此错误。
// apidoc 的错误基本上都是基于 SyntaxError。
type SyntaxError struct {
	Message string
	File    string
	Line    int
	Field   string
}

func (err *SyntaxError) Error() string {
	// TODO 根据是否有 filed，返回不同的提示内容

	// ErrMessage = "%s 位次于 %s:%d 的 %s"
	return locale.Sprintf(locale.ErrMessage, err.Message, err.File, err.Line, err.Field)
}

// NewError 声明新的 SyntaxError 实例
func NewError(file, field string, line int, msg message.Reference, val ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message: locale.Sprintf(msg, val...),
		File:    file,
		Line:    line,
		Field:   field,
	}
}

// WithError 声明 SyntaxError 实例，其中的提示信息由 err 返回
func WithError(file, field string, line int, err error) *SyntaxError {
	if serr, ok := err.(*xml.SyntaxError); ok {
		return &SyntaxError{
			Message: serr.Msg,
			File:    file,
			Line:    line + serr.Line,
			Field:   field,
		}
	}

	return &SyntaxError{
		Message: err.Error(),
		File:    file,
		Line:    line,
		Field:   field,
	}
}
