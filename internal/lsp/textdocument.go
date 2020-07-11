// SPDX-License-Identifier: MIT

package lsp

import (
	"strings"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/lsp/search"
)

// textDocument/didChange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
func (s *server) textDocumentDidChange(notify bool, in *protocol.DidChangeTextDocumentParams, out *interface{}) error {
	f := s.getMatchFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	matched := f.doc.URI == in.TextDocument.URI
	if !matched {
		for _, api := range f.doc.APIs {
			if api.URI == in.TextDocument.URI {
				matched = true
				break
			}
		} // end for
	}

	if !matched {
		return nil
	}

	f.doc.DeleteURI(in.TextDocument.URI)
	return f.openFile(in.TextDocument.URI)
}

func (s *server) getMatchFolder(uri core.URI) *folder {
	for _, f := range s.folders {
		if f.matchURI(uri) {
			return f
		}
	}
	return nil
}

// uri 是否与属于项目匹配
func (f *folder) matchURI(uri core.URI) bool {
	return strings.HasPrefix(string(uri), string(f.URI))
}

// textDocument/hover
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
func (s *server) textDocumentHover(notify bool, in *protocol.HoverParams, out *protocol.Hover) error {
	for _, f := range s.folders {
		if search.Hover(f.doc, in.TextDocument.URI, in.TextDocumentPositionParams.Position, out) {
			return nil
		}
	}
	return nil
}

// textDocument/publishDiagnostics
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_publishDiagnostics
func (s *server) textDocumentPublishDiagnostics(uri core.URI, errs []*core.SyntaxError, warns []*core.SyntaxError) error {

	if s.clientCapabilities.TextDocument.PublishDiagnostics.RelatedInformation == false {
		return nil
	}

	p := &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: make([]protocol.Diagnostic, 0, len(errs)+len(warns)),
	}

	for _, err := range errs {
		p.Diagnostics = append(p.Diagnostics, protocol.Diagnostic{
			Range:    err.Location.Range,
			Message:  err.Error(),
			Severity: protocol.DiagnosticSeverityError,
		})
	}

	for _, warn := range warns {
		p.Diagnostics = append(p.Diagnostics, protocol.Diagnostic{
			Range:    warn.Location.Range,
			Message:  warn.Error(),
			Severity: protocol.DiagnosticSeverityWarning,
		})
	}

	return s.Notify("textDocument/publishDiagnostics", p)
}
