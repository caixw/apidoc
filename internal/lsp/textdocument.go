// SPDX-License-Identifier: MIT

package lsp

import (
	"path/filepath"

	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// textDocument/didChange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
func (s *server) textDocumentDidChange(notify bool, in *protocol.DidChangeTextDocumentParams, out *interface{}) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil
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

	if !deleteURI(f.doc, in.TextDocument.URI) {
		return nil
	}

	for _, blk := range in.Blocks() {
		f.parseBlock(blk)
	}

	f.srv.textDocumentPublishDiagnostics(f, in.TextDocument.URI)
	return nil
}

func (f *folder) parseBlock(block core.Block) {
	var input *build.Input
	ext := filepath.Ext(block.Location.URI.String())
	for _, i := range f.cfg.Inputs {
		if sliceutil.Count(i.Exts, func(index int) bool { return i.Exts[index] == ext }) > 0 {
			input = i
			break
		}
	}
	if input == nil { // 无需解析
		return
	}

	f.doc.ParseBlocks(f.h, func(blocks chan core.Block) {
		lang.Parse(f.h, input.Lang, block, blocks)
	})

	if err := f.srv.apidocOutline(f); err != nil {
		f.srv.printErr(err)
	}
}

func deleteURI(doc *ast.APIDoc, uri core.URI) (deleted bool) {
	size := sliceutil.Delete(doc.APIs, func(i int) bool {
		api := doc.APIs[i]
		return api.URI == uri || (api.URI == "" && doc.URI == uri)
	})

	deleted = len(doc.APIs) > size
	doc.APIs = doc.APIs[:size]

	if doc.URI == uri {
		*doc = ast.APIDoc{
			APIs: doc.APIs,
		}
		deleted = true
	}

	return deleted
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
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil
	}

	lineFoldingOnly := s.clientParams.Capabilities.TextDocument.FoldingRange.LineFoldingOnly

	fr := make([]protocol.FoldingRange, 0, 10)
	if f.doc.URI == in.TextDocument.URI {
		fr = append(fr, protocol.BuildFoldingRange(f.doc.Base, lineFoldingOnly))
	}

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	for _, api := range f.doc.APIs {
		matched := api.URI == in.TextDocument.URI || (api.URI == "" && f.doc.URI == in.TextDocument.URI)
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
