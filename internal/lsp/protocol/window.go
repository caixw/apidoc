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

// ShowMessageParams window/showMessage 的参数
type ShowMessageParams struct {
	// The message type. See {@link MessageType}.
	Type MessageType `json:"type"`

	// The actual message.
	Message string `json:"message"`
}

// ShowMessageRequestParams window/showMessageRequest 传递的参数
type ShowMessageRequestParams struct {
	// The message type. See {@link MessageType}
	Type MessageType `json:"type"`

	// The actual message
	Message string `json:"message"`

	// The message action items to present.
	Actions []MessageActionItem `json:"actions,omitempty"`
}

// LogMessageParams window/logMessage 传递的参数
type LogMessageParams struct {
	// The message type. See {@link MessageType}
	Type MessageType `json:"type"`

	// The actual message
	Message string `json:"message"`
}

// MessageActionItem window/showMessageRequest 传递的按钮描述内容
type MessageActionItem struct {
	// A short title like 'Retry', 'Open Log' etc.
	Title string `json:"title"`
}
