// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/issue9/sliceutil"

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

	f.parsedMux.Lock()
	defer f.parsedMux.Unlock()

	// 清除相关的警告和错误信息
	size := sliceutil.QuickDelete(f.warns, func(i int) bool {
		return f.warns[i].Location.URI == in.TextDocument.URI
	})
	f.warns = f.warns[:size]
	size = sliceutil.QuickDelete(f.errors, func(i int) bool {
		return f.errors[i].Location.URI == in.TextDocument.URI
	})
	f.errors = f.errors[:size]

	if !search.DeleteURI(f.doc, in.TextDocument.URI) {
		return nil
	}

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

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

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

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

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

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

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

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	*out = search.References(f.doc, in.TextDocument.URI, in.Position, in.Context.IncludeDeclaration)
	return nil
}

// textDocument/definition
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_definition
func (s *server) textDocumentDefinition(notify bool, in *protocol.DefinitionParams, out *[]core.Location) error {
	// NOTE: LSP 允许 out 的值是 null，而 jsonrpc 模块默认情况下是空值，而不是 nil，
	// 所以在可能的情况下，都尽量将其返回类型改为数组，
	// 或是像 protocol.Hover 一样为返回类型实现 json.Marshaler 接口。
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	*out = search.Definition(f.doc, in.TextDocument.URI, in.Position)
	return nil
}
