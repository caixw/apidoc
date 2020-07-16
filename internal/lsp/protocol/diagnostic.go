// SPDX-License-Identifier: MIT

package protocol

import "github.com/caixw/apidoc/v7/core"

// DiagnosticSeverity 错误级别
type DiagnosticSeverity int

// DiagnosticSeverity 可用的常量
const (
	DiagnosticSeverityError       DiagnosticSeverity = iota + 1 // Reports an error
	DiagnosticSeverityWarning                                   // Reports a warning
	DiagnosticSeverityInformation                               // Reports an information
	DiagnosticSeverityHint                                      // Reports a hint
)

// DiagnosticTag the diagnostic tags.
//
// @since 3.15.0
type DiagnosticTag int

// DiagnosticTag 可用的常量列表
const (
	// DiagnosticTagUnnecessary unused or unnecessary code.
	//
	// Clients are allowed to render diagnostics with this tag faded out instead of having
	// an error squiggle.
	DiagnosticTagUnnecessary DiagnosticTag = 1

	// DiagnosticTagDeprecated deprecated or obsolete code.
	//
	// Clients are allowed to rendered diagnostics with this tag strike through.
	DiagnosticTagDeprecated DiagnosticTag = 2
)

// PublishDiagnosticsClientCapabilities 客户端有关错误信息的支持情况定义
type PublishDiagnosticsClientCapabilities struct {
	// Whether the clients accepts diagnostics with related information.
	RelatedInformation bool `json:"relatedInformation,omitempty"`

	// Client supports the tag property to provide meta data about a diagnostic.
	// Clients supporting tags have to handle unknown tags gracefully.
	//
	// @since 3.15.0
	TagSupport struct {
		// The tags supported by the client.
		ValueSet []DiagnosticTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`

	// Whether the client interprets the version property of the
	// `textDocument/publishDiagnostics` notification's parameter.
	//
	// @since 3.15.0
	VersionSupport bool `json:"versionSupport,omitempty"`
}

// PublishDiagnosticsParams textDocument/publishDiagnostics 事件发送的参数
type PublishDiagnosticsParams struct {
	// The URI for which diagnostic information is reported.
	URI core.URI `json:"uri"`

	// Optional the version number of the document the diagnostics are published for.
	//
	// @since 3.15.0
	Version int `json:"version,omitempty"`

	// An array of diagnostic information items.
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// Diagnostic represents a diagnostic,such as a compiler error or warning.
// Diagnostic objects are only valid in the scope of a resource.
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#diagnostic
type Diagnostic struct {
	// The range at which the message applies.
	Range core.Range `json:"range"`

	// The diagnostic's severity. Can be omitted. If omitted it is up to the
	// client to interpret diagnostics as error, warning, info or hint.
	Severity DiagnosticSeverity `json:"severity,omitempty"`

	// The diagnostic's code, which might appear in the user interface.
	Code string `json:"code,omitempty"`

	// A human-readable string describing the source of this
	// diagnostic, e.g. 'typescript' or 'super lint'.
	Source string `json:"source,omitempty"`

	// The diagnostic's message.
	Message string `json:"message"`

	// Additional metadata about the diagnostic.
	//
	// @since 3.15.0
	Tags []DiagnosticTag `json:"tags,omitempty"`

	// An array of related diagnostic information, e.g. when symbol-names within
	// a scope collide all definitions can be marked via this property.
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

// DiagnosticRelatedInformation represents a related message and source code location for a diagnostic
//
// This should be used to point to code locations that cause or are related to a diagnostics,
// e.g when duplicating a symbol in a scope.
type DiagnosticRelatedInformation struct {
	// The location of this related diagnostic information.
	Location core.Location `json:"location"`

	// The message of this related diagnostic information.
	Message string `json:"message"`
}
