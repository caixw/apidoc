// SPDX-License-Identifier: MIT

package search

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
)

var (
	referencerType   = reflect.TypeOf((*ast.Referencer)(nil)).Elem()
	definitionerType = reflect.TypeOf((*ast.Definitioner)(nil)).Elem()
)

// References 返回所在位置的引用列表
func References(doc *ast.APIDoc, uri core.URI, pos core.Position, include bool) (locations []core.Location) {
	r := doc.Search(uri, pos, referencerType)
	if r == nil {
		return
	}

	referencer := r.(ast.Referencer)
	if include {
		locations = append(locations, core.Location{
			URI:   uri,
			Range: referencer.R(),
		})
	}

	for _, ref := range referencer.References() {
		locations = append(locations, ref.Location)
	}

	return
}
