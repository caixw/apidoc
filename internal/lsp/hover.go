// SPDX-License-Identifier: MIT

package lsp

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

type usager interface {
	core.Searcher
	Usage() string
}

var usagerType = reflect.TypeOf((*usager)(nil)).Elem()

// textDocument/hover
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
func (s *server) textDocumentHover(notify bool, in *protocol.HoverParams, out *protocol.Hover) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil // 非项目文件，不应该出错

	}

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	if u := f.doc.Search(in.TextDocument.URI, in.TextDocumentPositionParams.Position, usagerType); u != nil {
		usage := u.(usager)
		if v := usage.Usage(); v != "" {
			out.Range = usage.Loc().Range
			out.Contents = protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: v,
			}
		}
	}
	return nil
}
