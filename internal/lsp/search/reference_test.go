// SPDX-License-Identifier: MIT

package search

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
)

const testDoc = `<apidoc version="1.1.1">
	<title>标题</title>
	<mimetype>xml</mimetype>
	<tag name="t1" title="tag1" />
	<tag name="t2" title="tag2" />
	<api method="GET">
		<tag>t1</tag>
		<path path="/users" />
		<response status="200" />
	</api>
	<api method="POST">
		<tag>t1</tag>
		<tag>t2</tag>
		<path path="/users" />
		<response status="200" />
	</api>
</apidoc>`

func loadDoc(a *assert.Assertion) *ast.APIDoc {
	blk := core.Block{Data: []byte(testDoc), Location: core.Location{URI: "doc.go"}}
	rslt := messagetest.NewMessageHandler()
	doc := &ast.APIDoc{}
	doc.Parse(rslt.Handler, blk)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	return doc
}

func TestReferences(t *testing.T) {
	a := assert.New(t)
	doc := loadDoc(a)

	pos := core.Position{}
	locs := References(doc, "doc.go", pos, false)
	a.Nil(locs)

	pos = core.Position{Line: 3, Character: 16}
	locs = References(doc, "doc.go", pos, false)
	a.Equal(len(locs), 2).
		Equal(locs[0], core.Location{
			URI: "doc.go",
			Range: core.Range{
				Start: core.Position{Line: 6, Character: 2},
				End:   core.Position{Line: 6, Character: 15},
			},
		})

	pos = core.Position{Line: 3, Character: 16}
	locs = References(doc, "doc.go", pos, true)
	a.Equal(len(locs), 3)
}

func TestDefinition(t *testing.T) {
	a := assert.New(t)
	doc := loadDoc(a)

	pos := core.Position{}
	loc := Definition(doc, "doc.go", pos)
	a.Equal(loc, core.Location{})

	pos = core.Position{Line: 3, Character: 16}
	loc = Definition(doc, "doc.go", pos)
	a.Equal(loc, core.Location{})

	pos = core.Position{Line: 6, Character: 2}
	loc = Definition(doc, "doc.go", pos)
	a.Equal(loc, core.Location{
		URI: "doc.go",
		Range: core.Range{
			Start: core.Position{Line: 3, Character: 1},
			End:   core.Position{Line: 3, Character: 31},
		},
	})

	pos = core.Position{Line: 12, Character: 2}
	loc = Definition(doc, "doc.go", pos)
	a.Equal(loc, core.Location{
		URI: "doc.go",
		Range: core.Range{
			Start: core.Position{Line: 4, Character: 1},
			End:   core.Position{Line: 4, Character: 31},
		},
	})
}
