// SPDX-License-Identifier: MIT

package protocol

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestBuildAPIDocOutline(t *testing.T) {
	a := assert.New(t, false)
	f := WorkspaceFolder{Name: "test"}

	doc := &ast.APIDoc{}
	outline := BuildAPIDocOutline(f, doc)
	a.Nil(outline)

	doc = asttest.Get()
	outline = BuildAPIDocOutline(f, doc)
	a.NotNil(outline)

	a.Equal(outline.Title, doc.Title.V())
	a.Equal(2, len(outline.APIs))
}

func TestAPIDocOutline_appendAPI(t *testing.T) {
	a := assert.New(t, false)

	outline := &APIDocOutline{
		APIs: []*API{},
	}

	// api.Path 为空，无法构建 API.Path 变量，忽略该值
	outline.appendAPI(&ast.API{})
	a.Equal(len(outline.APIs), 1)
	api := outline.APIs[0]
	a.Equal(api.Path, "?").Empty(api.Method)

	outline = &APIDocOutline{
		APIs: []*API{},
	}
	outline.appendAPI(&ast.API{Method: &ast.MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}}})
	a.Equal(len(outline.APIs), 1)
	api = outline.APIs[0]
	a.Equal(api.Path, "?").Equal(api.Method, http.MethodDelete)
}
