// SPDX-License-Identifier: MIT

// Package search 实现对 ast.APIDoc 内容的搜索
package search

import (
	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
)

// DeleteURI 删除与 uri 相关的文档内容
func DeleteURI(doc *ast.APIDoc, uri core.URI) {
	size := sliceutil.Delete(doc.APIs, func(i int) bool {
		api := doc.APIs[i]
		return api.URI == uri || (api.URI == "" && doc.URI == uri)
	})
	doc.APIs = doc.APIs[:size]

	if doc.URI == uri {
		*doc = ast.APIDoc{
			APIs: doc.APIs,
		}
	}
}
