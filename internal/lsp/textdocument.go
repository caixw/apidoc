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
	var f *folder
	for _, f = range s.folders {
		if strings.HasPrefix(string(in.TextDocument.URI), string(f.URI)) {
			break
		}
	}
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	if !search.DeleteURI(f.doc, in.TextDocument.URI) {
		return nil
	}

	for _, blk := range in.Blocks() {
		f.parseBlock(blk)
	}
	return nil
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
