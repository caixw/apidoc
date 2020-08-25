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
	RelatedInformation []core.RelatedInformation `json:"relatedInformation,omitempty"`
}

var typeTagsMap = map[core.ErrorType]DiagnosticTag{
	core.ErrorTypeDeprecated: DiagnosticTagDeprecated,
	core.ErrorTypeUnused:     DiagnosticTagUnnecessary,
}

// NewPublishDiagnosticsParams 声明空的 PublishDiagnosticsParams 对象
func NewPublishDiagnosticsParams(uri core.URI) *PublishDiagnosticsParams {
	return &PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []Diagnostic{},
	}
}

// AppendDiagnostic 将 core.Message 添加至诊断数据
func (p *PublishDiagnosticsParams) AppendDiagnostic(err *core.Error, msgType core.MessageType) {
	switch msgType {
	case core.Erro:
		p.Diagnostics = append(p.Diagnostics, buildDiagnostic(err, DiagnosticSeverityError))
	case core.Warn:
		p.Diagnostics = append(p.Diagnostics, buildDiagnostic(err, DiagnosticSeverityWarning))
	case core.Info:
		p.Diagnostics = append(p.Diagnostics, buildDiagnostic(err, DiagnosticSeverityInformation))
	case core.Succ:
		return
	default:
		panic("unreached")
	}
}

func buildDiagnostic(err *core.Error, severity DiagnosticSeverity) Diagnostic {
	var tags []DiagnosticTag
	if len(err.Types) > 0 {
		for _, typ := range err.Types {
			if tag, found := typeTagsMap[typ]; found {
				tags = append(tags, tag)
			}
		}
	}

	msg := err.Error()
	if err.Err != nil {
		msg = err.Err.Error()
	}

	d := Diagnostic{
		Range:    err.Location.Range,
		Message:  msg,
		Severity: severity,
		Source:   core.Name,
	}
	if len(tags) > 0 {
		d.Tags = tags
	}

	if len(err.Related) > 0 {
		d.RelatedInformation = err.Related
	}

	return d
}
