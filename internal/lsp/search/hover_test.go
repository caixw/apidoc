// SPDX-License-Identifier: MIT

package search

import (
	"reflect"
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
	hover := &protocol.Hover{}
	pos := core.Position{Line: 1, Character: 1}
	Hover(doc, core.URI("doc.go"), pos, hover)
	a.Equal(hover.Range, core.Range{
		Start: core.Position{Line: 1, Character: 1},
		End:   core.Position{Line: 1, Character: 18},
	})
	a.Equal(hover.Contents.(protocol.MarkupContent).Value, locale.Sprintf("usage-apidoc-title"))

	// apis[0]
	hover = &protocol.Hover{}
	pos = core.Position{Line: 4, Character: 2}
	Hover(doc, core.URI("doc.go"), pos, hover)
	a.Equal(hover.Range, core.Range{
		Start: core.Position{Line: 4, Character: 1},
		End:   core.Position{Line: 7, Character: 7},
	})
	a.Equal(hover.Contents.(protocol.MarkupContent).Value, locale.Sprintf("usage-apidoc-apis"))

	// 改变了 api[0].URI
	doc.APIs[0].URI = core.URI("api0.go")

	// 改变了 api[0].URI，不再匹配 apis[0]，取其父元素 apidoc
	hover = &protocol.Hover{}
	pos = core.Position{Line: 4, Character: 1}
	Hover(doc, core.URI("doc.go"), pos, hover)
	a.Equal(hover.Range, core.Range{
		Start: core.Position{Line: 0, Character: 0},
		End:   core.Position{Line: 12, Character: 9},
	})
	a.Equal(hover.Contents.(protocol.MarkupContent).Value, locale.Sprintf("usage-apidoc"))

	// 与 apis[0] 相同的 URI
	hover = &protocol.Hover{}
	pos = core.Position{Line: 4, Character: 1}
	Hover(doc, core.URI("api0.go"), pos, hover)
	a.Equal(hover.Range, core.Range{
		Start: core.Position{Line: 4, Character: 1},
		End:   core.Position{Line: 7, Character: 7},
	})
	a.Equal(hover.Contents.(protocol.MarkupContent).Value, locale.Sprintf("usage-apidoc-apis"))
}

func TestUsage(t *testing.T) {
	a := assert.New(t)

	b := `<apidoc version="1.1.1">
	<title>标题</title>
	<mimetype>xml</mimetype>
	<mimetype>json</mimetype>
	<api method="GET">
		<path path="/users" />
		<response status="200" />
	</api>
</apidoc>`
	blk := core.Block{Data: []byte(b)}
	rslt := messagetest.NewMessageHandler()
	doc := &ast.APIDoc{}
	doc.Parse(rslt.Handler, blk)
	rslt.Handler.Stop()

	a.Empty(rslt.Errors)
	a.Equal(doc.Start.Line, 0).
		Equal(doc.Start.Character, 0)
	a.Equal(doc.End.Line, 8).
		Equal(doc.End.Character, 9)
	v := reflect.ValueOf(doc)

	pos := core.Position{Line: 0, Character: 1}
	base := usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc"))

	// version
	pos = core.Position{Line: 0, Character: 8}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-version"))

	// 1.1.1
	pos = core.Position{Line: 0, Character: 18}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-version"))

	// title
	pos = core.Position{Line: 1, Character: 1}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-title"))

	// 标题
	pos = core.Position{Line: 1, Character: 10}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-title"))

	// mimetype[0]
	pos = core.Position{Line: 2, Character: 1}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-mimetypes"))

	// mimetype[1]
	pos = core.Position{Line: 3, Character: 2}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-mimetypes"))

	// mimetype[0].xml
	pos = core.Position{Line: 2, Character: 11}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-mimetypes"))

	// apidoc
	pos = core.Position{Line: 8, Character: 1}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc"))

	// exclude api
	pos = core.Position{Line: 4, Character: 2}
	base = usage(v, pos, "APIs")
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc"))

	// api[0]
	pos = core.Position{Line: 4, Character: 2}
	base = usage(v, pos)
	a.NotNil(base).Equal(base.Usage(), locale.Sprintf("usage-apidoc-apis"))
}
