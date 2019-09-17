// SPDX-License-Identifier: MIT

// Package message 各类输出消息的处理
package message

import (
	"log"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Type 表示氐的类型
type Type int8

// 消息的分类
const (
	Erro Type = iota
	Warn
	Info
)

// Message 输出消息的具体结构
type Message struct {
	Type    Type
	Message string
}

// HandlerFunc 错误处理函数
type HandlerFunc func(*Message)

// Handler 异步的消息处理机制
//
// 包含了本地化的信息，输出时，会以指定的本地化内容输出
type Handler struct {
	tag     language.Tag
	printer *message.Printer

	messages chan *Message
	f        HandlerFunc
}

// NewHandler 声明新的 Handler 实例
func NewHandler(f HandlerFunc, tag language.Tag) *Handler {
	h := &Handler{
		tag:     tag,
		printer: message.NewPrinter(tag),

		messages: make(chan *Message, 100),
		f:        f,
	}

	go func() {
		for msg := range h.messages {
			h.f(msg)
		}
	}()

	return h
}

// Tag 关联的 language.Tag
func (h *Handler) Tag() language.Tag {
	return h.tag
}

// Printer 关联的 message.Printer
func (h *Handler) Printer() *message.Printer {
	return h.printer
}

// Stop 停止处理错误内容
func (h *Handler) Stop() {
	close(h.messages)
}

// Message 发送普通的文本信息，内容由 key 和 val 组成本地化信息
func (h *Handler) Message(t Type, key message.Reference, val ...interface{}) {
	h.messages <- &Message{
		Type:    t,
		Message: h.printer.Sprintf(key, val...),
	}
}

// Error 将一条错误信息作为消息发送出去
func (h *Handler) Error(t Type, err error) {
	msg := &Message{Type: t}

	if l, ok := err.(LocaleError); ok {
		msg.Message = l.LocaleError(h.printer)
	} else {
		msg.Message = err.Error()
	}

	h.messages <- msg
}

// NewLogHandlerFunc 生成一个将错误信息输出到日志的 HandlerFunc
//
// 该实例仅仅是将语法错误和语法警告信息输出到指定的日志通道。
func NewLogHandlerFunc(errolog, warnlog, infolog *log.Logger) HandlerFunc {
	return func(msg *Message) {
		switch msg.Type {
		case Erro:
			errolog.Println(msg)
		case Warn:
			warnlog.Println(msg)
		case Info:
			infolog.Println(msg)
		default:
			panic("代码错误，不应该有其它错误类型")
		}
	}
}
