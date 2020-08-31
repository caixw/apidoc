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

	if !deleteURI(f.doc, in.TextDocument.URI) {
		return nil
	}

	f.clearDiagnostics()
	for _, blk := range in.Blocks() {
		f.parseBlock(blk)
	}
	f.srv.textDocumentPublishDiagnostics(f)

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
		return doc.APIs[i].URI == uri
	})

	deleted = len(doc.APIs) > size
	doc.APIs = doc.APIs[:size]

	if doc.URI == uri {
		*doc = ast.APIDoc{APIs: doc.APIs}
		deleted = true
	}

	return deleted
}

// textDocument/publishDiagnostics
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_publishDiagnostics
func (s *server) textDocumentPublishDiagnostics(f *folder) {
	for _, p := range f.diagnostics {
		if err := s.Notify("textDocument/publishDiagnostics", p); err != nil {
			s.erro.Println(err)
		}
	}
}

// 清空所有的诊断信息
func (f *folder) clearDiagnostics() {
	for _, p := range f.diagnostics {
		p.Diagnostics = p.Diagnostics[:0]
		if err := f.srv.Notify("textDocument/publishDiagnostics", p); err != nil {
			f.srv.erro.Println(err)
		}
	}

	f.diagnostics = make(map[core.URI]*protocol.PublishDiagnosticsParams, 0)
}

// textDocument/foldingRange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_foldingRange
func (s *server) textDocumentFoldingRange(notify bool, in *protocol.FoldingRangeParams, out *[]protocol.FoldingRange) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil
	}

	lineFoldingOnly := s.clientParams != nil &&
		s.clientParams.Capabilities.TextDocument.FoldingRange != nil &&
		s.clientParams.Capabilities.TextDocument.FoldingRange.LineFoldingOnly

	folds := make([]protocol.FoldingRange, 0, 10)

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	if f.doc.URI == in.TextDocument.URI {
		folds = append(folds, protocol.BuildFoldingRange(f.doc.Base, lineFoldingOnly))
	}

	for _, api := range f.doc.APIs {
		if api.URI == in.TextDocument.URI {
			folds = append(folds, protocol.BuildFoldingRange(api.Base, lineFoldingOnly))
		}
	}

	*out = folds

	return nil
}

// textDocument/completion
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
func (s *server) textDocumentCompletion(notify bool, in *protocol.CompletionParams, out *protocol.CompletionList) error {
	// TODO
	return nil
}
