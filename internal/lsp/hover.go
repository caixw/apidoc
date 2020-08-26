// SPDX-License-Identifier: MIT

package lsp

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

type usager interface {
	core.Searcher
	Usage() string
}

var usagerType = reflect.TypeOf((*usager)(nil)).Elem()

func hover(doc *ast.APIDoc, uri core.URI, pos core.Position, h *protocol.Hover) {
	u := doc.Search(uri, pos, usagerType)
	if u == nil {
		return
	}

	usage := u.(usager)
	if v := usage.Usage(); v != "" {
		h.Range = usage.Loc().Range
		h.Contents = protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: v,
		}
	}
}

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

	hover(f.doc, in.TextDocument.URI, in.TextDocumentPositionParams.Position, out)
	return nil
}
