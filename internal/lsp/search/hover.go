// SPDX-License-Identifier: MIT

package search

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

type usager interface {
	core.Ranger
	Usage() string
}

var usagerType = reflect.TypeOf((*usager)(nil)).Elem()

// Hover 从 doc 查找最符合 uri 和 pos 条件的元素并赋值给 hover
//
// 返回值表示是否找到了相应在的元素。
func Hover(doc *ast.APIDoc, uri core.URI, pos core.Position, h *protocol.Hover) {
	u := doc.Search(uri, pos, usagerType)
	if u != nil {
		usage := u.(usager)
		h.Range = usage.R()
		h.Contents = protocol.MarkupContent{
			Kind:  protocol.MarkupKinMarkdown,
			Value: usage.Usage(),
		}
	}
}
