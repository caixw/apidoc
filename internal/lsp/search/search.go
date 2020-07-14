// SPDX-License-Identifier: MIT

// Package search 实现对 ast.APIDoc 内容的搜索
package search

import (
	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
)

// DeleteURI 删除与 uri 相关的文档内容
//
// deleted 表示是否有内容被删除
func DeleteURI(doc *ast.APIDoc, uri core.URI) (deleted bool) {
	size := sliceutil.Delete(doc.APIs, func(i int) bool {
		api := doc.APIs[i]
		return api.URI == uri || (api.URI == "" && doc.URI == uri)
	})

	deleted = len(doc.APIs) > size
	doc.APIs = doc.APIs[:size]

	if doc.URI == uri {
		*doc = ast.APIDoc{
			APIs: doc.APIs,
		}
		deleted = true
	}

	return deleted
}
