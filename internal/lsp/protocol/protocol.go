// SPDX-License-Identifier: MIT

// Package protocol 协议内容的定义
package protocol

// DocumentURI https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#uri
type DocumentURI string

// Position https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#position
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#range
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#location
type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

// LocationLink https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#locationlink
type LocationLink struct {
	OriginSelectionRange Range       `json:"originSelectionRange,omitempty"`
	TargetURI            DocumentURI `json:"targetUri,omitempty"`
	TargetRange          Range       `json:"targetRange"`
	TargetSelectionRange Range       `json:"targetSelectionRange"`
}

// Diagnostic https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#diagnostic
type Diagnostic struct {
	Range              Range                          `json:"range"`
	Severity           DiagnosticSeverity             `json:"severity,omitempty"`
	Code               string                         `json:"code,omitempty"`
	Source             string                         `json:"source,omitempty"`
	Message            string                         `json:"message"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

type DiagnosticSeverity int

const (
	DiagnosticSeverityError DiagnosticSeverity = iota + 1
	DiagnosticSeverityWarning
	DiagnosticSeverityInformation
	DiagnosticSeverityHint
)

type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// Command https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#command
type Command struct {
	Title   string        `json:"title"`
	Command string        `json:"command"`
	Args    []interface{} `json:"arguments,omitempty"`
}

// TextEdit https://microsoft.github.io/language-server-protocol/specifications/specification-3-14/#textedit
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}
