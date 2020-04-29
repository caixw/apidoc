// SPDX-License-Identifier: MIT

package protocol

type MessageType int

const (
	Error   MessageType = iota + 1 // An error message.
	Warning                        // A warning message.
	Info                           // An information message.
	Log                            // A log message.
)

type ShowMessageParams struct {
	// The message type. See {@link MessageType}.
	Type MessageType `json:"type"`

	// The actual message.
	Message string `json:"message"`
}

type ShowMessageRequestParams struct {
	// The message type. See {@link MessageType}
	Type MessageType `json:"type"`

	// The actual message
	Message string `json:"message"`

	// The message action items to present.
	Actions []MessageActionItem `json:"actions,omitempty"`
}

type LogMessageParams struct {
	// The message type. See {@link MessageType}
	Type MessageType `json:"type"`

	// The actual message
	Message string `json:"message"`
}

type MessageActionItem struct {
	// A short title like 'Retry', 'Open Log' etc.
	Title string `json:"title"`
}
