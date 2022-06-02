// SPDX-License-Identifier: MIT

package core

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/internal/locale"
)

// MessageType 表示消息的类型
type MessageType int8

// 消息的分类
const (
	Erro MessageType = iota
	Warn
	Info
	Succ
)

func (t MessageType) String() string {
	switch t {
	case Erro:
		return "ERRO"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Succ:
		return "SUCC"
	default:
		return "<unknown>"
	}
}

// Message 输出消息的具体结构
type Message struct {
	Type    MessageType
	Message interface{}
}

// HandlerFunc 错误处理函数
type HandlerFunc func(*Message)

// MessageHandler 异步的消息处理机制
//
// 包含了本地化的信息，输出时，会以指定的本地化内容输出
type MessageHandler struct {
	messages chan *Message
	stop     chan struct{}
}

// NewMessageHandler 声明新的 MessageHandler 实例
func NewMessageHandler(f HandlerFunc) *MessageHandler {
	h := &MessageHandler{
		messages: make(chan *Message, 100),
		stop:     make(chan struct{}),
	}

	go func() {
		for msg := range h.messages {
			f(msg)
		}
		h.stop <- struct{}{}
	}()

	return h
}

// Stop 停止处理错误内容
//
// 只有在消息处理完成之后，才会返回。
func (h *MessageHandler) Stop() {
	close(h.messages)

	// Stop() 调用可能是在主程序结束处。
	// 通过 h.stop 阻塞函数返回，直到所有消息都处理完成。
	<-h.stop
}

// Message 发送消息
func (h *MessageHandler) Message(t MessageType, msg interface{}) {
	h.messages <- &Message{
		Type:    t,
		Message: msg,
	}
}

// Locale 发送普通的文本信息
func (h *MessageHandler) Locale(t MessageType, key message.Reference, val ...interface{}) {
	h.Message(t, locale.New(key, val...))
}

// Error 发送错误类型的值
func (h *MessageHandler) Error(err interface{}) { h.Message(Erro, err) }

// Warning 发送错误类型的值
func (h *MessageHandler) Warning(err interface{}) { h.Message(Warn, err) }

// Success 发送错误类型的值
func (h *MessageHandler) Success(err interface{}) { h.Message(Succ, err) }

// Info 发送错误类型的值
func (h *MessageHandler) Info(err interface{}) { h.Message(Info, err) }
