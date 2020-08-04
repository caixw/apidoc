// SPDX-License-Identifier: MIT

package core

import (
	"github.com/issue9/sliceutil"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/internal/locale"
)

// HTTPError 表示 HTTP 状态码的错误
type HTTPError struct {
	locale.Err
	Code int
}

// Error 用于描述 apidoc 中的大部分错误信息
//
// 无论是配置文件的错误，还是文档的语法错误，都将返回此错误。
type Error struct {
	Err      error       // 具体的错误信息
	Location Location    // 错误的详细定位
	Field    string      // 出错的字段
	Types    []ErrorType // 该错误的类型
}

// ErrorType 语法错误的类型
type ErrorType int

// 语法错误类型可用的枚举值
const (
	ErrorTypeDeprecated ErrorType = iota + 1
	ErrorTypeUnused
)

// NewHTTPError 声明 HTTPError 实例
func NewHTTPError(code int, key message.Reference, v ...interface{}) *HTTPError {
	return &HTTPError{
		Err:  locale.Err{Key: key, Values: v},
		Code: code,
	}
}

func (err *Error) Error() string {
	detail := err.Location.String()

	if err.Field != "" {
		detail += ":" + err.Field
	}

	// ErrMessage = "%s 位次于 %s:%d"
	return locale.Sprintf(locale.ErrMessage, err.Err.Error(), detail)
}

// Unwrap 实现 errors.Unwrap 接口
func (err *Error) Unwrap() error {
	return err.Err
}

// Is 实现 errors.Is 接口
func (err *Error) Is(target error) bool {
	return err.Err == target
}

// WithField 为语法错误修改或添加具体的错误字段
func (err *Error) WithField(field string) *Error {
	err.Field = field
	return err
}

// WithLocation  为语法错误添加定位信息
func (err *Error) WithLocation(loc Location) *Error {
	err.Location = loc
	return err
}

// AddTypes 为语法错误添加错误类型
func (err *Error) AddTypes(t ...ErrorType) *Error {
	if err.Types == nil {
		err.Types = t
		return err
	}

	for _, typ := range t {
		if sliceutil.Count(err.Types, func(i int) bool { return err.Types[i] == typ }) <= 0 {
			err.Types = append(err.Types, typ)
		}
	}

	return err
}

// NewError 返回 *Error 实例
func NewError(key message.Reference, v ...interface{}) *Error {
	return &Error{Err: locale.NewError(key, v...)}
}

// WithError 采用 err 实例 *Error 实例
func WithError(err error) *Error {
	var types []ErrorType

	if serr, ok := err.(*Error); ok {
		err = serr.Err
		types = serr.Types
	}

	return (&Error{Err: err}).AddTypes(types...)
}

// NewError 在当前位置生成语法错误信息
//
// 其中的 msg 和 val 会被转换成本地化的内容保存。
func (loc Location) NewError(key message.Reference, v ...interface{}) *Error {
	return &Error{Err: locale.NewError(key, v...), Location: loc}
}

// WithError 在当前位置生成语法错误信息
//
// 若 err 本身就是 *Error 类型，则会更新其 location 和 Field 两个字段的信息。
func (loc Location) WithError(err error) *Error {
	return &Error{Err: err, Location: loc}
}
