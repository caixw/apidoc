// SPDX-License-Identifier: MIT

package lsp

import (
	"fmt"
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/token"
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

	f.doc.DeleteURI(in.TextDocument.URI)
	return f.openFile(in.TextDocument.URI)
}

// textDocument/hover
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
func (s *server) textDocumentHover(notify bool, in *protocol.HoverParams, out *protocol.Hover) error {
	fmt.Println("INFO")
	for _, f := range s.folders {
		f.searchHover(in.TextDocument.URI, in.TextDocumentPositionParams.Position, out)

	}
	return nil
}

func (f *folder) searchHover(uri core.URI, pos core.Position, hover *protocol.Hover) {
	var tip *token.Tip
	if f.doc.URI == uri {
		tip = token.SearchUsage(reflect.ValueOf(f.doc), pos, "APIs")
	}

	for _, api := range f.doc.APIs {
		if api.URI == uri {
			if tip = token.SearchUsage(reflect.ValueOf(api), pos); tip != nil {
				break
			}
		}
	}

	if tip != nil {
		hover.Range = tip.Range
		hover.Contents = protocol.MarkupContent{
			Kind:  protocol.MarkupKinMarkdown,
			Value: tip.Usage,
		}
	}
}
