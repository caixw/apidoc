// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
)

func TestBuildAPIDocOutline(t *testing.T) {
	a := assert.New(t)
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
