// SPDX-License-Identifier: MIT

package lsp

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func TestHover(t *testing.T) {
	a := assert.New(t)

	b := `<apidoc version="1.1.1">
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
	blk := core.Block{Data: []byte(b), Location: core.Location{URI: "doc.go"}}
	rslt := messagetest.NewMessageHandler()
	doc := &ast.APIDoc{}
	doc.Parse(rslt.Handler, blk)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// title
	h := &protocol.Hover{}
	pos := core.Position{Line: 1, Character: 1}
	hover(doc, core.URI("doc.go"), pos, h)
	a.Equal(h.Range, core.Range{
		Start: core.Position{Line: 1, Character: 1},
		End:   core.Position{Line: 1, Character: 18},
	})
	a.Equal(h.Contents.Value, locale.Sprintf("usage-apidoc-title"))

	// apis[0]
	h = &protocol.Hover{}
	pos = core.Position{Line: 4, Character: 2}
	hover(doc, core.URI("doc.go"), pos, h)
	a.Equal(h.Range, core.Range{
		Start: core.Position{Line: 4, Character: 1},
		End:   core.Position{Line: 7, Character: 7},
	})
	a.Equal(h.Contents.Value, locale.Sprintf("usage-apidoc-apis"))

	// 改变了 api[0].URI
	doc.APIs[0].URI = core.URI("api0.go")

	// 改变了 api[0].URI，不再匹配 apis[0]，取其父元素 apidoc
	h = &protocol.Hover{}
	pos = core.Position{Line: 4, Character: 1}
	hover(doc, core.URI("doc.go"), pos, h)
	a.Equal(h.Range, core.Range{
		Start: core.Position{Line: 0, Character: 0},
		End:   core.Position{Line: 12, Character: 9},
	})
	a.Equal(h.Contents.Value, locale.Sprintf("usage-apidoc"))

	// 与 apis[0] 相同的 URI
	h = &protocol.Hover{}
	pos = core.Position{Line: 4, Character: 1}
	hover(doc, core.URI("api0.go"), pos, h)
	a.Equal(h.Range, core.Range{
		Start: core.Position{Line: 4, Character: 1},
		End:   core.Position{Line: 7, Character: 7},
	})
	a.Equal(h.Contents.Value, locale.Sprintf("usage-apidoc-apis"))
}