// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/lsp/search"
)

// textDocument/didChange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
func (s *server) textDocumentDidChange(notify bool, in *protocol.DidChangeTextDocumentParams, out *interface{}) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	if !f.deleteURI(in.TextDocument.URI) {
		return nil
	}

	f.parsedMux.Lock()
	defer f.parsedMux.Unlock()

	for _, blk := range in.Blocks() {
		f.parseBlock(blk)
	}

	f.srv.textDocumentPublishDiagnostics(f, in.TextDocument.URI)
	return nil
}

// textDocument/hover
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
func (s *server) textDocumentHover(notify bool, in *protocol.HoverParams, out *protocol.Hover) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	f.parsedMux.Lock()
	defer f.parsedMux.Unlock()

	search.Hover(f.doc, in.TextDocument.URI, in.TextDocumentPositionParams.Position, out)
	return nil
}

// textDocument/publishDiagnostics
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_publishDiagnostics
func (s *server) textDocumentPublishDiagnostics(f *folder, uri core.URI) error {
	if s.clientParams.Capabilities.TextDocument.PublishDiagnostics.RelatedInformation == false {
		return nil
	}

	p := &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: make([]protocol.Diagnostic, 0, len(f.errors)+len(f.warns)),
	}

	for _, err := range f.errors {
		if err.Location.URI == uri {
			p.Diagnostics = append(p.Diagnostics, protocol.BuildDiagnostic(err, protocol.DiagnosticSeverityError))
		}
	}

	for _, err := range f.warns {
		if err.Location.URI == uri {
			p.Diagnostics = append(p.Diagnostics, protocol.BuildDiagnostic(err, protocol.DiagnosticSeverityWarning))
		}
	}

	return s.Notify("textDocument/publishDiagnostics", p)
}

// textDocument/foldingRange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_foldingRange
func (s *server) textDocumentFoldingRange(notify bool, in *protocol.FoldingRangeParams, out *[]protocol.FoldingRange) error {
	uri := in.TextDocument.URI

	f := s.findFolder(uri)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, uri)
	}

	lineFoldingOnly := s.clientParams.Capabilities.TextDocument.FoldingRange.LineFoldingOnly

	fr := make([]protocol.FoldingRange, 0, 10)
	if f.doc.URI == uri {
		fr = append(fr, protocol.BuildFoldingRange(f.doc.Base, lineFoldingOnly))
	}

	f.parsedMux.Lock()
	defer f.parsedMux.Unlock()

	for _, api := range f.doc.APIs {
		matched := api.URI == uri || (api.URI == "" && f.doc.URI == uri)
		if !matched {
			continue
		}

		fr = append(fr, protocol.BuildFoldingRange(api.Base, lineFoldingOnly))
	}

	*out = fr

	return nil
}

// textDocument/completion
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
func (s *server) textDocumentCompletion(notify bool, in *protocol.CompletionParams, out *protocol.CompletionList) error {
	// TODO
	return nil
}

// textDocument/semanticTokens
func (s *server) textDocumentSemanticTokens(notify bool, in *protocol.SemanticTokensParams, out *protocol.SemanticTokens) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	f.parsedMux.Lock()
	defer f.parsedMux.Unlock()

	out.Data = search.Tokens(f.doc, in.TextDocument.URI, 0, 1, 2)
	return nil
}

// textDocument/references
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_references
func (s *server) textDocumentReferences(notify bool, in *protocol.ReferenceParams, out *[]core.Location) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	f.parsedMux.Lock()
	defer f.parsedMux.Unlock()

	*out = search.References(f.doc, in.TextDocument.URI, in.Position, in.Context.IncludeDeclaration)
	return nil
}
