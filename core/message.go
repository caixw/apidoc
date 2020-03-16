// SPDX-License-Identifier: MIT

package core

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/internal/locale"
)

// MessageType 表示氐的类型
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
	Message string
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

// Message 发送普通的文本信息，内容由 key 和 val 组成本地化信息
func (h *MessageHandler) Message(t MessageType, key message.Reference, val ...interface{}) {
	h.messages <- &Message{
		Type:    t,
		Message: locale.Sprintf(key, val...),
	}
}

// Error 将一条错误信息作为消息发送出去
func (h *MessageHandler) Error(t MessageType, err error) {
	h.messages <- &Message{
		Type:    t,
		Message: err.Error(),
	}
}
