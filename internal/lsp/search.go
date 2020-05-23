// SPDX-License-Identifier: MIT

package lsp

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/token"
)

func (s *server) search(p *protocol.HoverParams) *protocol.Hover {
	for _, f := range s.folders {
		hover := f.search(p.TextDocument.URI, p.TextDocumentPositionParams.Position)
		if hover != nil {
			return hover
		}
	}

	return nil
}

func (f *folder) search(uri core.URI, pos core.Position) *protocol.Hover {
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

	if tip == nil {
		return nil
	}
	return &protocol.Hover{
		Range:    tip.Range,
		Contents: tip.Usage,
	}
}
