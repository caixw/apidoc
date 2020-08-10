// SPDX-License-Identifier: MIT

package lsp

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestDeleteURI(t *testing.T) {
	a := assert.New(t)

	d := &ast.APIDoc{}
	d.APIDoc = &ast.APIDocVersionAttribute{Value: xmlenc.String{Value: "1.0.0"}}
	d.URI = core.URI("uri1")
	d.APIs = []*ast.API{
		{ //1
			URI: "uri1",
		},
		{ //2
			URI: "uri2",
		},
		{ //3
			URI: "uri3",
		},
		{ //4
		},
	}

	a.True(deleteURI(d, "uri3"))
	a.Equal(3, len(d.APIs)).NotNil(d.APIDoc)

	// 同时会删除 1,4
	a.True(deleteURI(d, "uri1"))
	a.Equal(1, len(d.APIs)).Nil(d.APIDoc)

	a.True(deleteURI(d, "uri2"))
	a.Equal(0, len(d.APIs)).Nil(d.APIDoc)

	a.False(deleteURI(d, "uri2"))
}
