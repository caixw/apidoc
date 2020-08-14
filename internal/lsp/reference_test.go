// SPDX-License-Identifier: MIT

package lsp

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func loadReferencesDoc(a *assert.Assertion) *ast.APIDoc {
	const referenceDefinitionDoc = `<apidoc version="1.1.1">
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

	blk := core.Block{Data: []byte(referenceDefinitionDoc), Location: core.Location{URI: "file:///root/doc.go"}}
	rslt := messagetest.NewMessageHandler()
	doc := &ast.APIDoc{}
	doc.Parse(rslt.Handler, blk)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	return doc
}

func TestServer_textDocumentReferences(t *testing.T) {
	a := assert.New(t)
	s := &server{}
	var locs []core.Location
	err := s.textDocumentReferences(false, &protocol.ReferenceParams{}, &locs)
	a.Nil(err).Empty(locs)

	s.folders = []*folder{
		{
			WorkspaceFolder: protocol.WorkspaceFolder{Name: "test", URI: "file:///root"},
			doc:             loadReferencesDoc(a),
		},
	}

	err = s.textDocumentReferences(false, &protocol.ReferenceParams{TextDocumentPositionParams: protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file:///root/doc.go"},
		Position:     core.Position{Line: 3, Character: 16},
	}}, &locs)
	a.NotError(err).Equal(len(locs), 2)
	a.Equal(locs[0], core.Location{
		URI: "file:///root/doc.go",
		Range: core.Range{
			Start: core.Position{Line: 6, Character: 2},
			End:   core.Position{Line: 6, Character: 15},
		},
	})
}

func TestServer_textDocumentDefinition(t *testing.T) {
	a := assert.New(t)
	s := &server{}
	var locs []core.Location
	err := s.textDocumentDefinition(false, &protocol.DefinitionParams{}, &locs)
	a.Nil(err).Empty(locs)

	s.folders = []*folder{
		{
			WorkspaceFolder: protocol.WorkspaceFolder{Name: "test", URI: "file:///root"},
			doc:             loadReferencesDoc(a),
		},
	}

	err = s.textDocumentDefinition(false, &protocol.DefinitionParams{TextDocumentPositionParams: protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file:///root/doc.go"},
		Position:     core.Position{Line: 6, Character: 2},
	}}, &locs)
	a.NotError(err).Equal(len(locs), 1)
	a.Equal(locs[0], core.Location{
		URI: "file:///root/doc.go",
		Range: core.Range{
			Start: core.Position{Line: 3, Character: 1},
			End:   core.Position{Line: 3, Character: 31},
		},
	})
}

func TestReferences(t *testing.T) {
	a := assert.New(t)
	doc := loadReferencesDoc(a)

	pos := core.Position{}
	locs := references(doc, "file:///root/doc.go", pos, false)
	a.Nil(locs)

	pos = core.Position{Line: 3, Character: 16}
	locs = references(doc, "file:///root/doc.go", pos, false)
	a.Equal(len(locs), 2).
		Equal(locs[0], core.Location{
			URI: "file:///root/doc.go",
			Range: core.Range{
				Start: core.Position{Line: 6, Character: 2},
				End:   core.Position{Line: 6, Character: 15},
			},
		})

	pos = core.Position{Line: 3, Character: 16}
	locs = references(doc, "file:///root/doc.go", pos, true)
	a.Equal(len(locs), 3)
}

func TestDefinition(t *testing.T) {
	a := assert.New(t)
	doc := loadReferencesDoc(a)

	pos := core.Position{}
	locs := definition(doc, "file:///root/doc.go", pos)
	a.Empty(locs)

	pos = core.Position{Line: 3, Character: 16}
	locs = definition(doc, "file:///root/doc.go", pos)
	a.Empty(locs)

	pos = core.Position{Line: 6, Character: 2}
	locs = definition(doc, "file:///root/doc.go", pos)
	a.Equal(locs, []core.Location{
		{
			URI: "file:///root/doc.go",
			Range: core.Range{
				Start: core.Position{Line: 3, Character: 1},
				End:   core.Position{Line: 3, Character: 31},
			},
		},
	})

	pos = core.Position{Line: 12, Character: 2}
	locs = definition(doc, "file:///root/doc.go", pos)
	a.Equal(locs, []core.Location{
		{
			URI: "file:///root/doc.go",
			Range: core.Range{
				Start: core.Position{Line: 4, Character: 1},
				End:   core.Position{Line: 4, Character: 31},
			},
		},
	})
}
