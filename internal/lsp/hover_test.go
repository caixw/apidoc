// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func TestServer_textDocumentHover(t *testing.T) {
	a := assert.New(t, false)
	s := newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	h := &protocol.Hover{}
	err := s.textDocumentHover(false, &protocol.HoverParams{}, h)
	a.Nil(err)

	const b = `<apidoc version="1.1.1">
	<title>标题</title>
	<mimetype>xml</mimetype>
	<mimetype>json</mimetype>
	<api method="GET">
		<path path="/users" />
		<response status="200" />
	</api>
	<api method="POST">
		<path path="/users" />
		<response status="200" />
	</api>
</apidoc>`
	blk := core.Block{Data: []byte(b), Location: core.Location{URI: "file:///test/doc.go"}}
	rslt := messagetest.NewMessageHandler()
	doc := &ast.APIDoc{}
	doc.Parse(rslt.Handler, blk)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	s.folders = []*folder{
		{
			WorkspaceFolder: protocol.WorkspaceFolder{Name: "test", URI: "file:///test"},
			doc:             doc,
		},
	}

	h = &protocol.Hover{}
	err = s.textDocumentHover(false, &protocol.HoverParams{TextDocumentPositionParams: protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file:///test/doc.go"},
		Position:     core.Position{Line: 1, Character: 1},
	}}, h)
	a.NotError(err)
	a.Equal(h.Range, core.Range{
		Start: core.Position{Line: 1, Character: 1},
		End:   core.Position{Line: 1, Character: 18},
	})
	a.Equal(h.Contents.Value, locale.Sprintf("usage-apidoc-title"))
}
