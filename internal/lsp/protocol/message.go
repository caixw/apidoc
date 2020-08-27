// SPDX-License-Identifier: MIT

package protocol

// SetTraceParams.Value 可用的值
const (
	TraceValueOff     = "off"
	TraceValueMessage = "message"
	TraceValueVerbose = "verbose"
)

// LogTraceParams $/logTrace 服务端下发参数
type LogTraceParams struct {
	// The message to be logged.
	Message string `json:"message"`

	// Additional information that can be computed if the `trace` configuration is set to `'verbose'`
	Verbose string `json:"verbose,omitempty"`
}

// SetTraceParams $/setTrace 入口参数
type SetTraceParams struct {
	// The new value that should be assigned to the trace setting.
	Value string `json:"value"`
}

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

// IsValidTraceValue 是否是一个有效的 TraceValue
func IsValidTraceValue(v string) bool {
	return v == TraceValueOff || v == TraceValueMessage || v == TraceValueVerbose
}

// BuildLogTrace 生成 logTrace 对象
func BuildLogTrace(trace, message, verbose string) *LogTraceParams {
	switch trace {
	case TraceValueOff:
		return nil
	case TraceValueMessage:
		verbose = ""
	case TraceValueVerbose:
	default:
		panic("无效的 trace 值")
	}

	return &LogTraceParams{
		Message: message,
		Verbose: verbose,
	}
}
