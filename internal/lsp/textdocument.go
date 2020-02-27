// SPDX-License-Identifier: MIT

package lsp

import "github.com/caixw/apidoc/v6/internal/lsp/protocol"

// textDocument/didOpen
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didOpen
func (s *server) textDocumentDidOpen(notify bool, in *protocol.DidOpenTextDocumentParams, out *interface{}) error {
	// TODO
	return nil
}

// textDocument/didChange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
func (s *server) textDocumentDidChange(notify bool, in *protocol.DidChangeTextDocumentParams, out *interface{}) error {
	// TODO
	return nil
}

// textDocument/hover
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
func (s *server) textDocumentHover(notify bool, in *protocol.HoverParams, out *protocol.Hover) error {
	// TODO
	return nil
}
