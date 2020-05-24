// SPDX-License-Identifier: MIT

package lsp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/token"
)

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

func (f *folder) matchPosition(uri core.URI, pos core.Position) (bool, error) {
	var r core.Range
	if f.URI == uri {
		r = f.doc.Range
	} else {
		for _, api := range f.doc.APIs {
			if api.URI == uri {
				r = api.Range
				break
			}
		}
	}
	if r.IsEmpty() {
		return false, nil
	}

	return r.Contains(pos), nil
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
