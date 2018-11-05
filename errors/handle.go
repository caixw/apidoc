// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package errors

import "log"

// HandlerFunc 错误处理函数
type HandlerFunc func(err *Error)

// Handler 用于接收错误信息内容
type Handler struct {
	errors chan *Error
	f      HandlerFunc
}

// NewHandler 声明新的 Handler 实例
func NewHandler(f HandlerFunc) *Handler {
	h := &Handler{
		errors: make(chan *Error, 100),
		f:      f,
	}

	go h.handle()

	return h
}

// Stop 停止处理错误内容
func (h *Handler) Stop() {
	close(h.errors)
}

// Error 输出一条错误信息
func (h *Handler) Error(err *Error) {
	h.errors <- err
}

// SyntaxError 输出一条语法错误信息
func (h *Handler) SyntaxError(err *Error) {
	err.Type = SyntaxError
	h.Error(err)
}

// SyntaxWarn 输出一条语法警告信息
func (h *Handler) SyntaxWarn(err *Error) {
	err.Type = SyntaxWarn
	h.Error(err)
}

func (h *Handler) handle() {
	for err := range h.errors {
		h.f(err)
	}
}

// NewHandlerFunc 初始 HandlerFunc 实例。
//
// 该实例仅仅是将语法错误和语法警告信息输出到指定的日志通道。
func NewHandlerFunc(errolog, warnlog *log.Logger) HandlerFunc {
	return func(err *Error) {
		switch err.Type {
		case SyntaxError:
			errolog.Println(err)
		case SyntaxWarn:
			warnlog.Println(err)
		case Other:
			panic("代码错误，不应该有其它错误类型")
		}
	}
}
