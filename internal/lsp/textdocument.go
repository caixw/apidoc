// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// textDocument/didOpen
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didOpen
func (s *server) textDocumentDidOpen(notify bool, in *protocol.DidOpenTextDocumentParams, out *interface{}) error {
	f := s.getMatchFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	return f.openFile(in.TextDocument.URI)
}

// textDocument/didChange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
func (s *server) textDocumentDidChange(notify bool, in *protocol.DidChangeTextDocumentParams, out *interface{}) error {
	f := s.getMatchFolder(in.TextDocument.URI)
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	for _, content := range in.ContentChanges {
		ok, err := f.matchPosition(in.TextDocument.URI, content.Range.Start)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	f.doc.DeleteFile(in.TextDocument.URI)
	return f.openFile(in.TextDocument.URI)
}

// textDocument/hover
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
func (s *server) textDocumentHover(notify bool, in *protocol.HoverParams, out *protocol.Hover) error {
	// TODO
	return nil
}
