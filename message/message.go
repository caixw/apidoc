// SPDX-License-Identifier: MIT

// Package message 各类输出消息的处理
package message

import "log"

// Type 表示氐的类型
type Type int8

// 消息的分类
const (
	Erro Type = iota
	Warn
)

// Message 输出消息的具体结构
type Message struct {
	Type    Type
	Message string
}

// HandlerFunc 错误处理函数
type HandlerFunc func(*Message)

// Handler 用于接收错误信息内容
type Handler struct {
	messages chan *Message
	f        HandlerFunc
}

// NewHandler 声明新的 Handler 实例
func NewHandler(f HandlerFunc) *Handler {
	h := &Handler{
		messages: make(chan *Message, 100),
		f:        f,
	}

	go h.handle()

	return h
}

// Stop 停止处理错误内容
func (h *Handler) Stop() {
	close(h.messages)
}

// Error 输出一条语法错误信息
func (h *Handler) Error(err error) {
	h.messages <- &Message{
		Type:    Erro,
		Message: err.Error(),
	}
}

// Warn 输出一条语法警告信息
func (h *Handler) Warn(err error) {
	h.messages <- &Message{
		Type:    Warn,
		Message: err.Error(),
	}
}

func (h *Handler) handle() {
	for msg := range h.messages {
		h.f(msg)
	}
}

// NewLogHandlerFunc 生成一个将错误信息输出到日志的 HandlerFunc
//
// 该实例仅仅是将语法错误和语法警告信息输出到指定的日志通道。
func NewLogHandlerFunc(errolog, warnlog *log.Logger) HandlerFunc {
	return func(msg *Message) {
		switch msg.Type {
		case Erro:
			errolog.Println(msg)
		case Warn:
			warnlog.Println(msg)
		default:
			panic("代码错误，不应该有其它错误类型")
		}
	}
}
