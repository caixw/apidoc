// SPDX-License-Identifier: MIT

package search

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/token"
)

func TestAPIDoc_DeleteURI(t *testing.T) {
	a := assert.New(t)

	d := &ast.APIDoc{}
	d.APIDoc = &ast.APIDocVersionAttribute{Value: token.String{Value: "1.0.0"}}
	d.URI = core.URI("uri1")
	d.APIs = []*ast.API{
		{ //1
			URI: core.URI("uri1"),
		},
		{ //2
			URI: core.URI("uri2"),
		},
		{ //3
			URI: core.URI("uri3"),
		},
		{ //4
		},
	}

	a.True(DeleteURI(d, core.URI("uri3")))
	a.Equal(3, len(d.APIs)).NotNil(d.APIDoc)

	// 同时会删除 1,4
	a.True(DeleteURI(d, core.URI("uri1")))
	a.Equal(1, len(d.APIs)).Nil(d.APIDoc)

	a.True(DeleteURI(d, core.URI("uri2")))
	a.Equal(0, len(d.APIs)).Nil(d.APIDoc)

	a.False(DeleteURI(d, core.URI("uri2")))
}
