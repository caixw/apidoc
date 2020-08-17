// SPDX-License-Identifier: MIT

package protocol

// MessageType 传递的消息类型
type MessageType int

// MessageType 可能的值
const (
	MessageTypeError   MessageType = iota + 1 // An error message.
	MessageTypeWarning                        // A warning message.
	MessageTypeInfo                           // An information message.
	MessageTypeLog                            // A log message.
)

// LogMessageParams window/logMessage 传递的参数
type LogMessageParams struct {
	// The message type. See {@link MessageType}
	Type MessageType `json:"type"`

	// The actual message
	Message string `json:"message"`
}
